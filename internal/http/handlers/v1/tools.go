package v1

import (
	"encoding/json"
	"github.com/internet24/shadowsocks24/internal/coordinator"
	"github.com/internet24/shadowsocks24/internal/database"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"io"
	"net/http"
	"strconv"
)

type OutlineImportRequest struct {
	URL string `json:"url" validate:"url"`
}

type OutlineResponse struct {
	AccessKeys []struct {
		Name      string `json:"name"`
		Password  string `json:"password"`
		Port      int    `json:"port"`
		Method    string `json:"method"`
		DataLimit struct {
			Bytes int64 `json:"bytes"`
		}
	} `json:"accessKeys"`
}

func OutlineImport(coordinator *coordinator.Coordinator) echo.HandlerFunc {
	return func(c echo.Context) error {
		var r OutlineImportRequest
		if err := c.Bind(&r); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Cannot parse the request body.",
			})
		}
		if err := c.Validate(&r); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": err.Error(),
			})
		}

		request, err := http.NewRequest("GET", r.URL+"/access-keys", nil)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Invalid URL!",
			})
		}
		request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)

		response, err := coordinator.Http.Do(request)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Cannot load!",
			})
		}
		defer func() {
			_ = response.Body.Close()
		}()

		if response.StatusCode != http.StatusOK {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Error " + strconv.Itoa(response.StatusCode),
			})
		}

		body, err := io.ReadAll(response.Body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Cannot read!",
			})
		}

		var or OutlineResponse
		if err = json.Unmarshal(body, &or); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Cannot parse!",
			})
		}

		for _, ak := range or.AccessKeys {
			for _, k := range coordinator.Database.KeyTable.Keys {
				if k.Secret == ak.Password {
					if err = coordinator.Database.KeyTable.Delete(k.Id); err != nil {
						coordinator.Logger.Error("cannot delete key", zap.Error(err))
						return c.JSON(http.StatusBadRequest, map[string]string{
							"message": "Incomplete import!",
						})
					}
				}
			}
			key := database.Key{
				Cipher:  ak.Method,
				Secret:  ak.Password,
				Name:    ak.Name,
				Quota:   int(ak.DataLimit.Bytes / 1000000),
				Enabled: true,
			}
			if _, err = coordinator.Database.KeyTable.Store(key); err != nil {
				coordinator.Logger.Error("cannot store key", zap.Error(err))
				return c.JSON(http.StatusBadRequest, map[string]string{
					"message": "Incomplete import!",
				})
			}
		}

		return c.NoContent(http.StatusNoContent)
	}
}
