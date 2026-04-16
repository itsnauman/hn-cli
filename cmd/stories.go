package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/itsnauman/hn-cli/api"
	"github.com/itsnauman/hn-cli/models"
	"github.com/itsnauman/hn-cli/output"
)

func RunStories(storyType string, client *api.Client, gf *GlobalFlags) int {
	ctx := context.Background()
	format := GetFormat(gf)
	limit := gf.Limit

	if !api.IsValidStoryType(storyType) {
		output.RenderError(os.Stdout, models.NewValidationError(
			fmt.Sprintf("unknown story type: %s", storyType),
			fmt.Sprintf("valid types: %s", strings.Join(api.ValidStoryTypes(), ", ")),
		), format)
		return 1
	}

	ids, err := client.FetchStoryIDs(ctx, storyType)
	if err != nil {
		output.RenderError(os.Stdout, models.NewAPIError(err), format)
		return 1
	}

	total := len(ids)
	items, err := client.FetchItems(ctx, ids, limit)
	if err != nil {
		output.RenderError(os.Stdout, models.NewAPIError(err), format)
		return 1
	}

	stories := make([]models.StoryListItem, len(items))
	for i, it := range items {
		stories[i] = it.ToStoryListItem()
	}

	if len(gf.Fields) > 0 {
		filtered := output.SelectFieldsList(stories, gf.Fields)
		wrapper := map[string]any{
			"type":    storyType,
			"count":   len(filtered),
			"total":   total,
			"stories": filtered,
		}
		output.Render(os.Stdout, wrapper, format)
		return 0
	}

	result := models.StoriesOutput{
		Type:    storyType,
		Count:   len(stories),
		Total:   total,
		Stories: stories,
	}

	output.Render(os.Stdout, result, format)
	return 0
}
