package config

import (
	"sync"

	"github.com/hexiaopi/rdm-toy/internal/model"
	"github.com/spf13/pflag"
)

var AppEngine AppConfig

type AppConfig struct {
	Redis map[string]model.Conn
}

func (o *AppConfig) Flags(fs *pflag.FlagSet) {

}

var Conns sync.Map
