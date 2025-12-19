package config

import (
	"strconv"

	"github.com/ghosind/antdb/server"
)

var optionsBuilders = map[string]func(*Config) server.ServerOption{
	"port": func(cfg *Config) server.ServerOption { return buildIntOption(cfg, "port", server.WithPort) },
	"bind": func(cfg *Config) server.ServerOption { return buildStringOption(cfg, "bind", server.WithBind) },
}

func BuildOptionsByConfig(cfg *Config) []server.ServerOption {
	var options []server.ServerOption

	for _, builder := range optionsBuilders {
		option := builder(cfg)
		if option != nil {
			options = append(options, option)
		}
	}

	return options
}

func buildIntOption(cfg *Config, name string, setter func(int) server.ServerOption) server.ServerOption {
	directives := cfg.Get(name)
	if len(directives) == 0 || len(directives[0].Args) == 0 {
		return nil
	}
	value, err := strconv.Atoi(directives[0].Args[0])
	if err != nil {
		return nil
	}
	return setter(value)
}

func buildStringOption(cfg *Config, name string, setter func(string) server.ServerOption) server.ServerOption {
	directives := cfg.Get(name)
	if len(directives) == 0 || len(directives[0].Args) == 0 {
		return nil
	}
	return setter(directives[0].Args[0])
}
