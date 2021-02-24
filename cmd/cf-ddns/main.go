package main

import (
	"github.com/exepirit/cf-ddns/internal/app/web"
	"github.com/exepirit/cf-ddns/internal/bus"
	"github.com/exepirit/cf-ddns/internal/repository"
	"github.com/pkg/errors"
	"log"
)

func main() {
	// Load config
	cfg, err := loadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Create update worker object
	worker, err := makeWorker(cfg)
	if err != nil {
		log.Fatal(err)
	}

	// Create DDNS bindings repository object
	repo := makeDataRepository()
	bus.Get().Subscribe(bus.RepositoryConsumer{
		BindingRepository: repo,
	})

	// Add records from repository to worker
	records := repo.GetAll()
	for _, record := range records {
		worker.AddDomain(record.Domain, record.UpdatePeriod)
	}

	// Start worker and application
	go worker.Run()
	engine := web.New(repo)
	if err := engine.Run(cfg.BindAddress); err != nil {
		log.Fatal(err)
	}
}

func makeDataRepository() repository.BindingRepository {
	repo, err := repository.NewSqlite("cf-ddns.sqlite3")
	if err != nil {
		log.Fatalln(errors.WithMessage(err, "failed to create database"))
	}

	if err := repo.Init(); err != nil {
		log.Fatalln(errors.WithMessage(err, "failed to init database"))
	}

	return repo
}
