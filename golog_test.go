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

func TestEnableLogger(t *testing.T) {
	defer cleanupTest()

	ta := &testAppender{}
	logger := GetLogger("my/logger")
	logger.Enable(ta)

	Disable("my/logger")
	logger.Debug("some msg")
	logger.Info("some msg")
	logger.Warn("some msg")
	logger.Error("some msg")

	assert.Exactly(t, 0, ta.count)

	Enable("my/logger")
	logger.Debug("some msg")
	logger.Info("some msg")
	logger.Warn("some msg")
	logger.Error("some msg")

	assert.Exactly(t, 4, ta.count)

	Enable("some-unknown-name")
}

func TestDisableLogger(t *testing.T) {
	defer cleanupTest()

	ta := &testAppender{}
	logger := GetLogger("my/logger")
	logger.Enable(ta)
	logger.Info("some msg")
	logger.Info("some msg")

	assert.Exactly(t, 2, ta.count)

	Disable("my/logger")
	logger.Debug("some msg")
	logger.Info("some msg")
	logger.Warn("some msg")
	logger.Error("some msg")

	assert.Exactly(t, 2, ta.count)

	Enable("some-unknown-name")
}
