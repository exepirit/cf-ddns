package main

import (
	"github.com/exepirit/cf-ddns/internal/app/web"
	"github.com/exepirit/cf-ddns/internal/bus"
	"github.com/exepirit/cf-ddns/internal/repository"
	"log"
)

func main() {
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	applicationBus := bus.Get()
	worker, err := makeWorker(cfg)
	if err != nil {
		log.Fatal(err)
	}
	applicationBus.Subscribe(worker)

	repo := bus.RepositoryConsumer{
		DDNSRepository: repository.NewMemory(),
	}
	applicationBus.Subscribe(repo)

	go worker.Run()
	engine := web.New()
	if err := engine.Run(cfg.BindAddress); err != nil {
		log.Fatal(err)
	}
}
