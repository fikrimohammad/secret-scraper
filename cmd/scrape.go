package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/fikrimohammad/secret-scraper/config"
	"github.com/fikrimohammad/secret-scraper/model"
	configstaticrepository "github.com/fikrimohammad/secret-scraper/repository/config/static"
	githubclientrepository "github.com/fikrimohammad/secret-scraper/repository/github/client"
	"github.com/fikrimohammad/secret-scraper/usecase"
	scraperusecase "github.com/fikrimohammad/secret-scraper/usecase/scraper"
)

func runScrape(args []string) {
	fs := flag.NewFlagSet("scrape", flag.ExitOnError)
	provider := fs.String("provider", "", "Secret provider (e.g., anthropic, openai)")
	secretType := fs.String("type", "", "Secret type (e.g., anthropic_api_key)")
	iterations := fs.Int("iterations", 0, "Max search iterations (0 = unlimited)")
	limit := fs.Int("limit", 10, "Max results per iteration")
	all := fs.Bool("all", false, "Scrape all configured provider/type combinations")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: secret-scraper scrape [flags]\n\nFlags:\n")
		fs.PrintDefaults()
	}

	if err := fs.Parse(args); err != nil {
		os.Exit(1)
	}

	if !*all && (*provider == "" || *secretType == "") {
		fs.Usage()
		os.Exit(1)
	}

	cfg, err := config.Init()
	if err != nil {
		log.Fatalf("failed to init config: %v", err)
	}

	var (
		githubClientRepository = githubclientrepository.New(cfg)
		configStaticRepository = configstaticrepository.New(cfg)
		scraperUseCase         = scraperusecase.New(configStaticRepository, githubClientRepository)
	)

	type scrapeTarget struct {
		provider   string
		secretType string
	}

	var targets []scrapeTarget
	if *all {
		for _, sc := range cfg.SecretScraper {
			targets = append(targets, scrapeTarget{
				provider:   sc.SecretProvider,
				secretType: sc.SecretType,
			})
		}
	} else {
		targets = []scrapeTarget{{provider: *provider, secretType: *secretType}}
	}

	var allSecrets []model.Secret
	for _, t := range targets {
		log.Printf("scraping %s/%s...", t.provider, t.secretType)

		result, err := scraperUseCase.ScrapeSecret(context.Background(), usecase.ScrapeSecretParams{
			SecretProvider:  model.SecretProvider(t.provider),
			SecretType:      model.SecretType(t.secretType),
			MaxIterations:   *iterations,
			MaxLimitPerIter: *limit,
		})
		if err != nil {
			log.Printf("error scraping %s/%s: %v", t.provider, t.secretType, err)
			continue
		}

		allSecrets = append(allSecrets, result.Data...)
		log.Printf("found %d secrets for %s/%s", len(result.Data), t.provider, t.secretType)
	}

	output, err := json.MarshalIndent(map[string]any{"data": allSecrets}, "", "  ")
	if err != nil {
		log.Fatalf("failed to marshal result: %v", err)
	}

	fmt.Println(string(output))
}
