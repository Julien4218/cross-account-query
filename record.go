package main

import (
	"fmt"
	"log"
	"math"
	"strings"
)

type Record struct {
	result map[string]interface{}
}

func NewRecord(result map[string]interface{}) *Record {
	record := &Record{result: result}
	return record
}

func (r *Record) Merge(fetched map[string]interface{}) {
	for key, value := range fetched {
		_, isMapContainsKey := r.result[key]
		if !isMapContainsKey {
			r.result[key] = value
		}
	}
}

func (r *Record) String() string {
	output := "["
	for key, _ := range r.result {
		if len(output) > 1 {
			output += ","
		}
		output += key
	}
	output += "]"
	return output
}

func (r *Record) GetField(name string) string {
	for key, _ := range r.result {
		lowerKey := strings.ToLower(key)
		if strings.EqualFold(lowerKey, name) {
			untyped := r.result[key]
			switch untyped.(type) {
			case int:
				return fmt.Sprintf("%d", untyped.(int))
			case float64:
				r2 := math.Round(r.result[key].(float64))
				return fmt.Sprintf("%d", int(r2))
			case string:
				return untyped.(string)
			default:
				log.Fatal(fmt.Sprintf("Unknown cast type of %v", untyped))
			}
		}
	}
	return ""
}

func arrayContains(x string, array []string) bool {
	if len(array) > 0 {
		for _, item := range array {
			if strings.EqualFold(x, item) {
				return true
			}
		}
	}
	return false
}
