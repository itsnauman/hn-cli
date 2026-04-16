package api

import (
	"context"
	"fmt"
	"sync"

	"github.com/naumanahmad/hacker-news-cli/models"
)

// FetchItem fetches a single item by ID.
func (c *Client) FetchItem(ctx context.Context, id int) (*models.Item, error) {
	var item models.Item
	path := fmt.Sprintf("/item/%d.json", id)
	if err := c.Get(ctx, path, &item); err != nil {
		return nil, err
	}
	if item.ID == 0 {
		return nil, fmt.Errorf("not found: item %d", id)
	}
	return &item, nil
}

// FetchItems fetches multiple items concurrently using a 20-worker semaphore pool.
// If limit > 0, only the first `limit` IDs are fetched.
// Preserves the order of the input IDs. Skips deleted/dead/failed items.
// Returns an error only if all fetches fail.
func (c *Client) FetchItems(ctx context.Context, ids []int, limit int) ([]*models.Item, error) {
	if limit < 0 {
		return nil, fmt.Errorf("limit must be zero or greater")
	}
	if limit > 0 && limit < len(ids) {
		ids = ids[:limit]
	}
	if len(ids) == 0 {
		return nil, nil
	}

	results := make([]*models.Item, len(ids))
	sem := make(chan struct{}, 20)
	errs := make(chan error, len(ids))
	var wg sync.WaitGroup

	for i, id := range ids {
		wg.Add(1)
		go func(idx, itemID int) {
			defer wg.Done()
			select {
			case sem <- struct{}{}:
				defer func() { <-sem }()
			case <-ctx.Done():
				errs <- fmt.Errorf("item %d: %w", itemID, ctx.Err())
				return
			}

			item, err := c.FetchItem(ctx, itemID)
			if err != nil {
				errs <- fmt.Errorf("item %d: %w", itemID, err)
				return
			}
			results[idx] = item
		}(i, id)
	}
	wg.Wait()
	close(errs)

	var (
		failed   int
		firstErr error
	)
	for err := range errs {
		failed++
		if firstErr == nil {
			firstErr = err
		}
	}

	compact := make([]*models.Item, 0, len(results))
	for _, item := range results {
		if item != nil && !item.Deleted && !item.Dead {
			compact = append(compact, item)
		}
	}
	if failed > 0 {
		return nil, fmt.Errorf("failed to fetch %d of %d items: %w", failed, len(ids), firstErr)
	}
	return compact, nil
}
