package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type EndpointConfig struct {
	Name string `yaml:"name"`
	Handler string  `yaml:"handler"`
}

type KafkaConfig struct {
	Brokers string `yaml:"brokers"`
}
type RestConfig struct {
	ApiCfg map[string]EndpointConfig `yaml:"api"`
	KafkaCfg KafkaConfig `yaml:"kafka"`
}

func (c *RestConfig) ReadConfig(f string) *RestConfig {
	yamlFile, err := ioutil.ReadFile(f)
	if err != nil {
		log.Printf("Config file get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return c
}
