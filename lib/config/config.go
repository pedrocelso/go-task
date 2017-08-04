package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"runtime"

	yaml "gopkg.in/yaml.v2"
)

// Config contains all runtime configuration of lumen
type Config struct {
	Mysql  Mysql `yaml:"mysql"`
	Server struct {
		Port int `yaml:"port"`
	} `yaml:"server"`
}

// Mysql contain all the connection data to authenticate on MySql
type Mysql struct {
	Host     string `yaml:"host"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
	Port     int    `yaml:"port"`
}

var cfg *Config

// MustInit initializes the config
func Init() *Config {
	_, currentFilename, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatal(1, "Unable to get current directory for config.yaml\n")
	}
	configPath := fmt.Sprintf("%v/config.yaml", path.Dir(currentFilename))
	var file []byte
	file, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Fatal(1, err)
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		log.Fatal(1, err)
	}

	return cfg
}
