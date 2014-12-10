package golog

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"strings"
	"testing"
)

func TestStdoutAppender(t *testing.T) {
	old := os.Stdout
	defer func() {
		out = old
	}()

	r, w, _ := os.Pipe()
	out = w

	Default.Debug("debug")
	Default.Info("info")
	Default.Warn("warn")
	Default.Error("error")

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		io.Copy(&buf, r)
		outC <- buf.String()
	}()

	w.Close()
	out := <-outC
	fmt.Println(out)

	lines := strings.Split(out, "\n")
	for i, line := range lines {
		switch i {
		case 0:
			assert.Equal(t, "[34mdefault - [DEBUG] - debug[0m", line)
		case 1:
			assert.Equal(t, "[32mdefault - [INFO] - info[0m", line)
		case 2:
			assert.Equal(t, "[33mdefault - [WARN] - warn[0m", line)
		case 3:
			assert.Equal(t, "[31mdefault - [ERROR] - error[0m", line)
		default:
		}
	}
	assert.True(t, true)
}
