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
)

const (
	port = ":50051"
)

var (
	gRPCAddr = ":8081"
	httpAddr = ":8080"
)

func main() {
	ctx := context.Background()
	srv := iotmonitor.NewService()
	errChan := make(chan error)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errChan <- fmt.Errorf("%s", <-c)
	}()

	registerEndpoint := iotmonitor.MakeRegisterEndpoint(srv)
	updateEndpoint := iotmonitor.MakeUpdateEndpoint(srv)
	telemetryEndpoint := iotmonitor.MakeTelemetryEndpoint(srv)
	endpoints := iotmonitor.Endpoints{
		UpdateEndpoint:    updateEndpoint,
		TelemetryEndpoint: telemetryEndpoint,
		RegisterEndpoint:  registerEndpoint,
	}

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
