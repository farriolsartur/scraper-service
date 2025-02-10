package spiders

import (
	"time"

	"scraper-service/pkg/models"
)

// NewStaticSpider creates a new static spider with optional arguments.
func NewStaticSpider(spiderType SpiderType, url string, comm *Communicator, args ...interface{}) Spider {
	var waitTime time.Duration

	// Extract waitTime if provided
	for _, arg := range args {
		if v, ok := arg.(time.Duration); ok {
			waitTime = v
		}
	}

	switch spiderType {
	case MathomLink:
		return &MathomLinkSpider{
			Communicator: comm,
			URL:          url,
			WaitTime:     waitTime,
		}
	default:
		return nil
	}
}

// NewDynamicSpider creates a new dynamic spider with optional arguments.
func NewDynamicSpider(spiderType SpiderType, input chan models.Link, comm *Communicator, args ...interface{}) Spider {
	var waitTime time.Duration

	// Extract waitTime if provided
	for _, arg := range args {
		if v, ok := arg.(time.Duration); ok {
			waitTime = v
		}
	}

	switch spiderType {
	case MathomOffer:
		return &MathomOfferSpider{
			Communicator:     comm,
			InputLinkChannel: input,
			WaitTime:         waitTime,
		}
	case BGG:
		return &BoardGameGeekSpider{
			Communicator:     comm,
			InputLinkChannel: input,
			WaitTime:         waitTime,
		}
	default:
		return nil
	}
}
