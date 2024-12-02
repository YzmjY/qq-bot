package config

import "github.com/jinzhu/configor"

type Config struct {
	AgentConfig AgentConfig `yaml:"agent_config"`
}

type AgentConfig struct {
	Type string `yaml:"type"`
	Addr string `yaml:"addr"`
	Port string `yaml:"port"`
}

func MustNewConfig(path string) Config {
	cfg := configor.Config{
		Debug: true,
	}

	var config Config
	err := configor.New(&cfg).Load(&config, path)
	if err != nil {
		panic(err)
	}
	return config
}
