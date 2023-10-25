package logger

import (
	"bytes"
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoggLevel(t *testing.T) {
	debugMsg := "i'm debug msg"
	infoMsg := "i'm info msg"
	errMsg := "i'm error msg"
	warnMsh := "i'm warn msg"

	tests := []struct {
		level          string
		expContains    []string
		expNotContains []string
	}{
		{level: "DEBUG", expContains: []string{debugMsg, infoMsg, warnMsh, errMsg}, expNotContains: []string{}},
		{level: "INFO", expContains: []string{infoMsg, errMsg, warnMsh}, expNotContains: []string{debugMsg}},
		{level: "WARN", expContains: []string{warnMsh, errMsg}, expNotContains: []string{debugMsg, infoMsg}},
		{level: "ERROR", expContains: []string{errMsg}, expNotContains: []string{debugMsg, infoMsg, warnMsh}},
	}

	defer log.SetOutput(os.Stderr)

	for _, tc := range tests {
		tc := tc

		var buf bytes.Buffer
		log.SetOutput(&buf)

		t.Run(tc.level, func(t *testing.T) {
			l := New(tc.level, "test")

			l.Debug(debugMsg)
			l.Info(infoMsg)
			l.Error(errMsg)
			l.Warn(warnMsh)

			out := buf.String()
			defer buf.Reset()

			for _, msg := range tc.expContains {
				require.Contains(t, out, msg)
			}
			for _, msg := range tc.expNotContains {
				require.NotContains(t, out, msg)
			}
		})
	}
}
