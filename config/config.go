package config

import "strings"

type Directive struct {
	Name string
	Args []string
	Raw  string
}

type Config struct {
	Directives map[string][]Directive
}

func (c *Config) Get(name string) []Directive {
	if c == nil {
		return nil
	}
	name = strings.ToLower(name)
	return c.Directives[name]
}

func BuildDefaultConfig() *Config {
	return &Config{
		Directives: map[string][]Directive{
			"port":                  {{Name: "port", Args: []string{"6379"}}},
			"bind":                  {{Name: "bind", Args: []string{"127.0.0.1"}}},
			"databases":             {{Name: "databases", Args: []string{"16"}}},
			"hz":                    {{Name: "hz", Args: []string{"10"}}},
			"active-expire-samples": {{Name: "active-expire-samples", Args: []string{"20"}}},
		},
	}
}
