package static

import (
	"fmt"

	"github.com/fikrimohammad/secret-scraper/config"
	"github.com/fikrimohammad/secret-scraper/model"
	"github.com/fikrimohammad/secret-scraper/repository"
)

type repositoryObject struct {
	config map[string]*config.SecretScraper
}

func New(cfg *config.Config) repository.ConfigStaticRepository {
	configMap := map[string]*config.SecretScraper{}
	for _, scfg := range cfg.SecretScraper {
		configKey := buildConfigKey(model.SecretProvider(scfg.SecretProvider), model.SecretType(scfg.SecretType))
		configMap[configKey] = &scfg
	}

	return &repositoryObject{
		config: configMap,
	}
}

func buildConfigKey(secretProvider model.SecretProvider, secretType model.SecretType) string {
	return fmt.Sprintf("%s||%s", secretProvider, secretType)
}
