package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/itsnauman/hn-cli/api"
	"github.com/itsnauman/hn-cli/models"
	"github.com/itsnauman/hn-cli/output"
)

func RunComments(storyID int, depth int, client *api.Client, gf *GlobalFlags) int {
	format := GetFormat(gf)
	ctx := context.Background()

	if depth < 0 {
		output.RenderError(os.Stdout, models.NewValidationError(
			"depth must be zero or greater",
			"use --depth 0 to suppress comments, or pass a positive integer",
		), format)
		return 1
	}

	story, err := client.FetchItem(ctx, storyID)
	if err != nil {
		output.RenderError(os.Stdout, models.NewErrorFromFetch("item", fmt.Sprintf("%d", storyID), err), format)
		return 1
	}

	truncate := GetTruncator(gf)
	limit := gf.Limit
	comments, err := fetchCommentTree(ctx, client, story.Kids, 0, depth, truncate)
	if err != nil {
		output.RenderError(os.Stdout, models.NewAPIError(err), format)
		return 1
	}

	if limit > 0 && limit < len(comments) {
		comments = comments[:limit]
	}

	result := models.CommentsOutput{
		StoryID:  storyID,
		Count:    len(comments),
		Depth:    depth,
		Comments: comments,
	}

	if len(comments) == 0 {
		result.Comments = []models.CommentView{}
	}

	if len(gf.Fields) > 0 {
		wrapper := map[string]any{
			"story_id": result.StoryID,
			"count":    result.Count,
			"depth":    result.Depth,
			"comments": output.SelectFieldsList(comments, gf.Fields),
		}
		output.Render(os.Stdout, wrapper, format)
		return 0
	}

	output.Render(os.Stdout, result, format)
	return 0
}

func fetchCommentTree(ctx context.Context, client *api.Client, kids []int, level, maxDepth int, truncate func(string) string) ([]models.CommentView, error) {
	if level >= maxDepth || len(kids) == 0 {
		return nil, nil
	}

	items, err := client.FetchItems(ctx, kids, 0)
	if err != nil {
		return nil, err
	}
	var comments []models.CommentView

	for _, item := range items {
		if item.Type != "comment" || item.Deleted || item.Dead {
			continue
		}
		cv := item.ToCommentView(level, truncate)
		comments = append(comments, cv)
		children, err := fetchCommentTree(ctx, client, item.Kids, level+1, maxDepth, truncate)
		if err != nil {
			return nil, err
		}
		comments = append(comments, children...)
	}
	return comments, nil
}
