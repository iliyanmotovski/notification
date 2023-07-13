package config

import (
	"os"

	"github.com/gookit/config"
	"github.com/gookit/config/yaml"
)

type Config interface {
	Get(string) (interface{}, bool)
}

func NewConfigService(appName, localConfigFolder string) (Config, error) {
	c := config.NewEmpty(appName)
	c.AddDriver(yaml.Driver)

	yamlFileAppConfig, err := os.ReadFile(localConfigFolder + "/config.yaml")
	if err != nil {
		return nil, err
	}

	err = c.LoadSources(config.Yaml, yamlFileAppConfig)
	if err != nil {
		return nil, err
	}

	return &gookitConfig{c}, nil
}

type gookitConfig struct {
	*config.Config
}

func (g *gookitConfig) Get(key string) (interface{}, bool) {
	return g.Config.Get(key)
}
