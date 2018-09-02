package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

// Config is the users configuration
var Config *Configuration

// Configuration holds the configuration the player sets up when he boots the game on his machine
type Configuration struct {
	Port      string    `yaml:"port"`
	Questions Questions `yaml:"questions"`
	Storage   Storage   `yaml:"storage"`
}

// Questions todo:
type Questions struct {
	UseAWS bool   `yaml:"use-aws"`
	Table  string `yaml:"table"`
	Region string `yaml:"region"`
}

// Storage TODO:
type Storage struct {
	SaveSubmissions bool `yaml:"save-submissions"`
	SaveLogs        bool `yaml:"save-logs"`
}

// ParseConfiguration uses the config.yml file to make an object to use
func ParseConfiguration() (*Configuration, error) {
	var c Configuration
	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(file, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// SetupConfiguration sets the global configuration using the ParseConfiguration method
func SetupConfiguration(c *Configuration, err error) *Configuration {

	if err != nil {
		logger.Printf("error setting up configuration: %v\n", err)
		panic(err)
	}
	return c
}
