package config

import (
	"strconv"

	"github.com/ghosind/antdb/server"
)

type ServerOptionParamType int

const (
	ServerOptionParamTypeInt ServerOptionParamType = iota
	ServerOptionParamTypeString
)

type ServerOptionParam struct {
	Name          string
	Type          ServerOptionParamType
	OptionBuilder any
}

func (p *ServerOptionParam) BuildOption(cfg *Config) server.ServerOption {
	switch p.Type {
	case ServerOptionParamTypeInt:
		return buildIntOption(cfg, p.Name, p.OptionBuilder.(func(int) server.ServerOption))
	case ServerOptionParamTypeString:
		return buildStringOption(cfg, p.Name, p.OptionBuilder.(func(string) server.ServerOption))
	default:
		return nil
	}
}

var optionParams = map[string]ServerOptionParam{
	"port":      {Name: "port", Type: ServerOptionParamTypeInt, OptionBuilder: server.WithPort},
	"bind":      {Name: "bind", Type: ServerOptionParamTypeString, OptionBuilder: server.WithBind},
	"databases": {Name: "databases", Type: ServerOptionParamTypeInt, OptionBuilder: server.WithDatabases},
	"hz":        {Name: "hz", Type: ServerOptionParamTypeInt, OptionBuilder: server.WithHZ},
	"active-expire-samples": {
		Name:          "active-expire-samples",
		Type:          ServerOptionParamTypeInt,
		OptionBuilder: server.WithActiveExpireSamples,
	},
	"requirepass": {
		Name:          "requirepass",
		Type:          ServerOptionParamTypeString,
		OptionBuilder: server.WithRequirePass,
	},
}

func BuildOptionsByConfig(cfg *Config) []server.ServerOption {
	var options []server.ServerOption

	for _, param := range optionParams {
		option := param.BuildOption(cfg)
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
