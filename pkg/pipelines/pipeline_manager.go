package pipeline

import (
	"context"
	"scraper-service/pkg/models"
	"scraper-service/pkg/spiders"
)

// PipelineManager will build and run your pipeline.
type PipelineManager struct {
	Spiders             map[string]spiders.Spider
	inputChannels       map[string]chan models.Link // for dynamic spiders
	config              []PipelineNode
	notificationChannel string
}

// NewPipelineManager creates a new manager given a pipeline configuration.
func NewPipelineManager(config []PipelineNode, notificationChannel string) *PipelineManager {
	return &PipelineManager{
		Spiders:             make(map[string]spiders.Spider),
		inputChannels:       make(map[string]chan models.Link),
		config:              config,
		notificationChannel: notificationChannel,
	}
}

// BuildPipeline creates spiders and wires their communication channels.
func (pm *PipelineManager) BuildPipeline() {
	// First pass: create channels for dynamic spiders.
	for _, node := range pm.config {
		// Static spiders (like MathomLink) don't need an input channel.
		if node.Type != spiders.MathomLink {
			pm.inputChannels[node.ID] = make(chan models.Link, 100)
		}
	}

	// Second pass: create each spider and set up its communicator.
	for _, node := range pm.config {
		// Build the list of output channels from this spider.
		var childChannels []chan models.Link
		for _, childID := range node.ChildIDs {
			if ch, ok := pm.inputChannels[childID]; ok {
				childChannels = append(childChannels, ch)
			}
		}
		// Create a communicator with the output channels for this node.
		comm := spiders.NewCommunicator(pm.notificationChannel, len(childChannels))
		// Overwrite the channels set by the factory if needed.
		comm.OutputLinkChannels = childChannels

		// Create the spider
		// TODO create map of dynamic and static spiders
		var sp spiders.Spider
		if node.Type == spiders.MathomLink {
			// Static spider: pass URL and waitTime.
			sp = spiders.NewStaticSpider(node.Type, node.URL, comm, node.WaitTime)
		} else {
			// Dynamic spider: pass the pre-created input channel.
			inputCh := pm.inputChannels[node.ID]
			sp = spiders.NewDynamicSpider(node.Type, inputCh, comm, node.WaitTime)
		}
		pm.Spiders[node.ID] = sp
	}
}

func (pm *PipelineManager) RunAll() {
	for _, sp := range pm.Spiders {
		go sp.Run(context.Background())
	}
}
