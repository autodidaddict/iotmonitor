package iotmonitor

import (
	"fmt"

	"github.com/autodidaddict/iotmonitor/pb"
	"golang.org/x/net/context"
)

type grpcServer struct {
}

func NewGRPCServer(ctx context.Context) pb.MonitorServer {
	return &grpcServer{}
}

func (s *grpcServer) RegisterDevice(ctx context.Context, in *pb.RegisterDeviceRequest) (*pb.RegisterDeviceReply, error) {
	fmt.Printf("Received request to register device name %s, type %v\n", in.Name, in.Devicetype)
	return &pb.RegisterDeviceReply{Registered: true, Deviceid: 99}, nil
}

func (s *grpcServer) UpdateDeviceStatus(ctx context.Context, in *pb.StatusUpdateRequest) (*pb.StatusUpdateReply, error) {
	fmt.Printf("Received request to update device status %+v\n", in)
	return &pb.StatusUpdateReply{Acknowledged: true}, nil
}

func (s *grpcServer) SubmitTelemetry(ctx context.Context, in *pb.TelemetrySubmitRequest) (*pb.TelemetrySubmitReply, error) {
	fmt.Printf("Received telemetry %+v\n", in)
	return &pb.TelemetrySubmitReply{Acknowledged: true}, nil
}
