package coordinator

import (
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
	MetricsPort   int
	ServerMetrics map[string]*ServerMetric
	KeyMetrics    map[string]*KeyMetric
	SyncedAt      int64
}

func (c *Coordinator) Run() {
	c.SyncKeys(false)

	var err error
	if c.MetricsPort, err = utils.FreePort(); err != nil {
		c.Logger.Fatal("cannot find a free port for the shadowsocks metrics", zap.Error(err))
	}
	go c.Shadowsocks.Run(c.MetricsPort)

	c.SyncServers(false)

	go c.Prometheus.Reload()

	go c.startWorker()
}

func (c *Coordinator) Sync() {
	c.SyncKeys(true)
	c.SyncServers(true)
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
