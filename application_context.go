package main

import (
	"errors"
	"fmt"
	"os"

	yaml "gopkg.in/yaml.v2"
)

type ApplicationContext struct {
	Report *Report
	Config *Config
}

func Init(configFilename string) (*ApplicationContext, error) {
	config, err := readFile(configFilename)
	if err != nil {
		return nil, err
	}
	return &ApplicationContext{
		Report: NewReport(),
		Config: config,
	}, nil
}

func readFile(filename string) (*Config, error) {
	var config *Config = nil

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("couldn't read file %s detail:%v", filename, err))
	}

	err2 := yaml.Unmarshal([]byte(data), &config)
	if err2 != nil {
		return nil, errors.New(fmt.Sprintf("couldn't unmarshal yaml %s detail:%v", filename, err2))
	}

	return config, nil
}
