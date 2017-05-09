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

	c := pb.NewMonitorClient(conn)

	r, err := c.RegisterDevice(context.Background(),
		&pb.RegisterDeviceRequest{
			Devicetype:   pb.DeviceType_DRONE,
			Name:         "Drone Alpha 1",
			Owner:        "Kevin",
			Serialnumber: "DR-A1",
		})

	fmt.Printf("Registered Device, Reply %+v\n", r)
}
