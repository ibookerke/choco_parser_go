package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/ibookerke/choco_parser_go/internal/config"
	"github.com/ibookerke/choco_parser_go/internal/delivery/http"
	"github.com/ibookerke/choco_parser_go/internal/pkg/pgx"
	"github.com/ibookerke/choco_parser_go/internal/pkg/trm/manager"
	"github.com/ibookerke/choco_parser_go/internal/repository"
	"github.com/ibookerke/choco_parser_go/internal/server"
	"github.com/ibookerke/choco_parser_go/internal/service"
)

func main() {
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	conf, err := config.Get()
	if err != nil {
		slog.Error("couldn't get config", "err", err)
		return
	}

	slogHandler := slog.Handler(slog.NewTextHandler(os.Stdout, nil))
	if !conf.Project.Debug {
		slogHandler = slog.NewJSONHandler(os.Stdout, nil)
	}

	logger := slog.New(slogHandler).With("svc", conf.Project.ServiceName)

	pool, err := pgx.NewPgxPool(ctx, conf.Database.DSN)
	if err != nil || pool == nil {
		logger.Error("couldn't create pgx pool", "err", err)
		return
	}
	defer pool.Close()

	// migrating database scheme using migrate library
	m, err := migrate.New("file://migrations", conf.Database.DSN)
	if err != nil {
		logger.Error("couldn't create migrate instance", "err", err)
		return
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error("couldn't migrate database", "err", err)
		return
	}

	trManager := manager.Must(pgx.NewDefaultFactory(pool))

	branchRepo := repository.NewBranchRepository(pool, pgx.DefaultCtxGetter, trManager)
	customerRepo := repository.NewCustomerRepository(pool, pgx.DefaultCtxGetter, trManager)
	paymentRepo := repository.NewPaymentRepository(pool, pgx.DefaultCtxGetter, trManager)
	authRepo := repository.NewAuthRepository(pool, pgx.DefaultCtxGetter, trManager)
	companyCustomerRepo := repository.NewCompanyCustomerRepository(pool, pgx.DefaultCtxGetter, trManager)

	branchService := service.NewBranchService(branchRepo, authRepo, trManager, conf.Choco)
	paymentService := service.NewPaymentService(paymentRepo, customerRepo, authRepo, trManager, conf.Choco)
	companyCustomerService := service.NewCompanyCustomersService(companyCustomerRepo, trManager, conf.Choco)

	e := echo.New()

	e.IPExtractor = echo.ExtractIPFromRealIPHeader()
	e.Use(middleware.Decompress())
	e.Use(middleware.Gzip())
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(5))) // 5 requests per sec
	e.Use(middleware.Recover())
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())
	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		Timeout:      30 * time.Second,
		ErrorMessage: "request timed out",
		OnTimeoutRouteErrorHandler: func(err error, c echo.Context) {
			logger.Error("request timed out", "err", err, "path", c.Path())
		},
	}))

	http.NewPaymentsHandler(logger, conf.Choco, *paymentService).RegisterRoutes(e)

	httpServer := server.NewHTTPServer(logger, conf.Rest, e)
	go func() {
		if err := httpServer.Run(); err != nil {
			logger.Error("failed to start http server", "err", err)
			cancelFn()
		}
	}()

	if len(os.Args) < 2 {
		fmt.Println("invalid number of parameters passed")
		fmt.Println("It should be in the format: go run cmd/main/go <commandAction>")
		return
	}

	commandName := os.Args[0]
	commandAction := os.Args[1]

	fmt.Println("commandName: ", commandName)
	fmt.Println("commandAction: ", commandAction)

	switch commandAction {
	case "customers":
		fetchCustomers(ctx, branchService, paymentService)
		break
	case "company_customers":
		company_name := os.Args[2]
		fetchCompanyCustomers(ctx, branchService, companyCustomerService, company_name)
	default:
		fmt.Println("invalid command name")
		break
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case v := <-quit:
		logger.Info("received exit signal", "signal", v)
		cancelFn()
	case <-ctx.Done():
		logger.Info("context done")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err = httpServer.Shutdown(ctx); err != nil {
		logger.Error("failed to shutdown server", "err", err)
	}
	logger.Info("server stopped")

}

func fetchCompanyCustomers(
	ctx context.Context,
	branchService *service.BranchService,
	companyCustomerService *service.CompanyCustomersService,
	companyName string,
) {
	_, err := branchService.FetchBranches(ctx)
	if err != nil {
		fmt.Println("error fetching terminals: ", err)
		return
	}

	terminals, err := branchService.GetBranchTerminals(ctx, companyName)
	if err != nil {
		fmt.Println("error fetching terminals: ", err)
		return
	}

	err = companyCustomerService.FetchCompanyCustomers(ctx, terminals, companyName)
	if err != nil {
		fmt.Println("error fetching company customers: ", err)
		return
	}
}

func fetchCustomers(
	ctx context.Context,
	branchService *service.BranchService,
	paymentService *service.PaymentService,

) {
	terminals, err := branchService.FetchBranches(ctx)
	if err != nil {
		fmt.Println("error fetching terminals: ", err)
		return
	}

	err = paymentService.FetchPayments(ctx, terminals)
	if err != nil {
		fmt.Println("error fetching payments: ", err)
		return
	}

	fmt.Println("fetching branches and payments completed successfully")
}
