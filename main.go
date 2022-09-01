package main

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/newrelic/newrelic-client-go/pkg/nrdb"
)

func main() {
	if len(os.Args) != 2 {
		displaySyntax()
		os.Exit(1)
	}
	fileName := os.Args[1]
	file, err := os.Stat(fileName)
	if err != nil {
		log.Fatalf("config file received %s does not exist, detail:%s", fileName, err)
	}
	if file.IsDir() {
		log.Fatalf("config file received %s is a directory, not a file", fileName)
	}

	fmt.Println(fmt.Sprintf("cross-account-query program. Got config fileName:%s", fileName))

	appContext, err := Init(fileName)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	config := appContext.Config
	report := appContext.Report

	for _, selectField := range config.Base.SelectFields {
		report.AddHeader(selectField)
	}

	client, err := NewClient(config.Base)
	if err != nil {
		log.Fatal(fmt.Sprintf("Couldn't load newrelic client, detail:%s", err))
	}
	query := config.Base.Query
	result, err := client.Nrdb.Query(config.Base.AccountId, nrdb.NRQL(query))
	if err != nil {
		log.Fatal(fmt.Sprintf("Couldn't execute query, detail:%s", err))
	}
	maxResult := len(result.Results)
	fmt.Println(fmt.Sprintf("Results(%d):", maxResult))
	for index, result := range result.Results {
		fmt.Println(fmt.Sprintf("Processing row %d/%d:", index+1, maxResult))
		record := NewRecord(result)
		report.AddRow(record)

		for _, configColumn := range config.Columns {
			for _, selectField := range configColumn.SelectFields {
				report.AddHeader(selectField)
			}

			client, err := NewClient(configColumn)
			if err != nil {
				log.Fatal(fmt.Sprintf("Couldn't load newrelic client, detail:%s", err))
			}
			query := replaceQueryFields(configColumn.Query, record)
			if query != "" {
				attempt := 1
				maxAttempt := 10
				for attempt < maxAttempt {
					subResult, err := client.Nrdb.Query(configColumn.AccountId, nrdb.NRQL(query))
					if err == nil {
						if len(subResult.Results) > 0 {
							record.Merge(subResult.Results[0])
						}
						break
					}
					if attempt == maxAttempt {
						log.Fatal(fmt.Sprintf("Couldn't execute query, detail:%s", err))
					}
					attempt++
					log.Println(fmt.Sprintf("Error while executing, retrying attempt %d, detail:%s", attempt, err))
					time.Sleep(2 * time.Second)
				}
			}
		}
	}

	fmt.Println(report)

	fmt.Println("Done")
}

func displaySyntax() {
	fmt.Println("cross-account-query program executes multiple account query and display results.")
	fmt.Println("syntax: cross-account-query config.yml")
}

func replaceQueryFields(query string, record *Record) string {
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
