# hn

A Hacker News CLI built for AI agents. Fast, minimal output, two formats.

Defaults to [TOON](https://github.com/toon-format/toon) output, which uses ~40% fewer tokens than JSON. Pass `-o json` when you need JSON.

## Install

```bash
go install github.com/itsnauman/hn-cli@latest
```

Or build it yourself:

```bash
git clone https://github.com/itsnauman/hn-cli.git
cd hacker-news-cli
go build -o hn .
```

## Usage

```bash
hn stories              # top stories (default)
hn stories ask          # Ask HN
hn stories show         # Show HN
hn stories job          # jobs
hn item 47769089        # single item
hn comments 47769089    # comment tree
hn user dang            # user profile
hn updates              # recently changed items
```

## Output

TOON (default):

```
type: top
count: 3
total: 500
stories[#3]{id,title,score,comments}:
  47769089,Cybersecurity looks like proof of work now,250,96
  47786164,I made a terminal pager,63,13
  47785948,Ohio prison inmates 'built computers and hid them in ceiling (2017),67,48
```

JSON (`-o json`):

```json
{
  "type": "top",
  "count": 2,
  "total": 500,
  "stories": [
    {
      "id": 47769089,
      "title": "Cybersecurity looks like proof of work now",
      "score": 250,
      "comments": 96
    },
    {
      "id": 47786164,
      "title": "I made a terminal pager",
      "score": 63,
      "comments": 13
    }
  ]
}
```

## Commands

**`hn stories [type]`** -- `top` (default), `new`, `best`, `ask`, `show`, `job`

```bash
hn stories new -n 5
hn stories best -o json
```

**`hn item <id>`** -- single item (story, comment, job, poll)

```bash
hn item 47769089
hn item 47769089 --full              # full text, no truncation
hn item 47769089 --fields id,title   # only these fields
```

**`hn comments <id>`** -- comment tree for a story

```bash
hn comments 47769089 --depth 1    # top-level only
hn comments 47769089 --depth 3    # three levels deep
```

**`hn user <username>`** -- user profile

```bash
hn user dang
hn user dang --full
```

**`hn updates`** -- recently changed items and profiles

```bash
hn updates -n 5
```

## Flags

| Flag | Short | Default | |
|------|-------|---------|-|
| `--output` | `-o` | `toon` | Output format: `toon` or `json` |
| `--limit` | `-n` | `10` | Number of items |
| `--fields` | | all | Comma-separated field list |
| `--full` | | off | Show full text, skip truncation |

## Field selection

Use `--fields` to cut the output down to what you need:

```bash
hn stories -n 5 --fields id,title
```

```
count: 5
stories[#5]{id,title}:
  47769089,Cybersecurity looks like proof of work now
  47786164,I made a terminal pager
  47785948,Ohio prison inmates 'built computers and hid them in ceiling (2017)
  47786791,YouTube now lets you turn off Shorts
  47782570,Google broke its promise to me – now ICE has my data
total: 500
type: top
```

## Design notes

Built with agents in mind, following the [AXI](https://axi.md/) principles.

List items have 3-4 fields by default instead of the full 10+ from the API. Text fields are truncated to 300 chars with a character count so agents know how much they're missing. Every list includes `count` and `total` to avoid extra round trips. Errors come back as structured objects with `error`, `code`, and `hint` fields. Nothing is interactive. Exit codes are 0 or 1.

Items are fetched in parallel (20 concurrent requests), so even a `--limit 50` comes back in about a second.

## API

Wraps the [Hacker News API](https://github.com/HackerNews/API). Two dependencies: [toon-go](https://github.com/toon-format/toon-go) for TOON output, [go-arg](https://github.com/alexflint/go-arg) for argument parsing. Everything else is standard library.

## License

MIT
