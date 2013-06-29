package main

import (
	"fmt"
	"github.com/FreeFlow/sphero"
	"os"
	"os/signal"
	"time"
)

type SensorData struct {
	AccelX, AccelY, AccelZ int16
}

func main() {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	// Set up a basic async listener
	async := make(chan *sphero.AsyncResponse, 256)
	go func() {
		var d SensorData
		for r := range async {
			r.Sensors(&d)
			fmt.Printf("Async: %#v\n", r, d)
		}
	}()

	fmt.Println("Connecting...")

	var err error
	var s *sphero.Sphero

	// TODO: Replace the first argument of `NewSphero` with your Sphero device name:
	if s, err = sphero.NewSphero("/dev/cu.Sphero-YBR-RN-SPP", async); err != nil {
		fmt.Println(err)
		fmt.Println("NOTE: Try toggling Bluetooth and/or Sphero to fix 'resource busy' errors")
		return
	}

	/*
		Remember to close. Sometimes OS X has trouble letting go...
		Toggling bluetooth and/or sleeping the Sphero fixes this.
	*/
	defer func() {
		fmt.Println("Closing...")
		s.Close()
		fmt.Println("Closed.")
	}()

	fmt.Println("Connected.")

	var res *sphero.Response
	ch := make(chan *sphero.Response, 1)

	// Send a ping command to verify that we're working
	fmt.Println("Ping...")
	s.Ping(ch)
	res = <-ch
	fmt.Printf("Pong %#x\n", res)

	// Enable data streaming - async messages are captured in the above goroutine
	fmt.Println("Enabling streaming...")
	// 400hz / (N = 400): 1hz or 1 async response per second
	s.SetDataStreaming(400, 1, 0, []uint32{sphero.ACCEL_RAW}, []uint32{}, ch)
	res = <-ch
	fmt.Printf("Streaming enabled %#x\n", res)

	fmt.Println("Press Ctrl+C to QUIT")
	<-sig

	// Sleeping the Sphero
	fmt.Println("Sleeping...")
	s.Sleep(time.Duration(0), 0, 0, ch)
	res = <-ch
	fmt.Printf("Slept %#x\n", res)
}
