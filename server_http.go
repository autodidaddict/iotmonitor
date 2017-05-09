package iotmonitor

import (
	"net/http"

	"github.com/gorilla/mux"

	httptransport "github.com/go-kit/kit/transport/http"
	"golang.org/x/net/context"
)

func NewHTTPServer(ctx context.Context, endpoints Endpoints) http.Handler {
	m := mux.NewRouter()

	registerHandler := httptransport.NewServer(
		endpoints.RegisterEndpoint,
		decodeRegisterRequest,
		encodeResponse,
	)

	statusUpdateHandler := httptransport.NewServer(
		endpoints.UpdateEndpoint,
		decodeUpdateRequest,
		encodeResponse,
	)

	telemetryUpdateHandler := httptransport.NewServer(
		endpoints.TelemetryEndpoint,
		decodeTelemetryRequest,
		encodeResponse,
	)

	m.Handle("/v1/devices", registerHandler).Methods("POST")
	m.Handle("/v1/devices/{id}/status", statusUpdateHandler).Methods("PUT")
	m.Handle("/v1/devices/{id}/telemetry", telemetryUpdateHandler).Methods("PUT")
	return m
}
