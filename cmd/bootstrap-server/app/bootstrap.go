package app

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/criticalstack/crit/pkg/cluster/bootstrap"
)

type bootstrapConfig struct {
	Provider   string
	Filters    map[string]string
	Kubeconfig string
}

type bootstrapRouter struct {
	*echo.Echo
	cfg *bootstrapConfig
}

func newBootstrapRouter(cfg *bootstrapConfig) *bootstrapRouter {
	r := &bootstrapRouter{
		Echo: echo.New(),
		cfg:  cfg,
	}
	r.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Skipper: func(c echo.Context) bool {
			return c.Path() == "/healthz"
		},
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	r.Use(middleware.Recover())
	r.Use(middleware.BodyLimit("2M"))
	r.Use(middleware.Secure())

	r.GET("/healthz", func(c echo.Context) error {
		return c.JSON(http.StatusOK, nil)
	})

	r.GET("/authorize", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"provider": cfg.Provider,
		})
	})
	r.POST("/authorize", func(c echo.Context) error {
		var auth bootstrap.Request
		if err := c.Bind(&auth); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": err.Error(),
			})
		}
		switch auth.Type {
		case bootstrap.AmazonIdentityDocumentAndSignature:
			return r.handleAmazonIdentityDocumentAndSignature(auth.Body)(c)
		default:
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": fmt.Sprintf("unknown auth type: %q", auth.Type)},
			)
		}
	})
	return r
}
