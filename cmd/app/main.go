package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/ibookerke/choco_parser_go/internal/config"
	"github.com/ibookerke/choco_parser_go/internal/pkg/pgx"
	"github.com/ibookerke/choco_parser_go/internal/pkg/trm/manager"
	"github.com/ibookerke/choco_parser_go/internal/repository"
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
