package pipeline

import (
	"scraper-service/pkg/spiders"
	"time"
)

type PipelineNode struct {
	ID       string
	Type     spiders.SpiderType
	Kind     spiders.SpiderKind // Explicitly indicates if the spider is static or dynamic.
	URL      string             // Only used for static spiders.
	WaitTime time.Duration
	ParentID string   // For dynamic spiders: the ID of the spider that sends links.
	ChildIDs []string // IDs of spiders that should receive links from this spider.
}
