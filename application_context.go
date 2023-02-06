package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

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

	content := string(data)
	replaced, err := replaceEnvVar(content)
	if err != nil {
		return nil, err
	}

	replaced = replaceSystemVar(replaced, time.Now())

	err2 := yaml.Unmarshal([]byte(replaced), &config)
	if err2 != nil {
		return nil, errors.New(fmt.Sprintf("couldn't unmarshal yaml %s detail:%v", filename, err2))
	}

	return config, nil
}

func replaceEnvVar(value string) (string, error) {
	if len(value) == 0 {
		return "", nil
	}
	result := value
	re := regexp.MustCompile(`env::(\w+)`)
	for found := re.Find([]byte(result)); found != nil; found = re.Find([]byte(result)) {
		match := string(found)
		key := strings.ReplaceAll(match, "env::", "")
		replaced := os.Getenv(key)
		if replaced == "" {
			return "", errors.New(fmt.Sprintf("Environment variable %s is defined but not found, cannot replace", match))
		}
		result = strings.ReplaceAll(result, match, replaced)
	}
	return result, nil
}

func replaceSystemVar(value string, date time.Time) string {
	if len(value) == 0 {
		return ""
	}
	result := value
	re := regexp.MustCompile(`sys::(\w+)`)
	for found := re.Find([]byte(result)); found != nil; found = re.Find([]byte(result)) {
		match := string(found)
		key := strings.ReplaceAll(match, "sys::", "")
		replaced := ""
		switch key {
		case "now":
			replaced = fmt.Sprintf("%02d%02d%04d%02d%02d", date.Month(), date.Day(), date.Year(), date.Hour(), date.Minute())
		}
		result = strings.ReplaceAll(result, match, replaced)
	}
	return result
}
