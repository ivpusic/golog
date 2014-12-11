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
	assert.Exactly(t, len(lines), 5)
}
