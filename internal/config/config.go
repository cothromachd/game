package config

import (
	"go.uber.org/multierr"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Srv HTTPServer `yaml:"httpSrv"`
	DB  DB         `yaml:"db"`
}

type HTTPServer struct {
	Addr string `yaml:"addr"`
}

type DB struct {
	Conn string `yaml:"conn"`
}

func New(path string) (_ *Config, err error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0755)
	if err != nil {
		return nil, err
	}
	defer multierr.AppendInto(&err, file.Close())

	var cfg Config
	return &cfg, yaml.NewDecoder(file).Decode(&cfg)
}
