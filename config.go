package main

import (
	"errors"
	"github.com/kelseyhightower/envconfig"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type Config struct {
	SoundCloud `yaml:"soundcloud"`
	Telegram   `yaml:"telegram"`
}

type SoundCloud struct {
	Token string `yaml:"token" envconfig:"SOUNDCLOUD_TOKEN"`
}

type Telegram struct {
	Token string `yaml:"token" envconfig:"TELEGRAM_TOKEN"`
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
	return cfg
}

func readFile(fileName string, cfg *Config) (e error) {
	f, err := os.Open(fileName)
	if err != nil {
		return errors.New("can't read config file")
	}
	defer f.Close()

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
