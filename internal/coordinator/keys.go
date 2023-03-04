package coordinator

import (
	"github.com/internet24/shadowsocks24/pkg/shadowsocks"
	"go.uber.org/zap"
	"time"
)

func (c *Coordinator) SyncKeys(reconfigure bool) {
	c.Logger.Debug("syncing keys with the shadowsocks service...")

	keys := make([]shadowsocks.Key, 0, len(c.Database.KeyTable.Keys))
	for _, k := range c.Database.KeyTable.Keys {
		if !k.Enabled {
			continue
		}
		keys = append(keys, shadowsocks.Key{
			Id:     k.Id,
			Secret: k.Secret,
			Cipher: k.Cipher,
			Port:   c.Database.SettingTable.ShadowsocksPort,
		})
	}

	if err := c.Shadowsocks.Update(keys); err != nil {
		c.Logger.Fatal("cannot sync keys with the shadowsocks service", zap.Error(err))
	}

	if reconfigure {
		c.Shadowsocks.Reconfigure()
	}

	c.SyncedAt = time.Now().Unix()
}
