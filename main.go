package main

import (
	"fmt"
	"log"
)

func main() {
	config := &Config{BaseUrl: "http://localhost:8000"}
	client := New(config)

	devices, err := client.ListDevices()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(devices)

	device, err := client.GetDevice("HT83G1C00069")
	if err != nil && err.Error() != "device not found" {
		log.Fatal(err)
	}

	fmt.Println(device)

	deviceApps, err := client.GetDeviceApps("HT83G1C00069")
	if err != nil && err.Error() != "device not found" {
		log.Fatal(err)
	}

	fmt.Println(deviceApps)

	session, err := client.StartSession("foo", "bar")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(session)

	err = client.StopSession("foo")
	if err != nil {
		log.Fatal(err)
	}
 }
