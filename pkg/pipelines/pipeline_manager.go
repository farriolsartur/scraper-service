package pipeline

import (
	"context"
	"fmt"
	"scraper-service/pkg/models"
	"scraper-service/pkg/spiders"
	"time"
)

type PipelineManager struct {
	Spiders             map[string]spiders.Spider
	inputChannels       map[string]chan models.Link // for dynamic spiders
	config              []PipelineNode
	notificationChannel string
}

func NewPipelineManager(config []PipelineNode, notificationChannel string) *PipelineManager {
	return &PipelineManager{
		Spiders:             make(map[string]spiders.Spider),
		inputChannels:       make(map[string]chan models.Link),
		config:              config,
		notificationChannel: notificationChannel,
	}
}

func (pm *PipelineManager) BuildPipeline() error {
	// First pass: create channels for dynamic spiders.
	for _, node := range pm.config {
		if node.Kind == spiders.Dynamic {
			pm.inputChannels[node.ID] = make(chan models.Link, 100)
		}
	}

	// Second pass: create each spider and set up its communicator.
	for _, node := range pm.config {
		var childChannels []chan models.Link
		// Collect channels for designated child spiders.
		for _, childID := range node.ChildIDs {
			if ch, ok := pm.inputChannels[childID]; ok {
				childChannels = append(childChannels, ch)
			}
		}

		// Create a communicator with the number of output channels equal to the number of child spiders.
		comm, err := spiders.NewCommunicator(pm.notificationChannel, len(childChannels))
		if err != nil {
			return fmt.Errorf("failed to create communicator for node %s: %w", node.ID, err)
		}
		comm.OutputLinkChannels = childChannels

		var sp spiders.Spider
		if node.Kind == spiders.Static {
			sp, err = spiders.NewStaticSpider(node.Type, node.URL, comm, (time.Duration(node.WaitTime))*time.Second)
			if err != nil {
				return fmt.Errorf("failed to create static spider for node %s: %w", node.ID, err)
			}
		} else {
			inputCh := pm.inputChannels[node.ID]
			sp, err = spiders.NewDynamicSpider(node.Type, inputCh, comm, (time.Duration(node.WaitTime))*time.Second)
			if err != nil {
				return fmt.Errorf("failed to create dynamic spider for node %s: %w", node.ID, err)
			}
		}
		pm.Spiders[node.ID] = sp
	}
	return nil
}

func (pm *PipelineManager) RunAll() {
	for _, sp := range pm.Spiders {
		go sp.Run(context.Background())
	}
}
