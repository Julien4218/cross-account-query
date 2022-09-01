package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

type Record struct {
	result map[string]interface{}
}

func NewRecord(result map[string]interface{}) *Record {
	if result == nil {
		result = make(map[string]interface{})
	}
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
	for key := range r.result {
		if len(output) > 1 {
			output += ","
		}
		output += key
	}
	output += "]"
	return output
}

func (r *Record) GetField(name string) string {
	value, _ := r.GetFieldOrError(name)
	return value
}

func (r *Record) GetFieldOrError(name string) (string, error) {
	for key := range r.result {
		lowerKey := strings.ToLower(key)
		if strings.EqualFold(lowerKey, name) {
			untyped := r.result[key]
			switch untyped.(type) {
			case int:
				return fmt.Sprintf("%d", untyped.(int)), nil
			case float64:
				r2 := math.Round(r.result[key].(float64))
				return fmt.Sprintf("%d", int(r2)), nil
			case string:
				return untyped.(string), nil
			default:
				return "", errors.New(fmt.Sprintf("Unknown cast type of %v", untyped))
			}
		}
	}
	return "", nil
}
