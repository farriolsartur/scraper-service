package spiders

import (
	"context"
	"fmt"

	"scraper-service/pkg/cache"  // Adjust the module path as needed.
	"scraper-service/pkg/models" // Adjust the module path as needed.
)

// Communicator provides a common interface for spiders (static or dynamic) to send links and signal completion.
type Communicator struct {
	MergerCache        *cache.MergerCache
	OutputLinkChannels []chan models.Link
}

func NewCommunicator(notificationChannel string, numOutputChannels int) (*Communicator, error) {
	mergerCache, err := cache.NewMergerCache(notificationChannel)
	if err != nil {
		return nil, fmt.Errorf("failed to create MergerCache: %w", err)
	}

	channels := make([]chan models.Link, numOutputChannels)
	for i := range channels {
		channels[i] = make(chan models.Link, 100) // buffered channel with capacity 100; adjust as needed.
	}

	return &Communicator{
		MergerCache:        mergerCache,
		OutputLinkChannels: channels,
	}, nil
}

// SendLink sends a link to all output channels and increments the Redis counter for the item.
// The counter key is derived from the provided itemID. It increments by the number of channels.
func (c *Communicator) SendLink(ctx context.Context, link models.Link) error {
	counterKey := fmt.Sprintf("counter:%s", link.ItemID)
	// Increase the counter by the number of output channels.
	if _, err := c.MergerCache.IncrementCounter(ctx, counterKey, int64(len(c.OutputLinkChannels))); err != nil {
		return fmt.Errorf("failed to increment counter: %w", err)
	}

	// Send the link to every output channel.
	for _, ch := range c.OutputLinkChannels {
		ch <- link
	}
	return nil
}

// FinishProcessing should be called by a spider when it has finished processing an item.
// It decrements the Redis counter for the given item (using the merger cache).
// If the counter reaches zero, the merger will be notified.
func (c *Communicator) FinishProcessing(ctx context.Context, itemID string, notificationMessage string) error {
	counterKey := fmt.Sprintf("counter:%s", itemID)
	if _, err := c.MergerCache.DecrementCounter(ctx, counterKey, notificationMessage); err != nil {
		return fmt.Errorf("failed to decrement counter: %w", err)
	}
	return nil
}
