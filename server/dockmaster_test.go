package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetContainer(t *testing.T) {
	t.SkipNow()
}

func TestSaveContainer(t *testing.T) {
	t.SkipNow()
}

func TestGetConfiguration(t *testing.T) {
	c := getConfiguration()

	assert.NotEmpty(t, c)
}
