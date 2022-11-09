package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

type Record struct {
	result     map[string]interface{}
	facetNames []string
}

func NewRecord(result map[string]interface{}, facetNames []string) *Record {
	if result == nil {
		result = make(map[string]interface{})
	}
	record := &Record{result: result, facetNames: facetNames}
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
		if key == "facet" {
			continue
		}
		if len(output) > 1 {
			output += ","
		}
		// DEBUG
		output += fmt.Sprintf("%s:%s", key, r.GetField(key))
	}
	for _, key := range r.facetNames {
		if len(output) > 1 {
			output += ","
		}
		// DEBUG
		output += fmt.Sprintf("%s:%s", key, r.GetField(key))
	}
	output += "]\n"
	return output
}

func (r *Record) GetField(name string) string {
	value, _ := r.GetFieldOrError(name)
	return value
}

func (r *Record) GetFieldOrError(name string) (string, error) {
	for key := range r.result {
		if name == "facet" {
			return "", nil
		}
		if len(r.facetNames) > 1 {
			if ok, index := arrayContains(r.facetNames, name); ok {
				facets := r.result["facet"]
				for i, j := range facets.([]interface{}) {
					if index == i {
						return castAsString(j)
					}
				}
			}
		}

		if strings.EqualFold(key, name) {
			return castAsString(r.result[key])
		}
	}
	return "", nil
}

func castAsString(value interface{}) (string, error) {
	switch value.(type) {
	case int:
		return fmt.Sprintf("%d", value.(int)), nil
	case float64:
		r2 := math.Round(value.(float64))
		return fmt.Sprintf("%d", int(r2)), nil
	case string:
		return value.(string), nil
	default:
		return "", errors.New(fmt.Sprintf("Unknown cast type of %v", value))
	}
}

func arrayContains(s []string, str string) (bool, int) {
	for i, v := range s {
		if v == str {
			return true, i
		}
	}

	return false, 0
}
