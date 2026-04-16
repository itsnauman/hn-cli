package api

import (
	"context"
	"fmt"
)

var storyEndpoints = map[string]string{
	"top":  "/topstories.json",
	"new":  "/newstories.json",
	"best": "/beststories.json",
	"ask":  "/askstories.json",
	"show": "/showstories.json",
	"job":  "/jobstories.json",
}

// ValidStoryTypes returns the list of valid story type names.
func ValidStoryTypes() []string {
	return []string{"top", "new", "best", "ask", "show", "job"}
}

// IsValidStoryType checks if the given type is a valid story type.
func IsValidStoryType(t string) bool {
	_, ok := storyEndpoints[t]
	return ok
}

// FetchStoryIDs fetches the list of story IDs for the given type.
func (c *Client) FetchStoryIDs(ctx context.Context, storyType string) ([]int, error) {
	endpoint, ok := storyEndpoints[storyType]
	if !ok {
		return nil, fmt.Errorf("unknown story type: %s", storyType)
	}
	var ids []int
	if err := c.Get(ctx, endpoint, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}
