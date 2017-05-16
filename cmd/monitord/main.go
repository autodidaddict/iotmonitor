package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"

	"google.golang.org/grpc"

	"golang.org/x/net/context"

	"github.com/autodidaddict/iotmonitor"

	"syscall"

	"github.com/autodidaddict/iotmonitor/pb"
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	port = ":50051"
)

var (
	gRPCAddr  = ":8081"
	httpAddr  = ":8080"
	debugAddr = ":8082"
)

func main() {
	ctx := context.Background()
	errChan := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	var telemetryUpdates, statusUpdates, devicesRegistered metrics.Counter
	{
		telemetryUpdates = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "iotmonitor",
			Name:      "telemetry_updates",
			Help:      "Total number of telemetry updates received.",
		}, []string{})
		statusUpdates = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "iotmonitor",
			Name:      "status_updates",
			Help:      "Total number of status updated received.",
		}, []string{})
		devicesRegistered = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "iotmonitor",
			Name:      "devices_registered",
			Help:      "Total number of devices registered.",
		}, []string{})
	}

	var srv iotmonitor.Service
	{
		srv = iotmonitor.NewService()
		srv = iotmonitor.ServiceInstrumentingMiddleware(telemetryUpdates, devicesRegistered, statusUpdates)(srv)
	}

	var duration metrics.Histogram
	{
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "iotmonitor",
			Name:      "request_duration_ns",
			Help:      "Request dureaction in nanoseconds.",
		}, []string{"method", "success"})
	}

	var registerEndpoint endpoint.Endpoint
	{
		registerDuration := duration.With("method", "register")
		registerEndpoint = iotmonitor.MakeRegisterEndpoint(srv)
		registerEndpoint = iotmonitor.EndpointInstrumentingMiddleware(registerDuration)(registerEndpoint)
	}

	var updateEndpoint endpoint.Endpoint
	{
		updateDuration := duration.With("method", "update")
		updateEndpoint = iotmonitor.MakeUpdateEndpoint(srv)
		updateEndpoint = iotmonitor.EndpointInstrumentingMiddleware(updateDuration)(updateEndpoint)
	}

	var telemetryEndpoint endpoint.Endpoint
	{
		telemetryDuration := duration.With("method", "telemetry")
		telemetryEndpoint = iotmonitor.MakeTelemetryEndpoint(srv)
		telemetryEndpoint = iotmonitor.EndpointInstrumentingMiddleware(telemetryDuration)(telemetryEndpoint)
	}

	endpoints := iotmonitor.Endpoints{
		UpdateEndpoint:    updateEndpoint,
		TelemetryEndpoint: telemetryEndpoint,
		RegisterEndpoint:  registerEndpoint,
	}

	// Debug/Diagnostics Transport
	go func() {
		log.Println("Debug http:", debugAddr)
		m := http.NewServeMux()
		m.Handle("/metrics", promhttp.Handler())
		errChan <- http.ListenAndServe(debugAddr, m)
	}()

	// HTTP Transport
	go func() {
		log.Println("http:", httpAddr)
		handler := iotmonitor.NewHTTPServer(ctx, endpoints)
		errChan <- http.ListenAndServe(httpAddr, handler)
	}()

	// gRPC Transport
	go func() {
		listener, err := net.Listen("tcp", gRPCAddr)
		if err != nil {
			errChan <- err
			return
		}
		log.Println("grpc:", gRPCAddr)
		handler := iotmonitor.NewGRPCServer(ctx, endpoints)
		gRPCServer := grpc.NewServer()
		pb.RegisterMonitorServer(gRPCServer, handler)
		errChan <- gRPCServer.Serve(listener)
	}()

	log.Fatalln(<-errChan)
}
