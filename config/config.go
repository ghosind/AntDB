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
