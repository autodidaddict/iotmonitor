package iotmonitor

import (
	"fmt"

	"golang.org/x/net/context"
)

type Service interface {
	RegisterDevice(ctx context.Context, name, owner, deviceType string) (id uint64, err error)
	UpdateStatus(ctx context.Context, id uint64, lat, long, alt float32, battery uint32) (bool, error)
	SubmitTelemetry(ctx context.Context, id uint64, readings map[string]float32) (bool, error)
}

func NewService() Service {
	return &monitorService{}
}

type monitorService struct{}

func (monitorService) RegisterDevice(ctx context.Context, name, owner, deviceType string) (id uint64, err error) {
	fmt.Printf("Registering device name %s, type %s\n", name, deviceType)
	return 99, nil
}

func (monitorService) UpdateStatus(ctx context.Context, id uint64, lat, long, alt float32, battery uint32) (bool, error) {
	fmt.Printf("Updating status for device %d, battery left %d .\n", id, battery)
	return true, nil
}

func (monitorService) SubmitTelemetry(ctx context.Context, id uint64, readings map[string]float32) (bool, error) {
	fmt.Printf("Submitting telemetry for device %d,  %+v\n", id, readings)
	return true, nil
}
