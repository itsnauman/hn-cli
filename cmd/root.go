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

const Version = "0.1.0"

// GlobalFlags holds flags shared across all subcommands.
type GlobalFlags struct {
	Output string
	Limit  int
	Fields []string
	Full   bool
}

func GetFormat(gf *GlobalFlags) output.Format {
	switch strings.ToLower(gf.Output) {
	case "json":
		return output.FormatJSON
	case "toon", "":
		return output.FormatTOON
	default:
		fmt.Fprintf(os.Stderr, "unknown output format %q, using toon\n", gf.Output)
		return output.FormatTOON
	}
}

func GetTruncator(gf *GlobalFlags) func(string) string {
	return output.MakeTruncator(output.DefaultTruncateLen, gf.Full)
}

// RunDashboard shows a compact dashboard with top stories.
func RunDashboard(client *api.Client, gf *GlobalFlags) int {
	ctx := context.Background()
	limit := gf.Limit
	format := GetFormat(gf)

	ids, err := client.FetchStoryIDs(ctx, "top")
	if err != nil {
		output.RenderError(os.Stdout, models.NewAPIError(err), format)
		return 1
	}

	items, err := client.FetchItems(ctx, ids, limit)
	if err != nil {
		output.RenderError(os.Stdout, models.NewAPIError(err), format)
		return 1
	}

	stories := make([]models.StoryListItem, len(items))
	for i, it := range items {
		stories[i] = it.ToStoryListItem()
	}

	result := models.DashboardOutput{
		Version: fmt.Sprintf("hn v%s", Version),
		Count:   len(stories),
		Top:     stories,
	}

	if len(gf.Fields) > 0 {
		wrapper := map[string]any{
			"version": result.Version,
			"count":   result.Count,
			"top":     output.SelectFieldsList(stories, gf.Fields),
		}
		output.Render(os.Stdout, wrapper, format)
		return 0
	}

	output.Render(os.Stdout, result, format)
	return 0
}
