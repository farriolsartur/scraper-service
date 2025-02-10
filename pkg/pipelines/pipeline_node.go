package pipeline

import (
	"scraper-service/pkg/spiders"
	"time"
)

type PipelineNode struct {
	ID       string
	Type     spiders.SpiderType
	URL      string // Only for static spiders.
	WaitTime time.Duration
	ParentID string   // For dynamic spiders, the ID of the spider that sends them links.
	ChildIDs []string // IDs of spiders that should receive links from this spider.
}
