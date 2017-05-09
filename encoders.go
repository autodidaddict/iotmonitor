package iotmonitor

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"strconv"

	"github.com/autodidaddict/iotmonitor/pb"
	"github.com/gorilla/mux"
)

type registerRequest struct {
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
	Owner        string `json:"owner"`
	DeviceType   string `json:"device_type"`
}

type registerReply struct {
	Registered bool   `json:"registered"`
	DeviceID   uint64 `json:"device_id"`
	Err        string `json:"err,omitempty"`
}

type updateRequest struct {
	DeviceID         uint64   `json:"device_id"`
	Location         location `json:"location"`
	BatteryRemaining uint32   `json:"battery_remaining"`
}

type location struct {
	Longitude float32 `json:"longitude"`
	Latitude  float32 `json:"latitude"`
	Altitude  float32 `json:"altitude"`
}

type updateReply struct {
	Acknowledged bool   `json:"acknowledged"`
	Err          string `json:"err,omitempty"`
}

type telemetryRequest struct {
	DeviceID uint64             `json:"device_id"`
	Readings map[string]float32 `json:"readings"`
}

type telemetryReply struct {
	Acknowledged bool   `json:"acknowledged"`
	Err          string `json:"err,omitempty"`
}

var errBadRoute = errors.New("bad route")

func decodeRegisterRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	var req registerRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}

	return req, nil
}

func decodeUpdateRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}

	var req updateRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	realID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	req.DeviceID = realID
	return req, nil
}

func decodeTelemetryRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}

	var req telemetryRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		return nil, err
	}
	realID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		return nil, err
	}
	req.DeviceID = realID

	return req, nil
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	return json.NewEncoder(w).Encode(response)
}

// GRPC Encode -> to protobuf
// GRPC Decode -> from protobuf

func EncodeGRPCRegisterRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(registerRequest)
	dt := pb.DeviceType_DRONE
	if req.DeviceType == "Sensor" {
		dt = pb.DeviceType_DRONE
	}
	return &pb.RegisterDeviceRequest{Devicetype: dt, Name: req.Name, Owner: req.Owner, Serialnumber: req.SerialNumber}, nil
}

func DecodeGRPCRegisterRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.RegisterDeviceRequest)
	dt := "Sensor"
	if req.Devicetype == pb.DeviceType_DRONE {
		dt = "Drone"
	}
	return registerRequest{DeviceType: dt, Name: req.Name, Owner: req.Owner, SerialNumber: req.Serialnumber}, nil
}

func EncodeGRPCRegisterResponse(ctx context.Context, r interface{}) (interface{}, error) {
	res := r.(registerReply)
	return &pb.RegisterDeviceReply{Deviceid: res.DeviceID, Registered: res.Registered, Err: res.Err}, nil
}

func DecodeGRPCRegisterResponse(ctx context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.RegisterDeviceReply)
	return registerReply{DeviceID: res.Deviceid, Registered: res.Registered, Err: res.Err}, nil
}

func EncodeGRPCUpdateRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(updateRequest)
	return &pb.StatusUpdateRequest{Batteryremaining: req.BatteryRemaining, Deviceid: req.DeviceID, Location: &pb.Location{
		Altitude:  req.Location.Altitude,
		Longitude: req.Location.Longitude,
		Latitude:  req.Location.Latitude,
	}}, nil
}

func DecodeGRPCUpdateRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.StatusUpdateRequest)
	return updateRequest{BatteryRemaining: req.Batteryremaining, DeviceID: req.Deviceid, Location: location{
		Altitude:  req.Location.Altitude,
		Longitude: req.Location.Longitude,
		Latitude:  req.Location.Latitude,
	}}, nil
}

func EncodeGRPCUpdateResponse(ctx context.Context, r interface{}) (interface{}, error) {
	res := r.(updateReply)
	return &pb.StatusUpdateReply{Acknowledged: res.Acknowledged, Err: res.Err}, nil
}

func DecodeGRPCUpdateResponse(ctx context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.StatusUpdateReply)
	return updateReply{Acknowledged: res.Acknowledged, Err: res.Err}, nil
}

func EncodeGRPCTelemetryRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(telemetryRequest)
	return &pb.TelemetrySubmitRequest{Deviceid: req.DeviceID, Readings: req.Readings}, nil
}

func DecodeGRPCTelemetryRequest(ctx context.Context, r interface{}) (interface{}, error) {
	req := r.(*pb.TelemetrySubmitRequest)
	return telemetryRequest{DeviceID: req.Deviceid, Readings: req.Readings}, nil
}

func EncodeGRPCTelemetryResponse(ctx context.Context, r interface{}) (interface{}, error) {
	res := r.(telemetryReply)
	return &pb.TelemetrySubmitReply{Acknowledged: res.Acknowledged, Err: res.Err}, nil
}

func DecodeGRPCTelemetryResponse(ctx context.Context, r interface{}) (interface{}, error) {
	res := r.(*pb.TelemetrySubmitReply)
	return telemetryReply{Acknowledged: res.Acknowledged, Err: res.Err}, nil
}
