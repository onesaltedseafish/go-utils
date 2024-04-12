package log

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoggerOptGetFilePath(t *testing.T) {
	currentAbs, _ := filepath.Abs(".")

	testcases := []struct {
		Opt      LoggerOpt
		WantPath string
	}{
		{LoggerOpt{Directory: "/tmp", Name: "log.log"}, "/tmp/log.log"},
		{LoggerOpt{Directory: ".", Name: "log.log"}, filepath.Join(currentAbs, "log.log")},
	}

	for _, testcase := range testcases {
		assert.Equal(t, testcase.WantPath, testcase.Opt.GetLogFilePath())
	}
}

func TestGetDefaultLogger(t *testing.T) {
	assert.Equal(t, true, GetDefaultLogger() == nil)
	logger := GetLogger("abc", &defaultLogOpt)
	assert.Equal(t, logger, GetDefaultLogger())
}
