// Package config - реализует логику создания конфигурационного файла.
package config

import (
	"encoding/json"
	"flag"
	"os"
)

type Config struct {
	ServerAddress string `json:"server_address" env:"SERVER_ADDRESS"`
	Mode          string `json:"server_mode" env:"SERVER_MODE"`
	Config        string `env:"CONFIG"`
}

// NewConfig - создает config.
func NewConfig() *Config {
	conf := &Config{
		ServerAddress: "",
		Mode:          "",
		Config:        "",
	}
	conf.parseFlags()

	conf.parseEnv()

	conf.readConfig()

	return conf
}

// ParseFlags - парсит значения флагов из аргументов командной строки.
func (c *Config) parseFlags() {
	flag.StringVar(&c.ServerAddress, "a", c.ServerAddress, "SERVER_ADDRESS")
	flag.StringVar(&c.Mode, "m", c.Mode, "SERVER_MODE")
	flag.StringVar(&c.Config, "c", c.Config, "CONFIG")
	flag.Parse()
}

// ParseEnv - парсит значения переменных окружения.
func (c *Config) parseEnv() {
	if c.ServerAddress == "" {
		c.ServerAddress = os.Getenv("SERVER_ADDRESS")
	}
	if c.Mode == "" {
		c.Mode = os.Getenv("SERVER_MODE")
	}
	if c.Config == "" {
		c.Config = os.Getenv("CONFIG")
	}
}

// ReadConfig - читает значения из конфигурационного файла если тот задан.
func (c *Config) readConfig() {
	configDataJSON, err := os.ReadFile("../pkg/" + c.Config)
	if err != nil {
		return
	}
	var configJSON Config
	if err = json.Unmarshal(configDataJSON, &configJSON); err != nil {
		return
	}
	if c.ServerAddress == "" {
		c.ServerAddress = configJSON.ServerAddress
	}
	if c.Mode == "" {
		c.Mode = configJSON.Mode
	}
}
