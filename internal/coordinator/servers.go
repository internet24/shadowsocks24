package coordinator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/internet24/shadowsocks24/internal/database"
	"github.com/internet24/shadowsocks24/pkg/utils"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"io"
	"net/http"
	"time"
)

func (c *Coordinator) SyncServers(reconfigure bool) {
	c.Logger.Debug("syncing servers with the prometheus service")

	servers := map[string]string{
		"s-0": fmt.Sprintf("127.0.0.1:%d", c.MetricsPort),
	}
	for _, s := range c.Database.ServerTable.Servers {
		servers[s.Id] = fmt.Sprintf("%s:%d", s.HttpHost, s.HttpPort)
	}

	if err := c.Prometheus.Update(servers); err != nil {
		c.Logger.Fatal("cannot sync servers with the prometheus service", zap.Error(err))
	}
	if reconfigure {
		c.Prometheus.Reload()
	}

	go c.pushServers()
}

func (c *Coordinator) CurrentServer() *database.Server {
	return &database.Server{
		Id:                 "s-0",
		HttpHost:           utils.IP(),
		HttpPort:           c.Config.HttpServer.Port,
		ShadowsocksEnabled: c.Database.SettingTable.ShadowsocksEnabled,
		ShadowsocksHost:    c.Database.SettingTable.ShadowsocksHost,
		ShadowsocksPort:    c.Database.SettingTable.ShadowsocksPort,
		ApiToken:           c.Database.SettingTable.ApiToken,
		Status:             database.ServerStatusActive,
		SyncedAt:           time.Now().Unix(),
	}
}

func (c *Coordinator) updateServerStatus(s *database.Server, newStatus string) {
	s.Status = newStatus
	if _, err := c.Database.ServerTable.Update(*s); err != nil {
		c.Logger.Error("cannot update server status", zap.String("server", s.Id), zap.Error(err))
	}
}

func (c *Coordinator) pullServers() {
	for _, s := range c.Database.ServerTable.Servers {
		c.pullServer(s)
	}
}

func (c *Coordinator) pullServer(s *database.Server) {
	url := fmt.Sprintf("http://%s:%d/v1/settings", s.HttpHost, s.HttpPort)

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		c.Logger.Error("cannot create health request", zap.String("url", url), zap.Error(err))
		c.updateServerStatus(s, database.ServerStatusUnavailable)
		return
	}

	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Add(echo.HeaderAuthorization, "Bearer "+s.ApiToken)

	response, err := c.Http.Do(request)
	if err != nil {
		c.Logger.Error("cannot check server health", zap.String("url", url), zap.Error(err))
		c.updateServerStatus(s, database.ServerStatusUnavailable)
		return
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusUnauthorized {
		c.updateServerStatus(s, database.ServerStatusUnauthorized)
		return
	}

	if response.StatusCode == http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			c.Logger.Error("cannot read pulled server", zap.String("url", url), zap.Error(err))
			c.updateServerStatus(s, database.ServerStatusUnavailable)
			return
		}

		var settings database.SettingTable
		if err = json.Unmarshal(body, &settings); err != nil {
			c.Logger.Error(
				"cannot unmarshall pulled server", zap.String("url", url),
				zap.Error(err), zap.String("body", string(body)),
			)
			c.updateServerStatus(s, database.ServerStatusUnavailable)
			return
		}

		s.Status = database.ServerStatusActive
		s.ShadowsocksEnabled = settings.ShadowsocksEnabled
		s.ShadowsocksHost = settings.ShadowsocksHost
		s.ShadowsocksPort = settings.ShadowsocksPort
		if _, err = c.Database.ServerTable.Update(*s); err != nil {
			c.Logger.Error("cannot update server status", zap.String("server", s.Id), zap.Error(err))
		}

		return
	}

	c.updateServerStatus(s, database.ServerStatusUnavailable)
	_ = response.Body.Close()

}

func (c *Coordinator) pushServers() {
	for _, s := range c.Database.ServerTable.Servers {
		if s.Status == database.ServerStatusUnknown || s.SyncedAt < c.Database.KeyTable.UpdatedAt {
			c.pushServer(s)
		}
	}
}

func (c *Coordinator) pushServer(s *database.Server) {
	url := fmt.Sprintf("http://%s:%d/v1/keys/refill", s.HttpHost, s.HttpPort)
	c.Logger.Debug("pushing keys to server...", zap.String("url", url))

	body, err := json.Marshal(c.Database.KeyTable.Keys)
	if err != nil {
		c.Logger.Fatal("cannot marshal database.keys", zap.Error(err))
	}

	request, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		c.Logger.Error("cannot create push request", zap.String("url", url), zap.Error(err))
		c.updateServerStatus(s, database.ServerStatusUnavailable)
		return
	}

	request.Header.Add(echo.HeaderContentType, echo.MIMEApplicationJSON)
	request.Header.Add(echo.HeaderAuthorization, "Bearer "+s.ApiToken)

	response, err := c.Http.Do(request)
	if err != nil {
		c.Logger.Error("cannot push keys to server", zap.String("url", url), zap.Error(err))
		c.updateServerStatus(s, database.ServerStatusUnavailable)
		return
	}
	defer func() {
		_ = response.Body.Close()
	}()

	if response.StatusCode == http.StatusUnauthorized {
		c.updateServerStatus(s, database.ServerStatusUnauthorized)
		return
	}

	if response.StatusCode != http.StatusNoContent {
		c.updateServerStatus(s, database.ServerStatusUnavailable)
		return
	}

	s.Status = database.ServerStatusActive
	s.SyncedAt = time.Now().Unix()
	if _, err = c.Database.ServerTable.Update(*s); err != nil {
		c.Logger.Error("cannot update server", zap.String("server", s.Id), zap.Error(err))
	}
}
