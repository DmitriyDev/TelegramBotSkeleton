package main

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type HandlerConfig struct {
	Name          string `yaml:"name"`
	Command       string `yaml:"command"`
	Description   string `yaml:"description"`
	HandlerMethod string `yaml:"handler"`
	Admin         bool   `yaml:"admin"`
	Visible       bool   `yaml:"visible"`
}

type Config struct {
	Commands []HandlerConfig
}

func GetConfig() Config {
	cfd, err := ioutil.ReadFile(ConfigFilePath)
	if err != nil {
		toLogFatalF("Read config error: %v", err)
	}

	c := Config{}
	err = yaml.Unmarshal([]byte(cfd), &c)
	if err != nil {
		toLogFatalF("Unmarshall config error: %v", err)
	}

	return c
}
