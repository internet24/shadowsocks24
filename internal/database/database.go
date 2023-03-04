package database

import (
	"fmt"
	"github.com/internet24/shadowsocks24/internal/config"
	"github.com/internet24/shadowsocks24/pkg/utils"
)

type DataError string

func (de DataError) Error() string {
	return string(de)
}

type Database struct {
	SettingTable *SettingTable
	KeyTable     *KeyTable
	ServerTable  *ServerTable
}

func New(c *config.Config) (*Database, error) {
	db := &Database{
		SettingTable: &SettingTable{
			AdminPassword:      "password",
			ApiToken:           utils.Token(),
			ShadowsocksHost:    utils.IP(),
			ShadowsocksPort:    1000,
			ShadowsocksEnabled: true,
			ExternalHttps:      "",
			ExternalHttp:       fmt.Sprintf("http://%s:%d", utils.IP(), c.HttpServer.Port),
			TrafficRatio:       1,
		},
		KeyTable: &KeyTable{
			Keys:   []*Key{},
			NextId: 1,
		},
		ServerTable: &ServerTable{
			Servers: []*Server{},
			NextId:  1,
		},
	}

	if err := db.SettingTable.Load(); err != nil {
		return nil, err
	}
	if err := db.KeyTable.Load(); err != nil {
		return nil, err
	}
	if err := db.ServerTable.Load(); err != nil {
		return nil, err
	}

	return db, nil
}
