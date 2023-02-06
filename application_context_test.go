package main

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShouldReplaceEnv(t *testing.T) {
	input := "some data with env::MY_KEY2 to be replaced"
	os.Setenv("MY_KEY2", "special_value")
	output, err := replaceEnvVar(input)

	assert.Nil(t, err)
	assert.Equal(t, "some data with special_value to be replaced", output)

}

func TestShouldNotReplaceEnvAndError(t *testing.T) {
	input := "some data with env::MY_KEY_NOT_FOUND to be replaced"
	_, err := replaceEnvVar(input)

	assert.NotNil(t, err)
}

func TestShouldReplaceDateTime(t *testing.T) {
	date, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z")
	input := "some data with sys::now to be replaced"
	output := replaceSystemVar(input, date)

	expected := fmt.Sprintf("some data with 010220061504 to be replaced")
	assert.Equal(t, expected, output)
}
