package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type config struct {
	System system `yaml:"system"`
}
type system struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
	Env  string `yaml:"env"`
}

var Config *config

func init() {
	yamlFile, err := ioutil.ReadFile("./config.yaml")
	if err != nil {
		return
	}
	yaml.Unmarshal(yamlFile, &Config)
}
