package instrumentation

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerMethods(t *testing.T) {
	msg := "test message"

	testCases := []struct {
		desc           string
		logFunc        func(Logger, string)
		expectedStdout string
		expectedStderr string
	}{
		{
			desc:           "info log",
			logFunc:        Logger.Info,
			expectedStdout: "^\\[INFO\\] \\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} test message\n$",
			expectedStderr: "",
		},
		{
			desc:           "warning log",
			logFunc:        Logger.Warning,
			expectedStdout: "^\\[WARNING\\] \\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} test message\n$",
			expectedStderr: "",
		},
		{
			desc:           "error log",
			logFunc:        Logger.Error,
			expectedStdout: "",
			expectedStderr: "^\\[ERROR\\] \\d{4}/\\d{2}/\\d{2} \\d{2}:\\d{2}:\\d{2} test message\n$",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			originalStdout := os.Stdout
			originalStderr := os.Stderr
			t.Cleanup(func() {
				os.Stdout = originalStdout
				os.Stderr = originalStderr
			})

			readStdout, writeStdout, err := os.Pipe()
			assert.NoError(t, err)
			os.Stdout = writeStdout

			readStderr, writeStderr, err := os.Pipe()
			assert.NoError(t, err)
			os.Stderr = writeStderr

			logger := NewLogger()
			tC.logFunc(logger, msg)

			writeStdout.Close()
			var bufStdout bytes.Buffer
			io.Copy(&bufStdout, readStdout)
			if tC.expectedStdout == "" {
				assert.Empty(t, bufStdout.String())
			}
			assert.Regexp(t, tC.expectedStdout, bufStdout.String())

			writeStderr.Close()
			var bufStderr bytes.Buffer
			io.Copy(&bufStderr, readStderr)
			if tC.expectedStderr == "" {
				assert.Empty(t, bufStderr.String())
			}
			assert.Regexp(t, tC.expectedStderr, bufStderr.String())
		})
	}
}
