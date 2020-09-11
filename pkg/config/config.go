package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	AppPort   string `yaml:"app_port"`
	MediaRoot string `yaml:"media_dir"`

	VlcHost     string `yaml:"vlc_host"`
	VlcPort     string `yaml:"vlc_port"`
	VlcPassword string `yaml:"vlc_password"`
	VlcLogin    string `yaml:"vlc_login"`
}

func InitConf() *AppConfig {
	path, err := os.Getwd()
	if err != nil {
		log.Fatal("Can't locate current dir", err)
	}

	log.Print("Current app dir: ", path)
	confFile := path + "/configs/app_conf.yml"
	confFile = filepath.Clean(confFile)
	log.Print("Try open config file: ", confFile)

	file, errP := os.Open(confFile)
	if errP != nil {
		log.Fatal("Can't open config file: ", confFile, errP)
	}
	defer file.Close()
	var cfg AppConfig
	decoder := yaml.NewDecoder(file)
	errD := decoder.Decode(&cfg)
	if errD != nil {
		log.Fatal("Invalid config file", errD, confFile)
	}

	return &cfg
}
