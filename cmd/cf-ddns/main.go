package main

import (
	"github.com/exepirit/cf-ddns/internal/app/web"
	"github.com/exepirit/cf-ddns/internal/bus"
	"github.com/exepirit/cf-ddns/internal/repository"
	"log"
)

func main() {
	// Load config
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create update worker object
	applicationBus := bus.Get()
	worker, err := makeWorker(cfg)
	if err != nil {
		log.Fatal(err)
	}
	applicationBus.Subscribe(worker)

	// Create DDNS records repository object
	repo := bus.RepositoryConsumer{
		DDNSRepository: repository.NewMemory(),
	}
	applicationBus.Subscribe(repo)

	// Add records from repository to worker
	records := repo.GetAll()
	for _, record := range records {
		worker.AddDomain(record.Domain, record.UpdatePeriod)
	}

	// Start worker and application
	go worker.Run()
	engine := web.New()
	if err := engine.Run(cfg.BindAddress); err != nil {
		log.Fatal(err)
	}
}
