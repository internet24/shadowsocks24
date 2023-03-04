package v1

import (
	"github.com/internet24/shadowsocks24/internal/coordinator"
	"github.com/internet24/shadowsocks24/internal/database"
	"github.com/labstack/echo/v4"
	"net/http"
)

type KeysStoreRequest struct {
	Cipher  string `json:"cipher"`
	Secret  string `json:"secret"`
	Name    string `json:"name"`
	Quota   int    `json:"quota"`
	Enabled bool   `json:"enabled"`
}

type KeysUpdateRequest struct {
	KeysStoreRequest
	Id string `json:"id"`
}

type KeyResponse struct {
	database.Key
	Used int64 `json:"used"`
}

func KeysIndex(coordinator *coordinator.Coordinator) echo.HandlerFunc {
	return func(c echo.Context) error {
		keys := make([]KeyResponse, 0, len(coordinator.Database.KeyTable.Keys))
		for _, k := range coordinator.Database.KeyTable.Keys {
			key := KeyResponse{Key: *k}
			if m, found := coordinator.KeyMetrics[k.Id]; found {
				key.Used = m.Total / 1000000
			}
			keys = append(keys, key)
		}

		return c.JSON(http.StatusOK, keys)
	}
}

func KeysStore(coordinator *coordinator.Coordinator) echo.HandlerFunc {
	return func(c echo.Context) error {
		var r KeysStoreRequest
		if err := c.Bind(&r); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Cannot parse the request body.",
			})
		}

		key, err := coordinator.Database.KeyTable.Store(database.Key{
			Cipher:  r.Cipher,
			Secret:  r.Secret,
			Name:    r.Name,
			Quota:   r.Quota,
			Enabled: r.Enabled,
		})
		if err != nil {
			if _, ok := err.(database.DataError); ok {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"message": err.Error(),
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Internal error.",
			})
		}

		go coordinator.Sync()

		return c.JSON(http.StatusCreated, KeyResponse{Key: *key, Used: 0})
	}
}

func KeysUpdate(coordinator *coordinator.Coordinator) echo.HandlerFunc {
	return func(c echo.Context) error {
		var r KeysUpdateRequest
		if err := c.Bind(&r); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Cannot parse the request body.",
			})
		}

		key, err := coordinator.Database.KeyTable.Update(database.Key{
			Id:      r.Id,
			Cipher:  r.Cipher,
			Secret:  r.Secret,
			Name:    r.Name,
			Quota:   r.Quota,
			Enabled: r.Enabled,
		})
		if err != nil {
			if _, ok := err.(database.DataError); ok {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"message": err.Error(),
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Internal error.",
			})
		}
		if key == nil {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "Key not found.",
			})
		}

		go coordinator.Sync()

		kr := KeyResponse{Key: *key}
		if m, found := coordinator.KeyMetrics[kr.Id]; found {
			kr.Used = m.Total / 1000000
		}

		return c.JSON(http.StatusOK, kr)
	}
}

func KeysDelete(coordinator *coordinator.Coordinator) echo.HandlerFunc {
	return func(c echo.Context) error {
		id := c.Param("id")

		err := coordinator.Database.KeyTable.Delete(id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"message": "Cannot update the database.",
			})
		}

		go coordinator.Sync()

		return c.NoContent(http.StatusNoContent)
	}
}

func KeysRefill(coordinator *coordinator.Coordinator) echo.HandlerFunc {
	return func(c echo.Context) error {
		var r []database.Key
		if err := c.Bind(&r); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"message": "Cannot parse the request body.",
			})
		}

		err := coordinator.Database.KeyTable.Refill(r)
		if err != nil {
			if _, ok := err.(database.DataError); ok {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"message": err.Error(),
				})
			}
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"message": "Internal error.",
			})
		}

		go coordinator.Sync()

		return c.NoContent(http.StatusNoContent)
	}
}
