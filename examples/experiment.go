package main

import (
	"fmt"
	"log"

	"github.com/mojbro/gosoundio"
)

func main() {
	fmt.Println("Testing libsoundio wrapper")
	soundio, err := gosoundio.CreateSoundIO()
	if err != nil {
		log.Fatal("Failed to create SoundIO context", err)
	}
	fmt.Printf("Created: %v\n", soundio)
	defer soundio.Destroy()

	if err := soundio.Connect(); err != nil {
		log.Fatal("Connect failed:", err)
	}
	defer soundio.Disconnect()

	idx, err := soundio.DefaultOutputDeviceIndex()
	if err != nil {
		log.Fatal("Couldn't get default output device", err)
	}
	fmt.Println("Got default device index:", idx)

	device, err := soundio.GetOutputDevice(idx)
	if err != nil {
		log.Fatal("GetOutputDevice: ", err)
	}
	defer device.Unref()

	fmt.Printf("Got device: %q, is raw: %v\n", device.Name(), device.IsRaw())

	outStream, err := gosoundio.NewOutStream(*device, func(outstream *gosoundio.OutStream, min int, max int) {
		fmt.Printf("callback!! min=%d, max=%d, outstream=%v\n", min, max, outstream)
		err := outstream.ClearBuffer()
		if err != nil {
			log.Fatal("error clearing buffer: ", err)
		}
	})
	if err != nil {
		log.Fatal("CreateOutstream: ", err)
	}
	defer outStream.Destroy()
	fmt.Printf("Created out stream: %v\n", outStream)

	err = outStream.Open()
	if err != nil {
		log.Fatal("open out stream: ", err)
	}
	fmt.Printf("Opened device: %v\n", outStream)

	outStream.Run()
}
