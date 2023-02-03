package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReportShouldCreateEmpty(t *testing.T) {
	r := NewReport()
	assert.NotNil(t, r)
	assert.Equal(t, len(r.headers), 0)
	assert.Equal(t, len(r.rows), 0)
	assert.Equal(t, r.String(), "")
}

func TestReportShouldAddHeader(t *testing.T) {
	r := NewReport()
	r.AddHeader("field1")
	assert.Equal(t, len(r.headers), 1)
	assert.Contains(t, r.String(), "field1")
}

func TestReportShouldNotAddDuplicateHeader(t *testing.T) {
	r := NewReport()
	r.AddHeader("field1")
	r.AddHeader("field1")
	r.AddHeader("field1")
	assert.Equal(t, len(r.headers), 1)
	assert.Contains(t, r.String(), "field1")
}

func TestReportShouldAddMultipleHeaders(t *testing.T) {
	r := NewReport()
	r.AddHeader("field1")
	r.AddHeader("field2")
	r.AddHeader("field3")
	assert.Equal(t, len(r.headers), 3)
	assert.Contains(t, r.String(), "field1,field2,field3")
}

func TestReportShouldAddRow(t *testing.T) {
	r := NewReport()
	r.AddHeader("field1")
	r.AddRow(NewRecord(NewDataMapBuilder().WithField("field1", 123).Build(), nil))
	assert.Equal(t, len(r.headers), 1)
	assert.Equal(t, len(r.rows), 1)
	assert.Contains(t, r.String(), "field1")
	assert.Contains(t, r.String(), "123")
}

func TestReportShouldAddRows(t *testing.T) {
	r := NewReport()
	r.AddHeader("field1")
	r.AddRow(NewRecord(NewDataMapBuilder().WithField("field1", 123).Build(), nil))
	r.AddRow(NewRecord(NewDataMapBuilder().WithField("field1", 231).Build(), nil))
	r.AddRow(NewRecord(NewDataMapBuilder().WithField("field1", 312).Build(), nil))
	assert.Equal(t, len(r.headers), 1)
	assert.Equal(t, len(r.rows), 3)
	assert.Contains(t, r.String(), "field1")
	assert.Contains(t, r.String(), "123\n")
	assert.Contains(t, r.String(), "231\n")
	assert.Contains(t, r.String(), "312\n")
}

func TestReportShouldAddMultipleHeadersRows(t *testing.T) {
	r := NewReport()
	r.AddHeader("field1")
	r.AddHeader("field2")
	r.AddHeader("field3")
	r.AddRow(NewRecord(NewDataMapBuilder().WithField("field1", 123).WithField("field2", 456).WithField("field3", 789).Build(), nil))
	r.AddRow(NewRecord(NewDataMapBuilder().WithField("field1", 231).WithField("field2", 564).WithField("field3", 897).Build(), nil))
	r.AddRow(NewRecord(NewDataMapBuilder().WithField("field1", 312).WithField("field2", 645).WithField("field3", 978).Build(), nil))
	assert.Equal(t, len(r.headers), 3)
	assert.Equal(t, len(r.rows), 3)
	assert.Contains(t, r.String(), "field1,field2,field")
	assert.Contains(t, r.String(), "123,456,789\n")
	assert.Contains(t, r.String(), "231,564,897\n")
	assert.Contains(t, r.String(), "312,645,978\n")
}

func TestReportShouldOutputJson(t *testing.T) {
	r := NewReport()
	r.AddHeader("field1")
	r.AddHeader("field2")
	r.AddHeader("field3")
	r.AddRow(NewRecord(NewDataMapBuilder().WithField("field1", 123).WithField("field2", 456).WithField("field3", 789).Build(), nil))
	r.AddRow(NewRecord(NewDataMapBuilder().WithField("field1", 231).WithField("field2", 564).WithField("field3", 897).Build(), nil))
	r.AddRow(NewRecord(NewDataMapBuilder().WithField("field1", 312).WithField("field2", 645).WithField("field3", 978).Build(), nil))
	assert.Contains(t, r.Json(), "\"field1\":123")
}

func TestReportShouldOutputAliasHeader(t *testing.T) {
	r := NewReport()
	r.AddHeader("field1:My Field 1")
	assert.Equal(t, len(r.headers), 1)
	assert.Contains(t, r.String(), "My Field 1")
}

func TestShouldHaveNoHeader(t *testing.T) {
	r := NewReport()
	assert.Equal(t, len(r.headers), 0)
	assert.Equal(t, r.String(), "")
}

func TestShouldNotFetchMissingRowData(t *testing.T) {
	r := NewReport()
	r.AddHeader("field1")
	r.AddHeader("field2")
	r.AddRow(NewRecord(NewDataMapBuilder().WithField("field3", 789).Build(), nil))
	assert.Equal(t, len(r.headers), 2)
	assert.Equal(t, r.String(), "field1,field2\n,\n")
}
