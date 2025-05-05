package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	AWSConfig   AWSConfig   `yaml:"aws"`
	ImageConfig ImageConfig `yaml:"image"`
}

type AWSConfig struct {
	AccountID string `yaml:"accountID"`
	Bucket    string `yaml:"bucket"`
	Region    string `yaml:"region"`
	AccessKey string `yaml:"accessKey"`
	SecretKey string `yaml:"secretKey"`
}

type ImageConfig struct {
	BaseURL         string `yaml:"baseURL"`
	ThumbnailParams string `yaml:"thumbnailParams"`
	FullsizeParams  string `yaml:"fullsizeParams"`
	Folder          string `yaml:"folder"`
}

func LoadConfig(path string) (*Config, error) {
	var cfg Config

	yamlFile, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
