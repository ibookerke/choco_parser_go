package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/ibookerke/choco_parser_go/internal/config"
)

type HttpServer struct {
	logger *slog.Logger
	conf   config.Rest
	srv    *http.Server

	hnd http.Handler
}

func NewHTTPServer(logger *slog.Logger, conf config.Rest, handler http.Handler) *HttpServer {
	return &HttpServer{
		logger: logger,
		conf:   conf,
		hnd:    handler,
	}
}

func (s *HttpServer) Run() error {
	address := fmt.Sprintf("%s:%d", s.conf.Host, s.conf.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		s.logger.Error("error listening on address", "err", err)
		return err
	}

	defer listener.Close()

	s.srv = &http.Server{
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 40 << 20, //
		Handler:        cors(s.hnd, s.conf.AllowedCORSOrigins),
	}

	s.logger.Info(fmt.Sprintf("HTTP server is running on %s", address))
	if err = s.srv.Serve(listener); err != nil {
		if errors.Is(err, http.ErrServerClosed) {
			s.logger.Info("HTTP server closed")
			return err
		}
		s.logger.Error("error serving", "err", err)
		return err
	}

	return nil
}

func (s *HttpServer) Shutdown(ctx context.Context) error {
	return s.srv.Shutdown(ctx)
}

func cors(h http.Handler, allowedOrigins []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		providedOrigin := r.Header.Get("Origin")
		matches := false
		for _, allowedOrigin := range allowedOrigins {
			if allowedOrigin == "*" {
				matches = true
				break
			}
			if providedOrigin == allowedOrigin {
				matches = true
				break
			}
		}

		if matches {
			w.Header().Set("Access-Control-Allow-Origin", providedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization, ResponseType, X-Request-ID, x-payload-digest, x-authorization-digest")
		}
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h.ServeHTTP(w, r)
	})
}
