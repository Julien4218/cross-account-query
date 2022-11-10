package main

import (
	"os"
	"testing"

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
