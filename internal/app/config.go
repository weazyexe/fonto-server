package app

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Postgres struct {
	Host   string `yaml:"host"`
	Port   string `yaml:"port"`
	DbName string `yaml:"dbname"`
}

type FontoConfig struct {
	Postgres Postgres `yaml:"postgres"`
	Port     string   `yaml:"port"`
}

func ReadConfig(path string) (*FontoConfig, error) {
	configYaml, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := FontoConfig{}
	if err := yaml.Unmarshal(configYaml, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
