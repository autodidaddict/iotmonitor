package iotmonitor

import (
	"context"
	"errors"

	"github.com/go-kit/kit/endpoint"
)

func MakeRegisterEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(registerRequest)
		v, err := srv.RegisterDevice(ctx, req.Name, req.Owner, req.DeviceType)
		if err != nil {
			return registerReply{DeviceID: 0, Registered: false, Err: err.Error()}, nil
		}
		return registerReply{DeviceID: v, Registered: true}, nil
	}
}

func MakeUpdateEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(updateRequest)
		v, err := srv.UpdateStatus(ctx, req.DeviceID, req.Location.Latitude, req.Location.Longitude, req.Location.Altitude,
			req.BatteryRemaining)
		if err != nil {
			return updateReply{Acknowledged: false, Err: err.Error()}, nil
		}
		return updateReply{Acknowledged: v}, nil
	}
}

func MakeTelemetryEndpoint(srv Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(telemetryRequest)
		v, err := srv.SubmitTelemetry(ctx, req.DeviceID, req.Readings)
		if err != nil {
			return telemetryReply{Acknowledged: false, Err: err.Error()}, nil
		}
		return telemetryReply{Acknowledged: v}, nil
	}
}

type Endpoints struct {
	RegisterEndpoint  endpoint.Endpoint
	UpdateEndpoint    endpoint.Endpoint
	TelemetryEndpoint endpoint.Endpoint
}

func (e Endpoints) RegisterDevice(ctx context.Context, name, owner, deviceType string) (id uint64, err error) {
	req := registerRequest{DeviceType: deviceType, Name: name, Owner: owner}
	resp, err := e.RegisterEndpoint(ctx, req)
	if err != nil {
		return 0, err
	}
	registerResp := resp.(registerReply)
	if registerResp.Err != "" {
		return 0, errors.New(registerResp.Err)
	}
	return registerResp.DeviceID, nil
}

func (e Endpoints) UpdateStatus(ctx context.Context, deviceId uint64, lat, long, alt float32, battery uint32) (bool, error) {
	req := updateRequest{BatteryRemaining: battery, DeviceID: deviceId, Location: location{
		Latitude: lat, Longitude: long, Altitude: alt},
	}
	resp, err := e.UpdateEndpoint(ctx, req)
	if err != nil {
		return false, err
	}
	updateResp := resp.(updateReply)
	if updateResp.Err != "" {
		return false, errors.New(updateResp.Err)
	}
	return updateResp.Acknowledged, nil
}

func (e Endpoints) SubmitTelemetry(ctx context.Context, id uint64, readings map[string]float32) (bool, error) {
	req := telemetryRequest{DeviceID: id, Readings: readings}
	resp, err := e.TelemetryEndpoint(ctx, req)
	if err != nil {
		return false, err
	}
	telemetryResp := resp.(telemetryReply)
	if telemetryResp.Err != "" {
		return false, errors.New(telemetryResp.Err)
	}
	return telemetryResp.Acknowledged, nil
}
