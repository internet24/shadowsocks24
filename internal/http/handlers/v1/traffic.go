package v1

import (
	b64 "encoding/base64"
	"github.com/internet24/shadowsocks24/internal/coordinator"
	"github.com/internet24/shadowsocks24/internal/database"
	"github.com/labstack/echo/v4"
	"net/http"
	"strings"
)

type TrafficResponse struct {
	coordinator.KeyMetric
	Name  string `json:"name"`
	Quota int    `json:"quota"`
}

func Traffic(cdr *coordinator.Coordinator) echo.HandlerFunc {
	return func(c echo.Context) error {
		var query, err = b64.StdEncoding.DecodeString(c.QueryParam("q"))
		if err != nil {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "Not found.",
			})
		}

		parts := strings.Split(string(query), ":")
		if len(parts) != 2 {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "Not found.",
			})
		}

		var key *database.Key
		for _, k := range cdr.Database.KeyTable.Keys {
			if k.Cipher == parts[0] && k.Secret == parts[1] {
				key = k
			}
		}
		if key == nil {
			return c.JSON(http.StatusNotFound, map[string]string{
				"message": "Not found.",
			})
		}

		var r TrafficResponse
		r.Name = key.Name
		r.Quota = key.Quota
		if m, found := cdr.KeyMetrics[key.Id]; found {
			r.KeyMetric = coordinator.KeyMetric{
				Id:      m.Id,
				DownTcp: int64(float64(m.DownTcp)*cdr.Database.SettingTable.TrafficRatio) / 1000000,
				DownUdp: int64(float64(m.DownUdp)*cdr.Database.SettingTable.TrafficRatio) / 1000000,
				UpTcp:   int64(float64(m.UpTcp)*cdr.Database.SettingTable.TrafficRatio) / 1000000,
				UpUdp:   int64(float64(m.UpUdp)*cdr.Database.SettingTable.TrafficRatio) / 1000000,
				Total:   int64(float64(m.Total)*cdr.Database.SettingTable.TrafficRatio) / 1000000,
			}
		}

		return c.JSON(http.StatusOK, r)
	}
}
