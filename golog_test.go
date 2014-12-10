package golog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLogger(t *testing.T) {
	defer cleanupTest()

	logger := GetLogger("logger")
	assert.NotNil(t, logger)

	another := GetLogger("another")
	assert.NotEqual(t, logger, another)

	_logger := GetLogger("logger")
	assert.Equal(t, logger, _logger)

	_another := GetLogger("another")
	assert.Equal(t, another, _another)
}
