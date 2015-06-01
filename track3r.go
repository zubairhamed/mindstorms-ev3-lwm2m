package main

import (
	. "github.com/zubairhamed/go-lwm2m/api"
	"github.com/zubairhamed/go-lwm2m/client"
	"github.com/zubairhamed/go-lwm2m/objects/ipso"
	"github.com/zubairhamed/go-lwm2m/objects/oma"
	"github.com/zubairhamed/go-lwm2m/registry"
	"github.com/zubairhamed/mindstorms-ev3-lwm2m/obj"
)

func main() {
	client := client.NewDefaultClient(":0", "localhost:5683")
	// client := NewLWM2MClient(":0", "192.168.1.212:5683")

	registry := registry.NewDefaultObjectRegistry()
	client.UseRegistry(registry)

	setupResources(client, registry)

	client.OnStartup(func() {
		client.Register("track3r")
	})

	client.Start()
}

func setupResources(client LWM2MClient, reg Registry) {
	device := ev3.NewDeviceObject(reg)

	client.EnableObject(oma.OBJECT_LWM2M_DEVICE, device)
	client.EnableObject(ipso.OBJECT_IPSO_ACTUATION, nil)

	act1 := reg.CreateObjectInstance(ipso.OBJECT_IPSO_ACTUATION, 0)
	client.AddObjectInstances(act1)

	instanceDevice := reg.CreateObjectInstance(oma.OBJECT_LWM2M_DEVICE, 0)
	client.AddObjectInstances(instanceDevice)
}

/*
Motors

Sound

LED

*/
