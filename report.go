package main

import (
	"strconv"
	"strings"
	"time"
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
	r.headers[len(r.headers)] = name
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
