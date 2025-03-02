package spiders

import (
	"fmt"
	"time"

	"scraper-service/pkg/models"
)

type StaticSpiderConstructor func(url string, comm *Communicator, waitTime time.Duration) Spider
type DynamicSpiderConstructor func(input chan models.Link, comm *Communicator, waitTime time.Duration) Spider

// --- Static Spider Constructors ---

var StaticSpiderConstructors = map[SpiderType]StaticSpiderConstructor{
	MathomLink: func(url string, comm *Communicator, waitTime time.Duration) Spider {
		return &MathomLinkSpider{
			Communicator: comm,
			URL:          url,
			WaitTime:     waitTime,
		}
	},
}

func NewStaticSpider(spiderType SpiderType, url string, comm *Communicator, waitTime time.Duration) (Spider, error) {
	constructor, ok := StaticSpiderConstructors[spiderType]
	if !ok {
		return nil, fmt.Errorf("no static spider constructor for spider type %d", spiderType)
	}
	return constructor(url, comm, waitTime), nil
}

// --- Dynamic Spider Constructors ---

var DynamicSpiderConstructors = map[SpiderType]DynamicSpiderConstructor{
	MathomOffer: func(input chan models.Link, comm *Communicator, waitTime time.Duration) Spider {
		return &MathomOfferSpider{
			Communicator:     comm,
			InputLinkChannel: input,
			WaitTime:         waitTime,
		}
	},
	BGG: func(input chan models.Link, comm *Communicator, waitTime time.Duration) Spider {
		return &BoardGameGeekSpider{
			Communicator:     comm,
			InputLinkChannel: input,
			WaitTime:         waitTime,
		}
	},
}

func NewDynamicSpider(spiderType SpiderType, input chan models.Link, comm *Communicator, waitTime time.Duration) (Spider, error) {
	constructor, ok := DynamicSpiderConstructors[spiderType]
	if !ok {
		return nil, fmt.Errorf("no dynamic spider constructor for spider type %d", spiderType)
	}
	return constructor(input, comm, waitTime), nil
}
