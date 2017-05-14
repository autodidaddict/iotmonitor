package main

import (
	"context"
	"fmt"

	"github.com/autodidaddict/iotmonitor/pb"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8081"
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	ctx := context.Background()

	c := pb.NewMonitorClient(conn)

	r, err := c.RegisterDevice(ctx,
		&pb.RegisterDeviceRequest{
			Devicetype:   pb.DeviceType_DRONE,
			Name:         "Drone Alpha 1",
			Owner:        "Kevin",
			Serialnumber: "DR-A1",
		})

	fmt.Printf("Registered Device, Reply %+v\n", r)

	for i := 0; i < 10; i++ {
		_, err := c.UpdateDeviceStatus(ctx, &pb.StatusUpdateRequest{
			Batteryremaining: 80,
			Deviceid:         r.Deviceid,
			Location: &pb.Location{
				Altitude:  100.0,
				Longitude: 35.4,
				Latitude:  20.1 * float32(i),
			},
		})
		if err != nil {
			fmt.Println(err)
		}
	}

	readingsMap := make(map[string]float32)
	readingsMap["temp"] = 64.0
	for j := 0; j < 10; j++ {
		readingsMap["readingCount"] = float32(j)
		_, err := c.SubmitTelemetry(ctx, &pb.TelemetrySubmitRequest{
			Deviceid: r.Deviceid,
			Readings: readingsMap,
		})
		if err != nil {
			fmt.Println(err)
		}
	}
}
