package config

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func Parse(r io.Reader) (*Config, error) {
	cfg := &Config{
		Directives: make(map[string][]Directive),
	}
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if idx := strings.Index(line, "#"); idx >= 0 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		d := Directive{
			Name: strings.ToLower(fields[0]),
			Args: fields[1:],
			Raw:  line,
		}
		cfg.Directives[d.Name] = append(cfg.Directives[d.Name], d)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return cfg, nil
}

func ParseFile(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

func ParseArgs(args []string) (*Config, error) {
	cfg := &Config{
		Directives: make(map[string][]Directive),
	}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--") {
			arg = arg[2:]
		} else if strings.HasPrefix(arg, "-") {
			arg = arg[1:]
		} else {
			fileCfg, err := ParseFile(arg)
			if err != nil {
				return nil, err
			}
			for name, directives := range fileCfg.Directives {
				cfg.Directives[name] = append(cfg.Directives[name], directives...)
			}
			continue
		}
		name := strings.ToLower(arg)
		var value string
		if i+1 < len(args) {
			value = args[i+1]
			i++
		}
		d := Directive{
			Name: name,
			Args: []string{value},
			Raw:  arg + " " + value,
		}
		cfg.Directives[d.Name] = append(cfg.Directives[d.Name], d)
	}
	return cfg, nil
}
