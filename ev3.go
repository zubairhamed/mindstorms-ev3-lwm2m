package main

import (
    . "github.com/zubairhamed/go-lwm2m"
    . "github.com/zubairhamed/go-lwm2m/api"
    "github.com/zubairhamed/go-lwm2m/objects/oma"
    "github.com/zubairhamed/go-lwm2m/registry"
    "github.com/zubairhamed/mindstorms-ev3-lwm2m/obj"
)

func main() {
    // client := NewLWM2MClient(":0", "localhost:5683")
    client := NewLWM2MClient(":0", "192.168.1.212:5683")

    registry := registry.NewDefaultObjectRegistry()
    client.UseRegistry(registry)

    setupResources(client, registry)

    client.OnStartup(func() {
        client.Register("EV3")
    })

    client.Start()
}

func setupResources(client LWM2MClient, reg Registry) {
    device := ev3.NewDeviceObject(reg)

    client.EnableObject(oma.OBJECT_LWM2M_DEVICE, device)

    instanceDevice := reg.CreateObjectInstance(oma.OBJECT_LWM2M_DEVICE, 0)
    client.AddObjectInstances( instanceDevice, )
}
