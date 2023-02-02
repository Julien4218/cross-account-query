package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

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

	records, err := o.query(o.appContext.Config.Base, report.rows)
	if err != nil {
		return err
	}
	log.Debug(fmt.Sprintf("Results(%d):", len(records)))
	for rindex, record := range records {
		log.Debug(fmt.Sprintf("Creating row %d/%d:", rindex+1, len(records)))
		report.AddRow(record)
	}

	for cindex, configColumn := range config.Columns {
		log.Debug(fmt.Sprintf("Processing column %d/%d:", cindex+1, len(config.Columns)))
		if configColumn.CanBatch {
			columnRecords, err := o.query(configColumn, report.rows)
			if err != nil {
				return err
			}
			// DEBUG
			// fmt.Printf("columnRecords:%v", columnRecords)
			for _, columnRecord := range columnRecords {
				add := NewRecord(columnRecord.result, columnRecord.facetNames)
				keys := o.getMatchingColumns(configColumn)
				found := report.FindMatchingAll(keys, add)
				if found != nil {
					found.Merge(add.result)
				}
			}
		} else {
			for rindex, record := range report.rows {
				log.Debug(fmt.Sprintf("Processing row %d/%d:", rindex+1, len(report.rows)))
				columnRecords, err := o.query(configColumn, []*Record{record})
				if err != nil {
					return err
				}
				if len(columnRecords) > 0 {
					// DEBUG
					// fmt.Printf("\nMerging %v\n", columnRecords[0].result)
					record.Merge(columnRecords[0].result)
				}
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

func (o *Orchestrator) query(config *ConfigQuery, records []*Record) ([]*Record, error) {
	client, err := o.clientFactory.GetOrCreate(config)
	if err != nil {
		return []*Record{}, errors.New(fmt.Sprintf("Couldn't load newrelic client, detail:%s", err))
	}

	output := []*Record{}
	query := o.getQueryReplacedFields(config, records)
	if query != "" {
		attempt := 1
		maxAttempt := 10
		for attempt < maxAttempt {
			log.Debug(fmt.Sprintf("Executing query:%s", query))
			results, err := client.Nrdb.Query(config.AccountId, nrdb.NRQL(query))
			if err == nil {
				if len(results.Results) > 0 {
					for _, result := range results.Results {
						record := NewRecord(result, results.Metadata.Facets)
						// Facet bug
						// record.Merge(results.OtherResult)
						output = append(output, record)
					}
				}
				break
			}
			if attempt == maxAttempt {
				return output, errors.New(fmt.Sprintf("Couldn't execute query, detail:%s", err))
			}
			attempt++
			log.Println(fmt.Sprintf("Error while executing, retrying attempt %d, detail:%s", attempt, err))
			time.Sleep(2 * time.Second)
		}
	}

	return output, nil
}

func (o *Orchestrator) getQueryReplacedFields(config *ConfigQuery, records []*Record) string {
	query := config.Query
	output := query
	keys := o.getMatchingColumns(config)
	for _, key := range keys {
		value := ""
		for _, record := range records {
			add := record.GetField(key.KeyName)
			if add != "" {
				if value != "" {
					value += ","
				}
				replaced := key.Replace(add)
				value += replaced
			}
		}
		if value == "" {
			return ""
		}
		output = strings.ReplaceAll(output, fmt.Sprintf("%s%s", key.Match, key.KeyName), value)
	}

	return output
}

func (o *Orchestrator) getMatchingColumns(config *ConfigQuery) []*MatchKey {
	keys := []*MatchKey{}
	re := regexp.MustCompile(`data::(\w+)`)
	for _, item := range re.FindAll([]byte(config.Query), -1) {
		key := strings.ReplaceAll(string(item), "data::", "")
		keys = append(keys, NewMatchKey("data::", key))
	}
	re = regexp.MustCompile(`add1UnixDay::(\w+)`)
	for _, item := range re.FindAll([]byte(config.Query), -1) {
		key := strings.ReplaceAll(string(item), "add1UnixDay::", "")
		keys = append(keys, NewMatchKey("add1UnixDay::", key))
	}
	return keys
}
