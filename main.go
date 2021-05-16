package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/exepirit/cf-ddns/control"
	"github.com/exepirit/cf-ddns/provider"
	"github.com/exepirit/cf-ddns/source"
)

func main() {
	stubProvider := &provider.Stub{
		Logger: log.Default(),
	}

	src, err := source.NewFromConfig(&source.Config{
		SourceType: "static",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	controller := control.Controller{
		Source:     src,
		Provider:   stubProvider,
		TimePeriod: time.Minute,
	}

	err = controller.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
