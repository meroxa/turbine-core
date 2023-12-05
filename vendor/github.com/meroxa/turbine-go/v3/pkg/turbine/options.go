package turbine

import (
	"maps"
)

type OptionConfig struct {
	Name                 string
	PluginConfig         map[string]string
	PlatformPluginConfig []string
}

type Option func(*OptionConfig)

func WithName(pluginName string) Option {
	return func(cfg *OptionConfig) {
		cfg.Name = pluginName
	}
}

func WithPluginConfig(pluginConfig map[string]string) Option {
	return func(cfg *OptionConfig) {
		maps.Copy(pluginConfig, cfg.PluginConfig)
	}
}

func WithPlatformConfig(ref string) Option {
	return func(cfg *OptionConfig) {
		cfg.PlatformPluginConfig = append(cfg.PlatformPluginConfig, ref)
	}
}
