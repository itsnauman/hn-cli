package cmd

import (
	"context"
	"os"

	"github.com/naumanahmad/hacker-news-cli/api"
	"github.com/naumanahmad/hacker-news-cli/models"
	"github.com/naumanahmad/hacker-news-cli/output"
)

func RunUpdates(client *api.Client, gf *GlobalFlags) int {
	ctx := context.Background()
	format := GetFormat(gf)
	limit := gf.Limit
	truncate := GetTruncator(gf)

	data, err := client.FetchUpdates(ctx)
	if err != nil {
		output.RenderError(os.Stdout, models.NewAPIError(err), format)
		return 1
	}

	items, err := client.FetchItems(ctx, data.Items, limit)
	if err != nil {
		output.RenderError(os.Stdout, models.NewAPIError(err), format)
		return 1
	}

	views := make([]models.UpdatedItemView, len(items))
	for i, it := range items {
		views[i] = it.ToUpdatedItemView(truncate)
	}

	profiles := data.Profiles
	if limit > 0 && limit < len(profiles) {
		profiles = profiles[:limit]
	}

	result := models.UpdatesOutput{
		Count:    len(views),
		Items:    views,
		Profiles: profiles,
	}

	if len(views) == 0 {
		result.Items = []models.UpdatedItemView{}
	}
	if len(profiles) == 0 {
		result.Profiles = []string{}
	}

	if len(gf.Fields) > 0 {
		wrapper := map[string]any{
			"count":    result.Count,
			"items":    output.SelectFieldsList(views, gf.Fields),
			"profiles": result.Profiles,
		}
		output.Render(os.Stdout, wrapper, format)
		return 0
	}

	output.Render(os.Stdout, result, format)
	return 0
}
