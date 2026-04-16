package api

import (
	"context"
	"fmt"

	"github.com/itsnauman/hacker-news-cli/models"
)

// FetchUser fetches a user profile by username.
func (c *Client) FetchUser(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	path := fmt.Sprintf("/user/%s.json", username)
	if err := c.Get(ctx, path, &user); err != nil {
		return nil, err
	}
	if user.ID == "" {
		return nil, fmt.Errorf("not found: user %s", username)
	}
	return &user, nil
}
