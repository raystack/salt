package otelgrpc_test

import (
	"reflect"
	"testing"

	"github.com/raystack/salt/telemetry/otelgrpc"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func Test_parseFullMethod(t *testing.T) {
	type args struct {
		fullMethod string
	}
	tests := []struct {
		name string
		args args
		want []attribute.KeyValue
	}{
		{name: "should parse correct method", args: args{
			fullMethod: "/test.service.name/MethodNameV1",
		}, want: []attribute.KeyValue{
			semconv.RPCService("test.service.name"),
			semconv.RPCMethod("MethodNameV1"),
		}},
		{name: "should return empty attributes on incorrect method", args: args{
			fullMethod: "incorrectMethod",
		}, want: nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := otelgrpc.ParseFullMethod(tt.args.fullMethod); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFullMethod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestExtractAddress(t *testing.T) {
	gotHost, gotPort := otelgrpc.ExtractAddress("localhost:1001")
	assert.Equal(t, "localhost", gotHost)
	assert.Equal(t, "1001", gotPort)
	gotHost, gotPort = otelgrpc.ExtractAddress("localhost")
	assert.Equal(t, "localhost", gotHost)
	assert.Equal(t, "80", gotPort)
	gotHost, gotPort = otelgrpc.ExtractAddress("some.address.golabs.io:15010")
	assert.Equal(t, "some.address.golabs.io", gotHost)
	assert.Equal(t, "15010", gotPort)
}
