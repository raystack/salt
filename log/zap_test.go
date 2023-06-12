package log_test

import (
	"bufio"
	"bytes"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"go.uber.org/zap"

	"github.com/raystack/salt/log"
)

type zapBufWriter struct {
	io.Writer
}

func (cw zapBufWriter) Close() error {
	return nil
}
func (cw zapBufWriter) Sync() error {
	return nil
}

type zapClock struct {
	t time.Time
}

func (m zapClock) Now() time.Time {
	return m.t
}

func (m zapClock) NewTicker(duration time.Duration) *time.Ticker {
	return time.NewTicker(duration)
}

func buildBufferedZapOption(writer io.Writer, t time.Time, bufWriterKey string) log.Option {
	config := zap.NewDevelopmentConfig()
	config.DisableCaller = true
	// register mock writer
	_ = zap.RegisterSink(bufWriterKey, func(u *url.URL) (zap.Sink, error) {
		return zapBufWriter{writer}, nil
	})
	// build a valid custom path
	customPath := fmt.Sprintf("%s:", bufWriterKey)
	config.OutputPaths = []string{customPath}

	return log.ZapWithConfig(config, zap.WithClock(&zapClock{
		t: t,
	}))
}

func TestZap(t *testing.T) {
	mockedTime := time.Date(2021, 6, 10, 11, 55, 0, 0, time.UTC)

	t.Run("should successfully print at info level", func(t *testing.T) {
		var b bytes.Buffer
		bWriter := bufio.NewWriter(&b)

		zapper := log.NewZap(buildBufferedZapOption(bWriter, mockedTime, randomString(10)))
		zapper.Info("hello", "wor", "ld")
		bWriter.Flush()

		assert.Equal(t, mockedTime.Format("2006-01-02T15:04:05.000Z0700")+"\tINFO\thello\t{\"wor\": \"ld\"}\n", b.String())
	})

	t.Run("should successfully print log from context", func(t *testing.T) {
		var b bytes.Buffer
		bWriter := bufio.NewWriter(&b)

		zapper := log.NewZap(buildBufferedZapOption(bWriter, mockedTime, randomString(10)))
		ctx := zapper.NewContext(context.Background())
		contextualLog := log.ZapFromContext(ctx)
		contextualLog.Info("hello", "wor", "ld")
		bWriter.Flush()

		assert.Equal(t, mockedTime.Format("2006-01-02T15:04:05.000Z0700")+"\tINFO\thello\t{\"wor\": \"ld\"}\n", b.String())
	})

	t.Run("should successfully print log from context with fields", func(t *testing.T) {
		var b bytes.Buffer
		bWriter := bufio.NewWriter(&b)

		zapper := log.NewZap(buildBufferedZapOption(bWriter, mockedTime, randomString(10)))
		ctx := zapper.NewContext(context.Background())
		ctx = log.ZapContextWithFields(ctx, zap.Int("one", 1))
		ctx = log.ZapContextWithFields(ctx, zap.String("two", "two"))
		log.ZapFromContext(ctx).Info("hello", "wor", "ld")
		bWriter.Flush()

		assert.Equal(t, mockedTime.Format("2006-01-02T15:04:05.000Z0700")+"\tINFO\thello\t{\"one\": 1, \"two\": \"two\", \"wor\": \"ld\"}\n", b.String())
	})
}

func randomString(n int) string {
	const alphabets = "abcdefghijklmnopqrstuvwxyz"
	var alphaBytes = make([]byte, n)
	rand.Read(alphaBytes)
	for i, b := range alphaBytes {
		alphaBytes[i] = alphabets[b%byte(len(alphabets))]
	}
	return string(alphaBytes)
}
