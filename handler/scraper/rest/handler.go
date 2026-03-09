package rest

import (
	"github.com/fikrimohammad/secret-scraper/handler"
	"github.com/fikrimohammad/secret-scraper/usecase"
)

type handlerObject struct {
	scraperUseCase usecase.Scraper
}

func New(scraperUseCase usecase.Scraper) handler.ScraperREST {
	return &handlerObject{
		scraperUseCase: scraperUseCase,
	}
}
