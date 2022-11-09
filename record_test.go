package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecordShouldCreateWhenUndefined(t *testing.T) {
	r := NewRecord(nil, nil)
	assert.NotNil(t, r)
}

func TestRecordShouldCreateEmptyWhenUndefined(t *testing.T) {
	r := NewRecord(nil, nil)
	assert.Equal(t, len(r.result), 0)
}

func TestRecordShouldGetIntField(t *testing.T) {
	result := NewDataMapBuilder().WithField("field1", 456).Build()
	r := NewRecord(result, nil)
	assert.Equal(t, len(r.result), 1)
	require.Equal(t, r.GetField("field1"), "456")
}

func TestRecordShouldGetMultipleFields(t *testing.T) {
	result := NewDataMapBuilder().WithField("field1", 456).WithField("field2", 678).Build()
	r := NewRecord(result, nil)
	assert.Equal(t, len(r.result), 2)
	require.Equal(t, r.GetField("field1"), "456")
	require.Equal(t, r.GetField("field2"), "678")
}

func TestRecordShouldNotGetField(t *testing.T) {
	result := NewDataMapBuilder().Build()
	r := NewRecord(result, nil)
	assert.Equal(t, len(r.result), 0)
	require.Equal(t, r.GetField("notexisting"), "")
}

func TestRecordShouldGetStringField(t *testing.T) {
	result := NewDataMapBuilder().WithField("fieldX", "something of value").Build()
	r := NewRecord(result, nil)
	assert.Equal(t, len(r.result), 1)
	require.Equal(t, r.GetField("fieldX"), "something of value")
}

func TestRecordShouldGetFieldCaseInsensitiveKey(t *testing.T) {
	result := NewDataMapBuilder().WithField("fieldXYZ", "something of HIGH value").Build()
	r := NewRecord(result, nil)
	assert.Equal(t, len(r.result), 1)
	require.Equal(t, r.GetField("fieldxyz"), "something of HIGH value")
}

func TestRecordShouldNotGetFieldArray(t *testing.T) {
	result := NewDataMapBuilder().WithField("fieldZ", []string{"unsupported arrays"}).Build()
	r := NewRecord(result, nil)
	assert.Equal(t, len(r.result), 1)
	_, err := r.GetFieldOrError("fieldZ")
	assert.Error(t, err)
}

func TestRecordShouldStringifyKeys(t *testing.T) {
	result := NewDataMapBuilder().WithField("field1", 456).WithField("field2", 678).Build()
	r := NewRecord(result, nil)
	assert.Equal(t, r.String(), "[field1:456,field2:678]\n")
}

func TestRecordShouldMerge(t *testing.T) {
	data1 := NewDataMapBuilder().WithField("field1", 456).Build()
	data2 := NewDataMapBuilder().WithField("field2", 678).Build()
	r := NewRecord(data1, nil)
	r.Merge(data2)
	assert.Equal(t, len(r.result), 2)
	require.Equal(t, r.GetField("field1"), "456")
	require.Equal(t, r.GetField("field2"), "678")
}
