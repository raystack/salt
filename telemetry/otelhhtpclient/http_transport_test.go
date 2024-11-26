package otelhttpclient_test

import (
	"testing"

	otelhttpclient "github.com/raystack/salt/telemetry/otelhhtpclient"
	"github.com/stretchr/testify/assert"
)

func TestNewHTTPTransport(t *testing.T) {
	tr := otelhttpclient.NewHTTPTransport(nil)
	assert.NotNil(t, tr)
}
