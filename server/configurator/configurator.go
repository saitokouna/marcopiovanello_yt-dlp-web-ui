package configurator

import (
	"log/slog"
	"os"
	"path/filepath"
	"sync"

	"github.com/marcopiovanello/yt-dlp-web-ui/v3/server/config"
	"gopkg.in/yaml.v3"
)

// A singleton holding configuration of the frontend component
// with optional persistence on a file.
type AppConfig struct {
	Title          string `yaml:"title" json:"title"`
	BaseURL        string `yaml:"base_url" json:"base_url"`
	Language       string `yaml:"language" json:"language"`
	RPCPollingTime int    `yaml:"rpc_polling_time" json:"rpc_polling_time"`
}

type Configurator struct {
	mu     sync.RWMutex
	Config AppConfig
}

var (
	instance     *Configurator
	instanceOnce sync.Once
)

func Instance() *Configurator {
	instanceOnce.Do(func() {
		if instance == nil {
			instance = &Configurator{}

			// TODO: move out of initialization
			err := instance.Load()
			if err != nil {
				slog.Error("failed initializating configurator", slog.Any("err", err))
			}
		}
	})
	return instance
}

func (c *Configurator) Load() error {
	fd, err := getConfigurationFile()
	if err != nil {
		return err
	}
	defer fd.Close()

	if err := yaml.NewDecoder(fd).Decode(&c.Config); err != nil {
		return err
	}

	return nil
}

func (c *Configurator) Persist() error {
	fd, err := getConfigurationFile()
	if err != nil {
		return err
	}
	defer fd.Close()

	if err := yaml.NewEncoder(fd).Encode(c.Config); err != nil {
		return err
	}

	return nil
}

func (c *Configurator) setAppConfig(ac *AppConfig) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// TODO: better validaitons
	if ac.BaseURL != "" {
		c.Config.BaseURL = ac.BaseURL
	}
	if ac.Language != "" {
		c.Config.Language = ac.Language
	}
	if ac.Title != "" {
		c.Config.Title = ac.Title
	}
	if ac.RPCPollingTime >= 250 && ac.RPCPollingTime <= 2000 {
		c.Config.RPCPollingTime = ac.RPCPollingTime
	}
}

func getConfigurationFile() (*os.File, error) {
	fd, err := os.OpenFile(
		filepath.Join(config.Instance().Dir(), "web_config.yml"),
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644,
	)
	if err != nil {
		return nil, err
	}

	return fd, nil
}
