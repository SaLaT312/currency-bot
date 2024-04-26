package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	TelegramToken string `yaml:"telegram_token"`
	Redis         struct {
		RedisAddr string `yaml:"addr"`
		RedisPass string `yaml:"pass"`
		RedisDB   int    `yaml:"db"`
	} `yaml:"redis"`
}

func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}
	//var testConfig interface{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	err = d.Decode(&config)

	if err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}

func main() {

	cfgPath, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	cfg, err := NewConfig(cfgPath)
	if err != nil {
		///log.Println("1111")
		log.Fatal(err)
	}

	RedisAddr = cfg.Redis.RedisAddr
	RedisPass = cfg.Redis.RedisPass
	RedisDB = cfg.Redis.RedisDB
	TelegramApiToken := cfg.TelegramToken

	RedisCheck() // Проверка доступности редиса перед запуском

	GetDataFromBank() // Заполнить редис при старте

	CronTabFunc()
	TelegramBot(TelegramApiToken)

}
