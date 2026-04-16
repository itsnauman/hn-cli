package api

import (
	"context"
)

// UpdatesData is the raw API response for /updates.json.
type UpdatesData struct {
	Items    []int    `json:"items"`
	Profiles []string `json:"profiles"`
}

// FetchUpdates fetches recent item and profile changes.
func (c *Client) FetchUpdates(ctx context.Context) (*UpdatesData, error) {
	var data UpdatesData
	if err := c.Get(ctx, "/updates.json", &data); err != nil {
		return nil, err
	}
	return &data, nil
}
