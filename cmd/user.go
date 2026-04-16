package cmd

import (
	"context"
	"os"

	"github.com/itsnauman/hn-cli/api"
	"github.com/itsnauman/hn-cli/models"
	"github.com/itsnauman/hn-cli/output"
)

func RunUser(username string, client *api.Client, gf *GlobalFlags) int {
	format := GetFormat(gf)
	ctx := context.Background()

	user, err := client.FetchUser(ctx, username)
	if err != nil {
		output.RenderError(os.Stdout, models.NewErrorFromFetch("user", username, err), format)
		return 1
	}

	view := user.ToUserView(GetTruncator(gf))
	submitted := len(user.Submitted)

	if len(gf.Fields) > 0 {
		filtered := output.SelectFields(view, gf.Fields)
		output.Render(os.Stdout, filtered, format)
		return 0
	}

	result := models.UserOutput{
		User:      view,
		Submitted: submitted,
	}

	output.Render(os.Stdout, result, format)
	return 0
}
