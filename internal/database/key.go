package database

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-playground/validator"
	"github.com/internet24/shadowsocks24/pkg/utils"
	"golang.org/x/exp/slices"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const KeyPath = "storage/database/keys.json"

type Key struct {
	Id      string `json:"id" validate:"required,hostname"`
	Cipher  string `json:"cipher" validate:"required,oneof=chacha20-ietf-poly1305 aes-128-gcm aes-256-gcm"`
	Secret  string `json:"secret" validate:"required,min=6,max=64"`
	Name    string `json:"name" validate:"required,min=1,max=64"`
	Quota   int64  `json:"quota" validate:"min=0"`
	Enabled bool   `json:"enabled"`
}

type KeyTable struct {
	Keys      []*Key `json:"keys" validate:"required"`
	NextId    int64  `json:"next_id" validate:"required,min=1"`
	UpdatedAt int64  `json:"updated_at" validate:"min=0"`
}

func (kt *KeyTable) Load() error {
	content, err := os.ReadFile(KeyPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if !utils.DirectoryExist(filepath.Dir(SettingPath)) {
				return errors.New(fmt.Sprintf("directory %s not found", filepath.Base(KeyPath)))
			}
			return kt.Save()
		}
		return errors.New(fmt.Sprintf("cannot load %s, err: %v", KeyPath, err))
	}

	err = json.Unmarshal(content, kt)
	if err != nil {
		return err
	}

	if err = validator.New().Struct(kt); err != nil {
		return errors.New(fmt.Sprintf("cannot validate %s, err: %v", KeyPath, err))
	}

	return nil
}

func (kt *KeyTable) Save() (err error) {
	defer func() {
		_ = kt.Load()
	}()

	if err = validator.New().Struct(kt); err != nil {
		return DataError(err.Error())
	}
	for _, k := range kt.Keys {
		if err = validator.New().Struct(k); err != nil {
			return DataError(err.Error())
		}
	}

	kt.UpdatedAt = time.Now().Unix()
	content, err := json.Marshal(kt)
	if err != nil {
		return err
	}

	if err = os.WriteFile(KeyPath, content, 0755); err != nil {
		return errors.New(fmt.Sprintf("cannot save %s, err: %v", KeyPath, err))
	}

	return kt.Load()
}

func (kt *KeyTable) Store(key Key) (*Key, error) {
	for _, k := range kt.Keys {
		if k.Secret == key.Secret {
			return nil, DataError(fmt.Sprintf("The secret `%s` already exists.", k.Secret))
		}
	}

	key.Id = fmt.Sprintf("k-%d", kt.NextId)

	kt.NextId++
	kt.Keys = append(kt.Keys, &key)

	return &key, kt.Save()
}

func (kt *KeyTable) Update(key Key) (*Key, error) {
	for _, k := range kt.Keys {
		if k.Id != key.Id && k.Secret == key.Secret {
			return nil, DataError(fmt.Sprintf("The secret %s already exists.", k.Secret))
		}
	}

	for i, k := range kt.Keys {
		if k.Id == key.Id {
			kt.Keys[i].Cipher = key.Cipher
			kt.Keys[i].Secret = key.Secret
			kt.Keys[i].Name = key.Name
			kt.Keys[i].Quota = key.Quota
			kt.Keys[i].Enabled = key.Enabled
			return kt.Keys[i], kt.Save()
		}
	}

	return nil, nil
}

func (kt *KeyTable) Refill(keys []Key) (err error) {
	var nextId int64 = 1
	for _, k := range keys {
		if err = validator.New().Struct(k); err != nil {
			return DataError(err.Error())
		}
		for _, k2 := range keys {
			if k.Id != k2.Id && k.Secret == k2.Secret {
				return DataError(fmt.Sprintf("The secret of %s and %s is %s.", k.Id, k2.Id, k.Secret))
			}
		}
		if nextId, err = strconv.ParseInt(k.Id[2:], 10, 64); err != nil {
			return DataError(fmt.Sprintf("Invalid key ID: %v", k.Id))
		}
	}

	kt.Keys = []*Key{}
	for i := range keys {
		kt.Keys = append(kt.Keys, &keys[i])
	}
	kt.NextId = nextId

	return kt.Save()
}

func (kt *KeyTable) Delete(id string) error {
	for i, k := range kt.Keys {
		if k.Id == id {
			kt.Keys = slices.Delete(kt.Keys, i, i+1)
			return kt.Save()
		}
	}
	return nil
}
