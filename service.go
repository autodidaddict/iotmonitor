package iotmonitor

import (
	"fmt"
	"time"

	"github.com/garyburd/redigo/redis"

	"github.com/go-kit/kit/metrics"
	"golang.org/x/net/context"
)

type Service interface {
	RegisterDevice(ctx context.Context, name, owner, deviceType string) (id uint64, err error)
	UpdateStatus(ctx context.Context, id uint64, lat, long, alt float32, battery uint32) (bool, error)
	SubmitTelemetry(ctx context.Context, id uint64, readings map[string]float32) (bool, error)
}

type Middleware func(Service) Service

func NewService() Service {
	return &monitorService{}
}

type monitorService struct{}

func (monitorService) RegisterDevice(ctx context.Context, name, owner, deviceType string) (id uint64, err error) {
	fmt.Printf("Registering device name %s, type %s\n", name, deviceType)

	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		// handle error
		return
	}
	defer c.Close()
	id, err = redis.Uint64(c.Do("INCR", "id:devices"))
	if err != nil {
		return
	}
	var newDevice struct {
		Name       string `redis:"name"`
		Owner      string `redis:"owner"`
		DeviceType string `redis:"device_type"`
		ID         uint64 `redis:"id"`
	}
	newDevice.DeviceType = deviceType
	newDevice.ID = id
	newDevice.Name = name
	newDevice.Owner = owner

	deviceKey := fmt.Sprintf("device:%d", id)
	if _, err := c.Do("HMSET", redis.Args{}.Add(deviceKey).AddFlat(&newDevice)...); err != nil {
		fmt.Printf("Failed to HMSET device %s\n", deviceKey)
		fmt.Println(err)
		return 0, err
	}

	_, err = c.Do("SADD", "devices", id)

	return
}

func (monitorService) UpdateStatus(ctx context.Context, id uint64, lat, long, alt float32, battery uint32) (bool, error) {
	fmt.Printf("Updating status for device %d, battery left %d .\n", id, battery)

	var lastStatus struct {
		Latitude  float32 `redis:"lat"`
		Longitude float32 `redis:"long"`
		Altitude  float32 `redis:"alt"`
		Battery   uint32  `redis:"battery"`
		Timestamp int64   `redis:"timestamp"`
	}
	lastStatus.Altitude = alt
	lastStatus.Longitude = long
	lastStatus.Latitude = lat
	lastStatus.Battery = battery
	lastStatus.Timestamp = makeTimestamp()

	statusKey := fmt.Sprintf("status:%d", id)
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		// handle error
		return false, err
	}
	defer c.Close()

	if _, err := c.Do("HMSET", redis.Args{}.Add(statusKey).AddFlat(&lastStatus)...); err != nil {
		fmt.Printf("Failed to HMSET status update %s\n", statusKey)
		fmt.Println(err)
		return false, err
	}

	return true, nil
}

func (monitorService) SubmitTelemetry(ctx context.Context, id uint64, readings map[string]float32) (bool, error) {
	fmt.Printf("Submitting telemetry for device %d,  %+v\n", id, readings)

	telemetryKey := fmt.Sprintf("telemetry:%d", id)
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		// handle error
		return false, err
	}
	defer c.Close()

	if _, err := c.Do("HMSET", redis.Args{}.Add(telemetryKey).AddFlat(readings)...); err != nil {
		fmt.Println(err)
		return false, err
	}

	if _, err := c.Do("HMSET", telemetryKey, "timestamp", makeTimestamp()); err != nil {
		fmt.Println(err)
		return false, err
	}

	return true, nil
}

func makeTimestamp() int64 {
	return time.Now().UTC().UnixNano() / (int64(time.Millisecond) / int64(time.Nanosecond))
}

type serviceInstrumentingMiddleware struct {
	telemetryUpdates  metrics.Counter
	devicesRegistered metrics.Counter
	statusUpdates     metrics.Counter
	next              Service
}

func ServiceInstrumentingMiddleware(telemetryUpdates, devicesRegistered, statusUpdates metrics.Counter) Middleware {
	return func(next Service) Service {
		return serviceInstrumentingMiddleware{
			telemetryUpdates:  telemetryUpdates,
			statusUpdates:     statusUpdates,
			devicesRegistered: devicesRegistered,
			next:              next,
		}
	}
}

func (mw serviceInstrumentingMiddleware) RegisterDevice(ctx context.Context, name, owner, deviceType string) (id uint64, err error) {
	v, err := mw.next.RegisterDevice(ctx, name, owner, deviceType)
	mw.devicesRegistered.Add(float64(1))
	return v, err
}
func (mw serviceInstrumentingMiddleware) UpdateStatus(ctx context.Context, id uint64, lat, long, alt float32, battery uint32) (bool, error) {
	v, err := mw.next.UpdateStatus(ctx, id, lat, long, alt, battery)
	mw.statusUpdates.Add(float64(1))
	return v, err
}
func (mw serviceInstrumentingMiddleware) SubmitTelemetry(ctx context.Context, id uint64, readings map[string]float32) (bool, error) {
	v, err := mw.next.SubmitTelemetry(ctx, id, readings)
	mw.telemetryUpdates.Add(float64(1))
	return v, err
}
