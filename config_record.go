package main

import (
	"fmt"
)

type Config struct {
	Base    *ConfigQuery   `yaml:"base"`
	Columns []*ConfigQuery `yaml:"columns"`
}

type ConfigQuery struct {
	ApiKey       string   `yaml:"api_key"`
	AccountId    int      `yaml:"account_id"`
	Region       string   `yaml:"region"`
	Query        string   `yaml:"query"`
	SelectFields []string `yaml:"select_fields"`
	CanBatch     bool     `yaml:"can_batch"`
}

type ColumnQuery struct {
	ApiKey       string   `yaml:"api_key"`
	AccountId    int      `yaml:"account_id"`
	Region       string   `yaml:"region"`
	Query        string   `yaml:"query"`
	SelectFields []string `yaml:"select_fields"`
}

func (r *ConfigQuery) String() string {
	return fmt.Sprintf("\nConfigQuery\n\taccountID:%d\n\tregion:%s\n\tapiKey:%s\n\tquery:%s\n\tselectFields:%v", r.AccountId, r.Region, r.GetSafeReadableApiKey(), r.Query, r.SelectFields)
}

func (b *Config) String() string {
	output := b.Base.String()
	for _, column := range b.Columns {
		output += fmt.Sprintf("\nColumnQuery\n\taccountID:%d\n\tregion:%s\n\tapiKey:%s\n\tquery:%s\n\tselectFields:%v", column.AccountId, column.Region, column.GetSafeReadableApiKey(), column.Query, column.SelectFields)
	}
	return output
}

func (r *ConfigQuery) GetSafeReadableApiKey() string {
	return getSecureString(r.ApiKey)
}

func (b *ColumnQuery) GetSafeReadableApiKey() string {
	return getSecureString(b.ApiKey)
}

func getSecureString(value string) string {
	if len(value) == 0 {
		return ""
	}
	if len(value) < 8 {
		return "********"
	}
	return value[0:8] + "******"
}
