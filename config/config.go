package config

import "github.com/BurntSushi/toml"

type Config struct {
	AWS           AWS           `toml:"aws"`
	StepFunctions StepFunctions `toml:"step_functions"`
}

type AWS struct {
	AccessKeyID     string `toml:"access_key_id"`
	SecretAccessKey string `toml:"secret_access_key"`
	Region          string `toml:"region"`
}

type StepFunctions struct {
	Arn string `toml:"arn"`
}

var cfg *Config

func Load(pathConfig string) (*Config, error) {

	cfg = &Config{}
	_, err := toml.DecodeFile(pathConfig, &cfg)
	if err != nil {
		return nil, err
	}

	return cfg, err
}
