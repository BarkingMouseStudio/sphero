package main

import (
	"fmt"
	"github.com/Freeflow/sphero"
	"time"
)

func main() {
	fmt.Println("Connecting...")

	async := make(chan *sphero.AsyncResponse, 256)
	go func() {
		for r := range async {
			fmt.Println(r)
		}
	}()

	var err error
	var s *sphero.Sphero
	if s, err = sphero.NewSphero("/dev/cu.Sphero-YBR-RN-SPP", async); err != nil {
		fmt.Println(err)
		fmt.Println("NOTE: Try toggling Bluetooth and/or Sphero to fix 'resource busy' errors")
		return
	}

	fmt.Println("Connected.")

	var res *sphero.Response
	ch := make(chan *sphero.Response, 1)

	// PING

	fmt.Println("Ping...")
	s.Ping(ch)
	res = <-ch
	fmt.Printf("Pong %#x\n", res)

	// COLOR

	fmt.Println("Setting color...")
	s.SetRGBLEDOutput(0, 0, 255, false, ch)
	res = <-ch
	fmt.Printf("Set Color %#x\n", res)

	fmt.Println("Getting color...")
	s.GetRGBLED(ch)
	res = <-ch
	c, _ := res.Color()
	fmt.Printf("Get Color %#x\n", res, c)

	<-time.Tick(3 * time.Second)

	// CLEANUP

	fmt.Println("Sleeping...")
	s.Sleep(time.Duration(0), 0, 0, ch)
	res = <-ch
	fmt.Printf("Slept %#x\n", res)

	fmt.Println("Closing...")
	s.Close()
	fmt.Println("Closed.")
}
