package main

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}

	req.Header.Set("gator", "1.0")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer res.Body.Close()

	dat, err := io.ReadAll(res.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	rssResp := RSSFeed{}
	err = xml.Unmarshal(dat, &rssResp)
	if err != nil {
		return &RSSFeed{}, err
	}

	rssResp.Channel.Title = html.UnescapeString(rssResp.Channel.Title)
	rssResp.Channel.Description = html.UnescapeString(rssResp.Channel.Description)
	for _, i := range rssResp.Channel.Item {
		i.Title = html.UnescapeString(i.Title)
		i.Description = html.EscapeString((i.Description))
	}

	return &rssResp, nil

}
