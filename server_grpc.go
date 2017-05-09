package iotmonitor

import (
	"github.com/autodidaddict/iotmonitor/pb"
	grpctransport "github.com/go-kit/kit/transport/grpc"
	"golang.org/x/net/context"
)

func NewGRPCServer(ctx context.Context, endpoints Endpoints) pb.MonitorServer {
	return &grpcServer{
		register: grpctransport.NewServer(
			endpoints.RegisterEndpoint,
			DecodeGRPCRegisterRequest,
			EncodeGRPCRegisterResponse,
		),
		update: grpctransport.NewServer(
			endpoints.UpdateEndpoint,
			DecodeGRPCUpdateRequest,
			EncodeGRPCUpdateResponse,
		),
		telemetry: grpctransport.NewServer(
			endpoints.TelemetryEndpoint,
			DecodeGRPCTelemetryRequest,
			EncodeGRPCTelemetryResponse,
		),
	}
}

type grpcServer struct {
	register  grpctransport.Handler
	update    grpctransport.Handler
	telemetry grpctransport.Handler
}

func (s *grpcServer) RegisterDevice(ctx context.Context, in *pb.RegisterDeviceRequest) (*pb.RegisterDeviceReply, error) {
	_, resp, err := s.register.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.RegisterDeviceReply), nil
}

func (s *grpcServer) UpdateDeviceStatus(ctx context.Context, in *pb.StatusUpdateRequest) (*pb.StatusUpdateReply, error) {
	_, resp, err := s.update.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.StatusUpdateReply), nil
}

func (s *grpcServer) SubmitTelemetry(ctx context.Context, in *pb.TelemetrySubmitRequest) (*pb.TelemetrySubmitReply, error) {
	_, resp, err := s.update.ServeGRPC(ctx, in)
	if err != nil {
		return nil, err
	}
	return resp.(*pb.TelemetrySubmitReply), nil
}
