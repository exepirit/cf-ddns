package main

import (
	"github.com/exepirit/cf-ddns/internal/app/web"
	"github.com/exepirit/cf-ddns/internal/repository"
	"log"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	worker, err := makeWorker(cfg)
	if err != nil {
		log.Fatal(err)
	}

	go worker.Run()
	repository.Set(repository.NewMemory())
	engine := web.New()
	if err := engine.Run(cfg.BindAddress); err != nil {
		log.Fatal(err)
	}
}
