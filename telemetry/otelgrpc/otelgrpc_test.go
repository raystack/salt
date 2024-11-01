package otelgrpc_test

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	commonv1 "github.com/goto/salt/server/example/proto/gotocompany/common/v1"
	"github.com/goto/salt/telemetry/otelgrpc"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
)

func Test_otelGRPCMonitor_Record(t *testing.T) {
	mt := otelgrpc.NewMeter("localhost:1001", otelgrpc.WithMeterName("/meterName"))
	assert.NotNil(t, mt)
	initialAttr := mt.GetAttributes()
	uc := mt.UnaryClientInterceptor()
	assert.NotNil(t, uc)
	assert.Equal(t, initialAttr, mt.GetAttributes())
	sc := mt.StreamClientInterceptor()
	assert.NotNil(t, sc)
	assert.Equal(t, initialAttr, mt.GetAttributes())
	mt.RecordUnary(context.Background(), otelgrpc.UnaryParams{
		Start:  time.Now(),
		Method: "/service.gojek.com/MethodName",
		Req:    nil,
		Res:    nil,
		Err:    nil,
	})
	assert.Equal(t, initialAttr, mt.GetAttributes())
	version := &commonv1.Version{
		Version: "1.0.0",
	}
	mt.RecordUnary(context.Background(), otelgrpc.UnaryParams{
		Start:  time.Now(),
		Method: "",
		Req:    &commonv1.GetVersionRequest{Client: version},
		Res:    nil,
		Err:    nil,
	})
	assert.Equal(t, initialAttr, mt.GetAttributes())
	mt.RecordStream(context.Background(), time.Now(), "", nil)
	assert.Equal(t, initialAttr, mt.GetAttributes())
	mt.RecordStream(context.Background(), time.Now(), "/service.gojek.com/MethodName", errors.New("dummy error"))
	assert.Equal(t, initialAttr, mt.GetAttributes())
}
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
func Test_getProtoSize(t *testing.T) {
	version := &commonv1.Version{
		Version: "1.0.0",
	}
	req := &commonv1.GetVersionResponse{Server: version}
	if got := otelgrpc.GetProtoSize(req); got != 9 {
		t.Errorf("getProtoSize() = %v, want %v", got, 9)
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
