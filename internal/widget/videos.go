package widget

import (
	"context"
	"html/template"
	"sort"
	"time"

	"github.com/glanceapp/glance/internal/assets"
	"github.com/glanceapp/glance/internal/feed"
)

type Videos struct {
	widgetBase        `yaml:",inline"`
	Videos            feed.Videos `yaml:"-"`
	VideoUrlTemplate  string      `yaml:"video-url-template"`
	Style             string      `yaml:"style"`
	CollapseAfterRows int         `yaml:"collapse-after-rows"`
	Channels          []string    `yaml:"channels"`
	BilibiliUIDs      []int       `yaml:"bilibili-uids"`
	Limit             int         `yaml:"limit"`
}

func (widget *Videos) Initialize() error {
	widget.withTitle("Videos").withCacheDuration(time.Hour)

	if widget.Limit <= 0 {
		widget.Limit = 25
	}

	if widget.CollapseAfterRows == 0 || widget.CollapseAfterRows < -1 {
		widget.CollapseAfterRows = 4
	}

	return nil
}

func (widget *Videos) Update(ctx context.Context) {
	bilibiliVideos, err := feed.FetchBilibiliUploads(widget.BilibiliUIDs)

	if err != nil {
		return
	}
	youtubeVideos, err := feed.FetchYoutubeChannelUploads(widget.Channels, widget.VideoUrlTemplate)
	if err != nil {
		return
	}

	videos := append(bilibiliVideos, youtubeVideos...)
	sort.Slice(videos, func(i, j int) bool {
		return videos[i].TimePosted.After(videos[j].TimePosted)
	})

	if !widget.canContinueUpdateAfterHandlingErr(err) {
		return
	}

	if len(videos) > widget.Limit {
		videos = videos[:widget.Limit]
	}

	widget.Videos = videos
}

func (widget *Videos) Render() template.HTML {
	if widget.Style == "grid-cards" {
		return widget.render(widget, assets.VideosGridTemplate)
	}

	return widget.render(widget, assets.VideosTemplate)
}
