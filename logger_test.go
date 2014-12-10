package golog

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type testAppender struct {
	count      int
	errorCount int
}

func (s *testAppender) Append(log Log) error {
	s.count += 1

	if log.Level.value >= 30 {
		s.errorCount += 1
	}
	return nil
}

func cleanupTest() {
	Default = &Logger{
		Name:  "default",
		Level: DEBUG,
	}

	Default.Enable(StdoutAppender())
}

func (s *testAppender) Id() string {
	return "github.com/ivpusic/golog/test"
}

func TestDefaultPresent(t *testing.T) {
	assert.NotNil(t, Default, "default logged should be authomatically created")
}

func TestEnable(t *testing.T) {
	defer cleanupTest()
	ta := &testAppender{}

	oldcount := len(Default.appenders)
	Default.Enable(ta)
	assert.True(t, oldcount == len(Default.appenders)-1)

	Default.Info("some msg")
	assert.Exactly(t, 1, ta.count)

	Default.Info("some msg")
	assert.Exactly(t, 2, ta.count)
}

func TestDisableByInstance(t *testing.T) {
	defer cleanupTest()
	ta := &testAppender{}

	oldcount := len(Default.appenders)
	Default.Enable(ta)
	assert.True(t, oldcount == len(Default.appenders)-1)

	Default.Disable(ta)
	assert.True(t, oldcount == len(Default.appenders))

	Default.Info("some msg")
	assert.Exactly(t, 0, ta.count)

	Default.Info("some msg")
	assert.Exactly(t, 0, ta.count)

	Default.Enable(ta)
	Default.Info("some msg")
	assert.Exactly(t, 1, ta.count)
}

func TestDisableById(t *testing.T) {
	defer cleanupTest()
	ta := &testAppender{}

	oldcount := len(Default.appenders)
	Default.Enable(ta)
	assert.True(t, oldcount == len(Default.appenders)-1)

	Default.Disable(ta.Id())
	assert.True(t, oldcount == len(Default.appenders))

	Default.Info("some msg")
	assert.Exactly(t, 0, ta.count)

	Default.Info("some msg")
	assert.Exactly(t, 0, ta.count)
}

func TestDisableInvalid(t *testing.T) {
	ta := &testAppender{}
	Default.Enable(ta)

	Default.Disable(123)
	assert.Exactly(t, 1, ta.errorCount)
}

func TestLogCalls(t *testing.T) {
	defer cleanupTest()

	defer func() {
		if r := recover(); r != nil {
		}
	}()

	ta := &testAppender{}
	Default.Enable(ta)

	Default.Debug("some msg")
	Default.Debug("some msg")

	Default.Info("some msg")
	Default.Info("some msg")

	Default.Warn("some msg")
	Default.Warn("some msg")

	Default.Error("some msg")
	Default.Error("some msg")

	Default.Panic("some msg")
	Default.Panic("some msg")

	assert.Exactly(t, 8, ta.count)
}

func TestLogCallsWithLevel(t *testing.T) {
	defer cleanupTest()

	ta := &testAppender{}
	Default.Enable(ta)

	Default.Level = WARN

	Default.Debug("some msg")
	Default.Debug("some msg")

	Default.Info("some msg")
	Default.Info("some msg")

	Default.Warn("some msg")
	Default.Warn("some msg")

	Default.Error("some msg")
	Default.Error("some msg")

	assert.Exactly(t, 4, ta.count)

	ta.count = 0
	Default.Level = DEBUG

	Default.Debug("some msg")
	Default.Debug("some msg")

	Default.Info("some msg")
	Default.Info("some msg")

	Default.Warn("some msg")
	Default.Warn("some msg")

	Default.Error("some msg")
	Default.Error("some msg")

	assert.Exactly(t, 8, ta.count)
}
