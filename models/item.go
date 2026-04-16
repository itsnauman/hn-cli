package models

import "time"

// Item is the raw HN API item model.
type Item struct {
	ID          int    `json:"id"`
	Type        string `json:"type"`
	By          string `json:"by"`
	Time        int64  `json:"time"`
	Text        string `json:"text"`
	URL         string `json:"url"`
	Score       int    `json:"score"`
	Title       string `json:"title"`
	Deleted     bool   `json:"deleted"`
	Dead        bool   `json:"dead"`
	Parent      int    `json:"parent"`
	Kids        []int  `json:"kids"`
	Descendants int    `json:"descendants"`
	Poll        int    `json:"poll"`
	Parts       []int  `json:"parts"`
}

// StoryListItem is the minimal view model for story lists.
type StoryListItem struct {
	ID       int    `toon:"id" json:"id"`
	Title    string `toon:"title" json:"title"`
	Score    int    `toon:"score" json:"score"`
	Comments int    `toon:"comments" json:"comments"`
}

// ItemDetail is the expanded view model for hn item.
type ItemDetail struct {
	ID       int    `toon:"id" json:"id"`
	Type     string `toon:"type" json:"type"`
	By       string `toon:"by" json:"by"`
	Title    string `toon:"title" json:"title"`
	URL      string `toon:"url" json:"url,omitempty"`
	Score    int    `toon:"score" json:"score"`
	Time     string `toon:"time" json:"time"`
	Text     string `toon:"text,omitempty" json:"text,omitempty"`
	Comments int    `toon:"comments" json:"comments"`
}

// CommentView is the view model for flattened comment trees.
type CommentView struct {
	ID    int    `toon:"id" json:"id"`
	By    string `toon:"by" json:"by"`
	Text  string `toon:"text" json:"text"`
	Time  string `toon:"time" json:"time"`
	Level int    `toon:"level" json:"level"`
}

// ToStoryListItem converts an API Item to a minimal story list item.
func (it *Item) ToStoryListItem() StoryListItem {
	return StoryListItem{
		ID:       it.ID,
		Title:    it.Title,
		Score:    it.Score,
		Comments: it.Descendants,
	}
}

// ToItemDetail converts an API Item to a detailed view model.
func (it *Item) ToItemDetail(truncate func(string) string) ItemDetail {
	text := it.Text
	if truncate != nil {
		text = truncate(text)
	}
	return ItemDetail{
		ID:       it.ID,
		Type:     it.Type,
		By:       it.By,
		Title:    it.Title,
		URL:      it.URL,
		Score:    it.Score,
		Time:     formatTime(it.Time),
		Text:     text,
		Comments: it.Descendants,
	}
}

// ToCommentView converts an API Item (comment) to a comment view model.
func (it *Item) ToCommentView(level int, truncate func(string) string) CommentView {
	text := it.Text
	if truncate != nil {
		text = truncate(text)
	}
	return CommentView{
		ID:    it.ID,
		By:    it.By,
		Text:  text,
		Time:  formatTime(it.Time),
		Level: level,
	}
}

func formatTime(unix int64) string {
	if unix == 0 {
		return ""
	}
	return time.Unix(unix, 0).UTC().Format(time.RFC3339)
}

// StoriesOutput wraps story list results with aggregates.
type StoriesOutput struct {
	Type    string          `toon:"type" json:"type"`
	Count   int             `toon:"count" json:"count"`
	Total   int             `toon:"total" json:"total"`
	Stories []StoryListItem `toon:"stories" json:"stories"`
}

// ItemOutput wraps a single item result.
type ItemOutput struct {
	Item ItemDetail `toon:"item" json:"item"`
}

// CommentsOutput wraps comment tree results with aggregates.
type CommentsOutput struct {
	StoryID  int           `toon:"story_id" json:"story_id"`
	Count    int           `toon:"count" json:"count"`
	Depth    int           `toon:"depth" json:"depth"`
	Comments []CommentView `toon:"comments" json:"comments"`
}

// DashboardOutput wraps the root command output.
type DashboardOutput struct {
	Version string          `toon:"version" json:"version"`
	Count   int             `toon:"count" json:"count"`
	Top     []StoryListItem `toon:"top" json:"top"`
}

// UpdatedItemView is a minimal view for recently changed items (mixed types).
type UpdatedItemView struct {
	ID    int    `toon:"id" json:"id"`
	Type  string `toon:"type" json:"type"`
	By    string `toon:"by" json:"by"`
	Title string `toon:"title,omitempty" json:"title,omitempty"`
	Text  string `toon:"text,omitempty" json:"text,omitempty"`
}

// ToUpdatedItemView converts an API Item to an updates view model.
func (it *Item) ToUpdatedItemView(truncate func(string) string) UpdatedItemView {
	v := UpdatedItemView{
		ID:   it.ID,
		Type: it.Type,
		By:   it.By,
	}
	if it.Title != "" {
		v.Title = it.Title
	}
	if it.Text != "" && it.Title == "" {
		text := it.Text
		if truncate != nil {
			text = truncate(text)
		}
		v.Text = text
	}
	return v
}

// UpdatesOutput wraps the updates command output.
type UpdatesOutput struct {
	Count    int               `toon:"count" json:"count"`
	Items    []UpdatedItemView `toon:"items" json:"items"`
	Profiles []string          `toon:"profiles" json:"profiles"`
}
