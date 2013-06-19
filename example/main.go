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
		fmt.Println("NOTE: Try toggling your Bluetooth or putting the Sphero asleep (and reawakening) to fix 'resource busy' errors")
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

	/* fmt.Println("Setting color...")
	s.SetRGBLEDOutput(0, 0, 255, true, ch)
	res = <-ch
	fmt.Printf("Set Color %#x\n", res) */

	fmt.Println("Getting color...")
	s.GetRGBLED(ch)
	res = <-ch
	c, _ := sphero.ParseColor(res.Data)
	fmt.Printf("Get Color %#x\n", res, c)

	fmt.Println("Power notifications...")
	s.SetPowerNotification(true, ch)
	res = <-ch
	fmt.Printf("Power notification %#x\n", res)

	// ASYNC

	mask := sphero.ApplyMasks32([]uint32{
		sphero.ACCEL_AXIS_X_RAW, sphero.ACCEL_AXIS_Y_RAW, sphero.ACCEL_AXIS_Z_RAW,
	})
	fmt.Printf("Set data streaming... %#x\n", mask)
	s.SetDataStreaming(40, 1, mask, 0, 0, ch)
	res = <-ch
	fmt.Printf("Data streaming %#x\n", res)

	<-time.Tick(10 * time.Second)

	// CLEANUP

	fmt.Println("Sleeping...")
	s.Sleep(time.Duration(0), 0, 0, ch)
	res = <-ch
	fmt.Printf("Slept %#x\n", res)

	fmt.Println("Closing...")
	s.Close()
	fmt.Println("Closed.")
}
