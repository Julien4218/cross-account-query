package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Report struct {
	headers map[int]string
	rows    []*Record
}

func NewReport() *Report {
	r := &Report{
		headers: map[int]string{},
		rows:    []*Record{},
	}
	return r
}

func (r *Report) Json() string {
	output := ""
	for _, row := range r.rows {
		rowOutput := ""
		for i := 0; i < len(r.headers); i++ {
			header := r.headers[i]
			colName := getHeaderColumnName(header)
			aliasName := getHeaderColumnAlias(header)
			value := row.GetField(colName)
			if value != "" {
				if rowOutput != "" {
					rowOutput += ","
				}
				rowOutput += "\""
				if aliasName != "" {
					rowOutput += aliasName
				} else {
					rowOutput += colName
				}
				rowOutput += "\":"
				isNative := isValueNative(value)
				if !isNative {
					rowOutput += "\""
				}
				if colName == "timestamp" && aliasName == "" {
					seconds, _ := strconv.ParseInt(value, 10, 64)
					unix := time.UnixMilli(seconds)
					value = unix.Format(time.RFC3339)
				}
				rowOutput += value
				if !isNative {
					rowOutput += "\""
				}
			}
		}
		if rowOutput != "" {
			if output != "" {
				output += ",\n"
			}
			output += "{" + rowOutput + "}"
		}
	}
	if output != "" {
		output = "[\n" + output + "\n]\n"
	}
	return output
}

func isValueNative(value string) bool {
	_, err := strconv.ParseInt(value, 10, 64)
	if err == nil {
		return true
	}
	_, err = strconv.ParseBool(value)
	if err == nil {
		return true
	}
	return false
}

func (r *Report) String() string {
	output := ""
	if len(r.headers) == 0 {
		return output
	}
	for i := 0; i < len(r.headers); i++ {
		if i > 0 {
			output += ","
		}
		header := r.headers[i]
		colName := getHeaderColumnName(header)
		aliasName := getHeaderColumnAlias(header)
		if aliasName != "" {
			output += aliasName
		} else {
			output += colName
		}
	}
	output += "\n"
	for _, row := range r.rows {
		for i := 0; i < len(r.headers); i++ {
			if i > 0 {
				output += ","
			}
			header := r.headers[i]
			colName := getHeaderColumnName(header)
			aliasName := getHeaderColumnAlias(header)
			value := row.GetField(colName)
			if colName == "timestamp" && aliasName == "" {
				seconds, _ := strconv.ParseInt(value, 10, 64)
				unix := time.UnixMilli(seconds)
				value = unix.Format(time.RFC3339)
			}
			if value != "" {
				output += value
			}
		}
		output += "\n"
	}
	return output
}

func (r *Report) AddHeader(name string) {
	for _, header := range r.headers {
		if strings.EqualFold(name, header) {
			return
		}
	}
	log.Debugf(fmt.Sprintf("Adding header %s", name))
	r.headers[len(r.headers)] = name
}

func (r *Report) LastRow() *Record {
	if len(r.rows) > 0 {
		return r.rows[len(r.rows)-1]
	}
	return nil
}

func (r *Report) FindMatchingAll(keys []*MatchKey, record *Record) *Record {
	if len(keys) == 0 || r.rows == nil {
		return nil
	}
	for _, row := range r.rows {
		found := false
		for _, key := range keys {
			rowValue := row.GetField(key.KeyName)
			recordValue := record.GetField(key.KeyName)
			if rowValue != "" && recordValue != "" && rowValue == recordValue {
				found = true
				continue
			} else {
				found = false
				break
			}
		}
		if found {
			return row
		}
	}
	return nil
}

func (r *Report) AddRow(record *Record) {
	r.rows = append(r.rows, record)

}

func getHeaderColumnName(header string) string {
	parts := strings.Split(header, ":")
	if len(parts) > 0 {
		return parts[0]
	}
	return header
}

func getHeaderColumnAlias(header string) string {
	parts := strings.Split(header, ":")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}
