package main

import (
	"fmt"
	"log"
	"os"
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

	err = NewOrchestrator(appContext).Execute()
	if err != nil {
		log.Fatal(err)
		os.Exit(2)
	}

	fmt.Println(appContext.Report)

	fmt.Println("Done")
}

func displaySyntax() {
	fmt.Println("cross-account-query program executes multiple account query and display results.")
	fmt.Println("syntax: cross-account-query config.yml")
}
