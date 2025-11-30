package config

import (
	"os"
	"sync"

	"gopkg.in/yaml.v2"
)

type Config struct {
	App      AppConfig      `yaml:"app"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	Kafka    KafkaConfig    `yaml:"kafka"`
}

type AppConfig struct {
	Port string `yaml:"port"`
}

type DatabaseConfig struct {
	PostgresURL string `yaml:"postgres_url"`
}

type RedisConfig struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
}

type KafkaConfig struct {
	BrokerAddress string `yaml:"broker_address"`
	Topic         string `yaml:"topic"`
	GroupId       string `yaml:"group_id"`
}

var (
	once     sync.Once
	instance *Config
)

func Load(configPath string) (*Config, error) {
	var cfg Config
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}

func GetConfig(configPath string) *Config {
	once.Do(func() {
		var err error
		instance, err = Load(configPath)
		if err != nil {
			panic(err)
		}
	})
	return instance
}
