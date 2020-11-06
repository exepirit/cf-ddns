package main

import (
	"github.com/exepirit/cf-ddns/internal/app/web"
	"github.com/exepirit/cf-ddns/internal/repository"
	"log"
)

func main() {
	repository.Set(repository.NewMemory())

	engine := web.New()
	if err := engine.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
