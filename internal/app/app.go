package app

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/labstack/echo/v4"
	v1 "github.com/zd4r/cbr_exchange_rate/internal/controller/http/v1"
	"github.com/zd4r/cbr_exchange_rate/internal/usecase"
	"github.com/zd4r/cbr_exchange_rate/internal/usecase/webapi"
	"github.com/zd4r/cbr_exchange_rate/pkg/config"
	"github.com/zd4r/cbr_exchange_rate/pkg/httpserver"
	"github.com/zd4r/cbr_exchange_rate/pkg/logger"
)

func Run(cfg *config.Config) {
	// Logger
	l := logger.New(cfg.Log.Level, cfg.Log.Structured)
	l.Info().Msg(fmt.Sprintf("⚡ init app [ name: %v, version: %v ]", cfg.Name, cfg.Version))

	// Use case
	collectionUseCase := usecase.NewDynamicQuote(
		webapi.New(),
	)

	// HTTP server
	handler := echo.New()
	v1.NewRouter(handler, l, collectionUseCase)
	httpServer := httpserver.New(handler, httpserver.Port(cfg.HTTP.Port))
	l.Info().Msg(fmt.Sprintf("🌏 http server started on :%v", cfg.HTTP.Port))

	// Graceful shutdown
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, os.Interrupt)

	select {
	case s := <-interrupt:
		l.Info().Msg(fmt.Sprintf("app - Run - signal: %v", s.String()))
	case err := <-httpServer.Notify():
		l.Err(err).Msg("app - Run - httpServer.Notify")
	}

	err := httpServer.Shutdown()
	if err != nil {
		l.Err(err).Msg("app - Run - httpServer.Shutdown")
	}
}
