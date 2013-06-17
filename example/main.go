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

	async := make(chan interface{})
	var s *sphero.Sphero
	if s, err = sphero.NewSphero(&conf, async); err != nil {
		panic(err)
	}

	fmt.Println("Pinging...")
	ping := make(chan *sphero.Response)
	s.Ping(ping)
	pong := <-ping
	fmt.Printf("Pong %#x\n", pong)

	fmt.Println("Setting color...")
	setColor := make(chan *sphero.Response)
	s.SetRGBLEDOutput(0, 0, 255, setColor)
	r := <-setColor
	fmt.Printf("Set Color %#x\n", r)

	fmt.Println("Getting color...")
	getColor := make(chan *sphero.Response)
	s.GetRGBLED(getColor)
	r = <-getColor
	fmt.Printf("Get Color %#x\n", r)

	<-time.Tick(10 * time.Second)

	fmt.Println("Sleeping...")
	sleep := make(chan *sphero.Response)
	s.Sleep(time.Duration(0), 0, 0, sleep)
	slept := <-sleep
	fmt.Printf("Slept %#x\n", slept)

	fmt.Println("Closing...")
	s.Close()
	fmt.Println("Done.")
}
