package main

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

type Orchestrator struct {
	appContext    *ApplicationContext
	clientFactory *ClientFactory
}

func NewOrchestrator(appContext *ApplicationContext) *Orchestrator {
	return newOrchestratorInternal(appContext, NewClientFactory())
}

func newOrchestratorInternal(appContext *ApplicationContext, clientFactory *ClientFactory) *Orchestrator {
	return &Orchestrator{
		appContext:    appContext,
		clientFactory: clientFactory,
	}
}

func (o *Orchestrator) Execute() error {
	config := o.appContext.Config
	report := o.appContext.Report

	o.addAllHeaders()

	records, err := o.query(o.appContext.Config.Base, 0, nil)
	if err != nil {
		return err
	}
	maxResult := len(records)
	fmt.Println(fmt.Sprintf("Results(%d):", maxResult))
	for index, record := range records {
		fmt.Println(fmt.Sprintf("Processing row %d/%d:", index+1, maxResult))
		report.AddRow(record)

		for _, configColumn := range config.Columns {
			columnRecords, err := o.query(configColumn, maxResult, record)
			if err != nil {
				return err
			}
			if len(columnRecords) > 0 {
				record.Merge(columnRecords[0].result)
			}
		}
	}
	return nil
}

func (o *Orchestrator) addAllHeaders() {
	for _, selectField := range o.appContext.Config.Base.SelectFields {
		o.appContext.Report.AddHeader(selectField)
	}
	for _, configColumn := range o.appContext.Config.Columns {
		for _, selectField := range configColumn.SelectFields {
			o.appContext.Report.AddHeader(selectField)
		}
	}
}

func (o *Orchestrator) query(config *ConfigQuery, maxResult int, record *Record) ([]*Record, error) {
	client, err := o.clientFactory.GetOrCreate(config)
	if err != nil {
		return []*Record{}, errors.New(fmt.Sprintf("Couldn't load newrelic client, detail:%s", err))
	}

	records := []*Record{}
	query := config.Query
	if query != "" {
		if record != nil {
			query = o.replaceQueryFields(config.Query, record)
		}
		attempt := 1
		maxAttempt := 10
		for attempt < maxAttempt {
			results, err := client.Nrdb.Query(config.AccountId, nrdb.NRQL(query))
			if err == nil {
				if len(results.Results) > 0 {
					for _, result := range results.Results {
						records = append(records, NewRecord(result))
					}
				}
				break
			}
			if attempt == maxAttempt {
				return records, errors.New(fmt.Sprintf("Couldn't execute query, detail:%s", err))
			}
			attempt++
			log.Println(fmt.Sprintf("Error while executing, retrying attempt %d, detail:%s", attempt, err))
			time.Sleep(2 * time.Second)
		}
	}

	return records, nil
}

func (o *Orchestrator) replaceQueryFields(query string, record *Record) string {
	output := query
	re := regexp.MustCompile(`env::(\w+)`)
	for _, item := range re.FindAll([]byte(query), -1) {
		key := strings.ReplaceAll(string(item), "env::", "")
		value := record.GetField(key)
		if value == "" {
			return ""
		}
		output = strings.ReplaceAll(output, string(item), value)
	}

	return output
}
