package coordinator

import (
	"fmt"
	"github.com/internet24/shadowsocks24/internal/config"
	"github.com/internet24/shadowsocks24/internal/database"
	"github.com/internet24/shadowsocks24/pkg/prometheus"
	"github.com/internet24/shadowsocks24/pkg/shadowsocks"
	"github.com/internet24/shadowsocks24/pkg/utils"
	"go.uber.org/zap"
	"net/http"
)

type Coordinator struct {
	Http          *http.Client
	Logger        *zap.Logger
	Config        *config.Config
	Prometheus    *prometheus.Prometheus
	Shadowsocks   *shadowsocks.Shadowsocks
	Database      *database.Database
	IP            string
	MetricsPort   int
	ServerMetrics map[string]*ServerMetric
	KeyMetrics    map[string]*KeyMetric
	SyncedAt      int64
}

func (c *Coordinator) Run() {
	c.initSettings()
	c.initMetricsPort()
	c.syncKeys(false)
	c.syncServers(false)
	go c.Shadowsocks.Run(c.MetricsPort)
	go c.Prometheus.Reload()
	go c.startWorkers()
}

func (c *Coordinator) Sync() {
	c.syncKeys(true)
	c.syncServers(true)
}

func (c *Coordinator) initSettings() {
	c.IP = utils.IP()
	local := "127.0.0.1"
	if c.IP != local && c.Database.SettingTable.ShadowsocksHost == local {
		c.Database.SettingTable.ShadowsocksHost = c.IP
	}

	if c.Database.SettingTable.ExternalHttp == "http://127.0.0.1:80" {
		c.Database.SettingTable.ExternalHttp = fmt.Sprintf("http://%s:%d", c.IP, c.Config.HttpServer.Port)
	}

	if c.Database.SettingTable.ApiToken == "SECRET0123456789" {
		c.Database.SettingTable.ApiToken = utils.Token()
	}

	if c.Database.SettingTable.ShadowsocksPort == 1 {
		var err error
		if c.Database.SettingTable.ShadowsocksPort, err = utils.FreePort(); err != nil {
			c.Logger.Fatal("cannot find a free port for the shadowsocks server", zap.Error(err))
		}
	}

	if err := c.Database.SettingTable.Save(); err != nil {
		c.Logger.Fatal("cannot save settings", zap.Error(err))
	}
}

func (c *Coordinator) initMetricsPort() {
	var err error
	if c.MetricsPort, err = utils.FreePort(); err != nil {
		c.Logger.Fatal("cannot find a free port for the shadowsocks metrics", zap.Error(err))
	}
}

func New(
	c *config.Config, l *zap.Logger, hc *http.Client, p *prometheus.Prometheus, db *database.Database,
	ss *shadowsocks.Shadowsocks,
) *Coordinator {
	return &Coordinator{
		Config:        c,
		Logger:        l,
		Http:          hc,
		Database:      db,
		Prometheus:    p,
		Shadowsocks:   ss,
		ServerMetrics: map[string]*ServerMetric{},
		KeyMetrics:    map[string]*KeyMetric{},
	}
}
