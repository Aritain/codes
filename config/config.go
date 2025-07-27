package config

import (
	"sync"

	"github.com/spf13/viper"
)

var (
	serviceInstance *ConfigService
	once            sync.Once
)

type ConfigService struct {
	mu     sync.RWMutex
	config Config
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

		var cfg Config
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

func (cs *ConfigService) Config() Config {
	cs.mu.RLock()
	defer cs.mu.RUnlock()
	return cs.config
}

type Config struct {
	TGToken        string `mapstructure:"TG_TOKEN"`
	BotDebug       bool   `mapstructure:"BOT_DEBUG"`
	GCodesEndpoint string `mapstructure:"GCODES_ENDPOINT"`
	HCodesEndpoint string `mapstructure:"HCODES_ENDPOINT"`
	ZCodesEndpoint string `mapstructure:"ZCODES_ENDPOINT"`
	GCodesName     string `mapstructure:"GCODES_NAME"`
	HCodesName     string `mapstructure:"HCODES_NAME"`
	ZCodesName     string `mapstructure:"ZCODES_NAME"`
	CodesDir       string `mapstructure:"CODES_DIR"`
}
