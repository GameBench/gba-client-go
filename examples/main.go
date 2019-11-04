package main

import (
	"fmt"
	"time"

	gba "github.com/GameBench/gba-client-go"
)

func main() {
	config := &Config{BaseUrl: "http://localhost:8000", Username: "ade@gamebench.net", Password: ""}
	client := New(config)

	devices, err := client.ListDevices()
	if err != nil {
		panic(err)
	}

	fmt.Println(devices)

	device, err := client.GetDevice("HT83G1C00069")
	if err != nil && err.Error() != "device not found" {
		panic(err)
	}

	fmt.Println(device)

	deviceApps, err := client.GetDeviceApps("HT83G1C00069")
	if err != nil && err.Error() != "device not found" {
		panic(err)
	}

	fmt.Println(deviceApps)

	session, err := client.StartSession("HT83G1C00069", "com.codigames.market.idle.tycoon", &StartSessionOptions{AutoSync: true, Screenshots: true})
	if err != nil {
		panic(err)
	}

	fmt.Println(session)

	time.Sleep(2 * time.Minute)

	err = client.StopSession(session.Id)
	if err != nil {
		panic(err)
	}
}
