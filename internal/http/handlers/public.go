package handlers

import (
	"github.com/internet24/shadowsocks24/internal/coordinator"
	"github.com/labstack/echo/v4"
	"net/http"
	"os"
)

func Public(coordinator *coordinator.Coordinator) echo.HandlerFunc {
	return func(c echo.Context) error {
		content, err := os.ReadFile("web/public.html")
		if err != nil {
			return err
		}

		return c.HTML(http.StatusOK, string(content))
	}
}
