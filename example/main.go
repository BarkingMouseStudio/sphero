package main

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/Freeflow/sphero"
	"time"
)

func main() {
	var conf sphero.Config
	var err error
	if _, err = toml.DecodeFile("config.toml", &conf); err != nil {
		panic(err)
	}

	fmt.Println("Connecting...")

	async := make(chan *sphero.AsyncResponse, 256)
	go func(async <-chan *sphero.AsyncResponse) {
		for r := range async {
			fmt.Println(r)
		}
	}(async)

	var s *sphero.Sphero
	if s, err = sphero.NewSphero(&conf, async); err != nil {
		panic(err)
	}

	var res *sphero.Response
	ch := make(chan *sphero.Response, 8)

	fmt.Println("Pinging...")
	s.Ping(ch)
	res = <-ch
	fmt.Printf("Pong %#x\n", res)

	fmt.Println("Setting color...")
	s.SetRGBLEDOutput(0, 0, 255, ch)
	res = <-ch
	fmt.Printf("Set Color %#x\n", res)

	fmt.Println("Getting color...")
	s.GetRGBLED(ch)
	res = <-ch
	fmt.Printf("Get Color %#x\n", res)

	<-time.Tick(10 * time.Second)

	fmt.Println("Sleeping...")
	s.Sleep(time.Duration(0), 0, 0, ch)
	res = <-ch
	fmt.Printf("Slept %#x\n", res)

	fmt.Println("Closing...")
	s.Close()
	fmt.Println("Done.")
}
