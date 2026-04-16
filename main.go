package main

import (
	"fmt"
	"os"
	"strings"

	arg "github.com/alexflint/go-arg"

	"github.com/naumanahmad/hacker-news-cli/api"
	"github.com/naumanahmad/hacker-news-cli/cmd"
	"github.com/naumanahmad/hacker-news-cli/models"
	"github.com/naumanahmad/hacker-news-cli/output"
)

type StoriesCmd struct {
	Type string `arg:"positional" default:"top" help:"story type: top, new, best, ask, show, job"`
}

type ItemCmd struct {
	ID int `arg:"positional,required" help:"item ID"`
}

type CommentsCmd struct {
	ID    int `arg:"positional,required" help:"story ID"`
	Depth int `arg:"--depth" default:"2" help:"comment tree depth"`
}

type UserCmd struct {
	Username string `arg:"positional,required" help:"username"`
}

type UpdatesCmd struct{}

type Args struct {
	Stories  *StoriesCmd  `arg:"subcommand:stories" help:"list stories"`
	Item     *ItemCmd     `arg:"subcommand:item" help:"view item details"`
	Comments *CommentsCmd `arg:"subcommand:comments" help:"view comment tree"`
	User     *UserCmd     `arg:"subcommand:user" help:"view user profile"`
	Updates  *UpdatesCmd  `arg:"subcommand:updates" help:"view recent changes"`

	Output  string `arg:"-o,--output" default:"toon" help:"output format: toon, json"`
	Limit   int    `arg:"-n,--limit" default:"10" help:"number of items to return" placeholder:"N"`
	Fields  string `arg:"--fields" help:"comma-separated list of fields to include"`
	Full    bool   `arg:"--full" help:"don't truncate text fields"`
	Version bool   `arg:"--version" help:"display version and exit"`
}

func (Args) Description() string {
	return "Hacker News CLI"
}

func main() {
	var args Args
	p, err := arg.NewParser(arg.Config{}, &args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	err = p.Parse(os.Args[1:])
	switch {
	case err == arg.ErrHelp:
		p.WriteHelpForSubcommand(os.Stdout, p.SubcommandNames()...)
		os.Exit(0)
	case err != nil:
		fmt.Fprintln(os.Stdout, err.Error())
		os.Exit(1)
	}

	if args.Version {
		fmt.Printf("hn v%s\n", cmd.Version)
		os.Exit(0)
	}
	if errOut := validateArgs(&args); errOut != nil {
		output.RenderError(os.Stdout, errOut, cmd.GetFormat(toGlobalFlags(&args)))
		os.Exit(1)
	}

	gf := toGlobalFlags(&args)
	client := api.NewClient()

	switch {
	case args.Stories != nil:
		os.Exit(cmd.RunStories(args.Stories.Type, client, gf))
	case args.Item != nil:
		os.Exit(cmd.RunItem(args.Item.ID, client, gf))
	case args.Comments != nil:
		os.Exit(cmd.RunComments(args.Comments.ID, args.Comments.Depth, client, gf))
	case args.User != nil:
		os.Exit(cmd.RunUser(args.User.Username, client, gf))
	case args.Updates != nil:
		os.Exit(cmd.RunUpdates(client, gf))
	default:
		// No subcommand — show dashboard
		os.Exit(cmd.RunDashboard(client, gf))
	}
}

func toGlobalFlags(args *Args) *cmd.GlobalFlags {
	gf := &cmd.GlobalFlags{
		Output: args.Output,
		Limit:  args.Limit,
		Full:   args.Full,
	}
	for f := range strings.SplitSeq(args.Fields, ",") {
		f = strings.TrimSpace(f)
		if f != "" {
			gf.Fields = append(gf.Fields, f)
		}
	}
	return gf
}

func validateArgs(args *Args) *models.ErrorOutput {
	if args.Limit < 0 {
		return models.NewValidationError(
			"limit must be zero or greater",
			"use -n 0 for no limit, or pass a positive integer",
		)
	}
	if args.Comments != nil && args.Comments.Depth < 0 {
		return models.NewValidationError(
			"depth must be zero or greater",
			"use --depth 0 to suppress comments, or pass a positive integer",
		)
	}
	return nil
}
