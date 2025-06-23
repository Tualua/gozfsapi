package main

import (
	"fmt"
	"os"

	"github.com/Tualua/gozfsapi/internal/server"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	oapimw "github.com/oapi-codegen/echo-middleware"
	"go.uber.org/zap"
)

func main() {
	swagger, err := server.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}
	swagger.Servers = nil
	logger, err := zap.NewProduction()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error creating logger: %s", err)
		os.Exit(1)
	}
	defer logger.Sync()
	zfsApi := server.NewApiServer(logger)
	e := echo.New()
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:    true,
		LogStatus: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger.Info("request",
				zap.String("URI", v.URI),
				zap.Int("status", v.Status),
			)
			return nil
		},
	}))
	e.Use(oapimw.OapiRequestValidator(swagger))
	server.RegisterHandlers(e, zfsApi)
	logger.Fatal("starting server",
		zap.Error(e.Start(":8080")),
	)
}
