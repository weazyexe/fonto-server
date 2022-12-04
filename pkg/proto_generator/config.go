package proto_generator

import (
	"gopkg.in/yaml.v3"
	"os"
)

type ProtoConfig struct {
	SourcesDirectory  string `yaml:"sourcesDirectory"`
	GenerateDirectory string `yaml:"generateDirectory"`
}

func ReadConfig(path string) (*ProtoConfig, error) {
	configYaml, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := ProtoConfig{}
	if err := yaml.Unmarshal(configYaml, &config); err != nil {
		return nil, err
	}
	return &config, nil
}
