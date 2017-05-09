package iotmonitor

import (
	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	m := http.NewServeMux()
	m.Handle("/v1/devices", httptransport.NewServer(
		endpoints.RegisterEndpoint,
		decodeRegisterRequest,
		encodeResponse,
	))
	m.Handle("/v1/devices/updates", httptransport.NewServer(
		endpoints.UpdateEndpoint,
		decodeUpdateRequest,
		encodeResponse,
	))
	m.Handle("/v1/devices/telemetry", httptransport.NewServer(
		endpoints.TelemetryEndpoint,
		decodeTelemetryRequest,
		encodeResponse,
	))
	return m
}
