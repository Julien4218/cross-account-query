package main

import "github.com/newrelic/newrelic-client-go/newrelic"

type ClientFactory struct {
	clients map[string]*newrelic.NewRelic
}

func NewClientFactory() *ClientFactory {
	return &ClientFactory{
		clients: make(map[string]*newrelic.NewRelic),
	}
}

func (f *ClientFactory) GetOrCreate(config *ConfigQuery) (*newrelic.NewRelic, error) {
	c, exist := f.clients[config.ApiKey]
	if exist {
		return c, nil
	}
	c, err := NewClient(config)
	if err == nil {
		f.clients[config.ApiKey] = c
		return c, nil
	}
	return nil, err
}
