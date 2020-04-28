package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
}

func TestContains(t *testing.T) {
        assert := assert.New(t)
        a := []string{"here", "there", "everywhere"}
        assert.True(contains(a, "there"))
        assert.False(contains(a, "nowhere"))
}

