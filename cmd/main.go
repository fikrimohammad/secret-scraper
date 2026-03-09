package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fikrimohammad/secret-scraper/config"
	scraperresthandler "github.com/fikrimohammad/secret-scraper/handler/scraper/rest"
	configstaticrepository "github.com/fikrimohammad/secret-scraper/repository/config/static"
	githubclientrepository "github.com/fikrimohammad/secret-scraper/repository/github/client"
	scraperusecase "github.com/fikrimohammad/secret-scraper/usecase/scraper"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	var (
		githubClientRepository = githubclientrepository.New(cfg)
		configStaticRepository = configstaticrepository.New(cfg)
		scraperUseCase         = scraperusecase.New(configStaticRepository, githubClientRepository)
		scraperRestHandler     = scraperresthandler.New(scraperUseCase)
	)

	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Post("/v1/scraper/scrape_secret", scraperRestHandler.ScrapeSecret)

	go func() {
		if err := app.Listen(":3000"); err != nil {
			log.Panic(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	app.Shutdown()
	log.Println("successfully shutting down app")
}
