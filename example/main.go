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

	fmt.Println("Connected")

	ping := make(chan interface{})
	s.Ping(ping)

	fmt.Println("Waiting...")
	<-time.Tick(10 * time.Second)

	fmt.Println("Closing...")
	s.Close()

	fmt.Println("Done.")
}
