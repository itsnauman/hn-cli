package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/itsnauman/hn-cli/api"
	"github.com/itsnauman/hn-cli/models"
	"github.com/itsnauman/hn-cli/output"
)

func RunItem(id int, client *api.Client, gf *GlobalFlags) int {
	format := GetFormat(gf)
	ctx := context.Background()

	item, err := client.FetchItem(ctx, id)
	if err != nil {
		output.RenderError(os.Stdout, models.NewErrorFromFetch("item", fmt.Sprintf("%d", id), err), format)
		return 1
	}

	detail := item.ToItemDetail(GetTruncator(gf))

	if len(gf.Fields) > 0 {
		filtered := output.SelectFields(detail, gf.Fields)
		output.Render(os.Stdout, filtered, format)
		return 0
	}

	result := models.ItemOutput{Item: detail}
	output.Render(os.Stdout, result, format)
	return 0
}
