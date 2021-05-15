package main

import (
	"github.com/exepirit/cf-ddns/control"
	"github.com/exepirit/cf-ddns/provider"
	"github.com/exepirit/cf-ddns/source/static"
	"log"
	"time"
)

func main() {
	stubProvider := &provider.Stub{
		Logger: log.Default(),
	}
	source := static.NewSourceFromEnv()

	controller := control.Controller{
		Source:     source,
		Provider:   stubProvider,
		TimePeriod: time.Minute,
	}

	err := controller.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
