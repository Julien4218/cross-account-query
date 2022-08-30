package main

import (
	"errors"
	"fmt"

	"github.com/newrelic/newrelic-client-go/newrelic"
)

var (
	NRClient    *newrelic.NewRelic
	serviceName = "newrelic-cli"
)

// NewClient initializes the New Relic client.
func NewClient(config *ConfigQuery) (*newrelic.NewRelic, error) {

	if config.ApiKey == "" {
		return nil, errors.New("a User API key is required")
	}

	region := config.Region
	userAgent := fmt.Sprintf("cross-account-query CLI")

	cfgOpts := []newrelic.ConfigOption{
		newrelic.ConfigPersonalAPIKey(config.ApiKey),
		newrelic.ConfigRegion(region),
		newrelic.ConfigUserAgent(userAgent),
		newrelic.ConfigServiceName(serviceName),
	}

	nrClient, err := newrelic.New(cfgOpts...)
	if err != nil {
		return nil, fmt.Errorf("unable to create New Relic client with error: %s", err)
	}

	return nrClient, nil
}
