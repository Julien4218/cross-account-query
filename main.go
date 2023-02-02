package main

import (
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	if len(os.Args) < 2 {
		displaySyntax()
		os.Exit(1)
	}
	isJsonOutput := false
	for _, arg := range os.Args {
		if arg == "--json" {
			isJsonOutput = true
		}
		if arg == "--debug" {
			log.SetLevel(log.DebugLevel)
		}
		if arg == "--trace" {
			log.SetLevel(log.TraceLevel)
		}
		if arg == "--error" {
			log.SetLevel(log.ErrorLevel)
		}
		if arg == "--warn" {
			log.SetLevel(log.WarnLevel)
		}
	}
	fileName := os.Args[1]
	file, err := os.Stat(fileName)
	if err != nil {
		log.Fatalf("config file received %s does not exist, detail:%s", fileName, err)
	}
	if file.IsDir() {
		log.Fatalf("config file received %s is a directory, not a file", fileName)
	}

	log.Debug(fmt.Sprintf("cross-account-query program. Got config fileName:%s", fileName))

	appContext, err := Init(fileName)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	err = NewOrchestrator(appContext).Execute()
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	if isJsonOutput {
		fmt.Printf(appContext.Report.Json())
	} else {
		fmt.Print(appContext.Report.String())
		log.Debug("Done")
	}
}

func displaySyntax() {
	log.Debug("cross-account-query program executes multiple account query and display results.")
	log.Debug("syntax: cross-account-query config.yml")
}
