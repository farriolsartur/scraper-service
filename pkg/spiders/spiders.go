package spiders

import (
	"context"
	"time"

	"scraper-service/pkg/models"
)

type SpiderType int

const (
	MathomLink SpiderType = iota
	MathomOffer
	BGG
)

type Spider interface {
	Run(ctx context.Context) error
}

type MathomLinkSpider struct {
	Communicator *Communicator
	URL          string
	WaitTime     time.Duration
}

type MathomOfferSpider struct {
	Communicator     *Communicator
	InputLinkChannel chan models.Link
	WaitTime         time.Duration
}

type BoardGameGeekSpider struct {
	Communicator     *Communicator
	InputLinkChannel chan models.Link
	WaitTime         time.Duration
}

func (s *MathomLinkSpider) Run(ctx context.Context) error {
	return nil
}

func (s *MathomOfferSpider) Run(ctx context.Context) error {
	return nil
}

func (s *BoardGameGeekSpider) Run(ctx context.Context) error {
	return nil
}
