package smarthome_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/betom84/go-alexa/smarthome"
)

func TestDefaultLoggerForTraceLevel(t *testing.T) {

	builder := &strings.Builder{}
	logger := smarthome.NewDefaultLogger(smarthome.Trace, builder)

	logger.Trace("This is a trace message")
	assert.True(t, strings.HasSuffix(builder.String(), "[go-alexa] [TRACE] This is a trace message\n"))
	builder.Reset()

	logger.Debug("This is a debug message")
	assert.True(t, strings.HasSuffix(builder.String(), "[go-alexa] [DEBUG] This is a debug message\n"))
	builder.Reset()

	logger.Info("This is a info message")
	assert.True(t, strings.HasSuffix(builder.String(), "[go-alexa] [INFO] This is a info message\n"))
	builder.Reset()

	logger.Warning("This is a warning message")
	assert.True(t, strings.HasSuffix(builder.String(), "[go-alexa] [WARN] This is a warning message\n"))
	builder.Reset()

	logger.Error("This is a %s message\n", "error")
	assert.True(t, strings.HasSuffix(builder.String(), "[go-alexa] [ERROR] This is a error message\n"))
	builder.Reset()

	logger.Fatal("This is a %s message\n", "fatal")
	assert.True(t, strings.HasSuffix(builder.String(), "[go-alexa] [FATAL] This is a fatal message\n"))
	builder.Reset()
}

func TestDefaultLoggerForFatalLevel(t *testing.T) {

	builder := &strings.Builder{}
	logger := smarthome.NewDefaultLogger(smarthome.Fatal, builder)

	logger.Trace("This is a trace message")
	logger.Debug("This is a debug message")
	logger.Info("This is a info message")
	logger.Warning("This is a warning message\n")
	logger.Error("This is a %s message\n", "error")
	assert.Equal(t, 0, builder.Len())

	logger.Fatal("This is a %s message\n", "fatal")
	assert.True(t, strings.HasSuffix(builder.String(), "[go-alexa] [FATAL] This is a fatal message\n"))
	builder.Reset()
}
