package config

import (
	m "main/models"
	"sync"

	"github.com/spf13/viper"
)

var (
	serviceInstance *ConfigService
	once            sync.Once
)

type ConfigService struct {
	mu     sync.RWMutex
	config m.Config
}

func Initialize() error {
	var initErr error
	once.Do(func() {
		v := viper.New()
		v.AddConfigPath("config/")
		v.SetConfigName("config")
		v.SetConfigType("toml")

		if err := v.ReadInConfig(); err != nil {
			initErr = err
			return
		}

		var cfg m.Config
		if err := v.Unmarshal(&cfg); err != nil {
			initErr = err
			return
		}

		serviceInstance = &ConfigService{
			config: cfg,
		}
	})
	return initErr
}
func Get() *ConfigService {
	return serviceInstance
}

func (cs *ConfigService) Config() m.Config {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.config
}
