package main

import (
	"errors"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	Segmentio `yaml:"segmentio"`
	Telegram  `yaml:"telegram"`
	Mode      string `yaml:"mode" envconfig:"MODE"`
}

type Segmentio struct {
	Token string `yaml:"token" envconfig:"SEGMENTIO_TOKEN"`
}

type Telegram struct {
	Token     string `yaml:"token" envconfig:"TELEGRAM_TOKEN"`
	TestToken string `yaml:"testToken" envconfig:"TELEGRAM_TEST_TOKEN"`
}

func loadConfig(file string) (cfg Config) {
	if err := readFile(file, &cfg); err != nil {
		log.Println(err)
	}
	if err := readEnv(&cfg); err != nil {
		log.Println(err)
	}
	if (cfg == Config{}) {
		log.Fatal("config not loaded!")
	}
	if cfg.Mode == "debug" {
		cfg.Telegram.Token = cfg.Telegram.TestToken
	}
	return cfg
}

func readFile(fileName string, cfg *Config) (e error) {
	f, err := os.Open(fileName)
	if err != nil {
		return errors.New("can't read config file")
	}
	defer closeFile(f)

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(cfg)
	if err != nil {
		return errors.New("can't decode config file")
	}
	return nil
}

func readEnv(cfg *Config) (e error) {
	log.Println("Loading env")
	err := envconfig.Process("", cfg)
	if err != nil {
		return errors.New("can't read environment variables")
	}
	return nil
}

func closeFile(f *os.File) {
	err := f.Close()

	if err != nil {
		log.Printf("error while closing file: %v\n", err)
	}
}
