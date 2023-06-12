package log_test

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"

	"github.com/sirupsen/logrus"

	"github.com/raystack/salt/log"

	"github.com/stretchr/testify/assert"
)

func TestLogrus(t *testing.T) {
	t.Run("should parse info messages at debug level correctly", func(t *testing.T) {
		var b bytes.Buffer
		foo := bufio.NewWriter(&b)

		logger := log.NewLogrus(log.LogrusWithLevel("debug"), log.LogrusWithWriter(foo), log.LogrusWithFormatter(&logrus.TextFormatter{
			DisableTimestamp: true,
		}))
		logger.Info("hello world")
		foo.Flush()

		assert.Equal(t, "level=info msg=\"hello world\"\n", b.String())
	})
	t.Run("should not parse debug messages at info level correctly", func(t *testing.T) {
		var b bytes.Buffer
		foo := bufio.NewWriter(&b)

		logger := log.NewLogrus(log.LogrusWithLevel("info"), log.LogrusWithWriter(foo), log.LogrusWithFormatter(&logrus.TextFormatter{
			DisableTimestamp: true,
		}))
		logger.Debug("hello world")
		foo.Flush()

		assert.Equal(t, "", b.String())
	})
	t.Run("should parse field maps correctly", func(t *testing.T) {
		var b bytes.Buffer
		foo := bufio.NewWriter(&b)

		logger := log.NewLogrus(log.LogrusWithLevel("debug"), log.LogrusWithWriter(foo), log.LogrusWithFormatter(&logrus.TextFormatter{
			DisableTimestamp: true,
		}))
		logger.Debug("current values", "day", 11, "month", "aug")
		foo.Flush()

		assert.Equal(t, "level=debug msg=\"current values\" day=11 month=aug\n", b.String())
	})
	t.Run("should handle errors correctly", func(t *testing.T) {
		var b bytes.Buffer
		foo := bufio.NewWriter(&b)

		logger := log.NewLogrus(log.LogrusWithLevel("info"), log.LogrusWithWriter(foo), log.LogrusWithFormatter(&logrus.TextFormatter{
			DisableTimestamp: true,
		}))
		var err = fmt.Errorf("request failed")
		logger.Error(err.Error(), "hello", "world")
		foo.Flush()
		assert.Equal(t, "level=error msg=\"request failed\" hello=world\n", b.String())
	})
	t.Run("should ignore params if malformed", func(t *testing.T) {
		var b bytes.Buffer
		foo := bufio.NewWriter(&b)

		logger := log.NewLogrus(log.LogrusWithLevel("info"), log.LogrusWithWriter(foo), log.LogrusWithFormatter(&logrus.TextFormatter{
			DisableTimestamp: true,
		}))
		var err = fmt.Errorf("request failed")
		logger.Error(err.Error(), "hello", "world", "!")
		foo.Flush()
		assert.Equal(t, "level=error msg=\"request failed\"\n", b.String())
	})
}
