package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Github        Github          `yaml:"github"`
	SecretScraper []SecretScraper `yaml:"secret_scraper"`
}

type Github struct {
	AccessToken string `yaml:"access_token"`
}

type SecretScraper struct {
	SecretProvider     string `yaml:"secret_provider"`
	SecretType         string `yaml:"secret_type"`
	SecretQueryKeyword string `yaml:"secret_query_keyword"`
	SecretRegexPattern string `yaml:"secret_regex_pattern"`
}

func Init() (*Config, error) {
	configRaw, err := os.ReadFile("files/config/config.yaml")
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(configRaw, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
