package main

import (
	"blogaggregator/internal/database"
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

func handlerAgg(s *State, cmd Command, user database.User) error {
	time_between_reqs := cmd.args[1]
	duration, err := time.ParseDuration(time_between_reqs)
	if err != nil {
		return err
	}
	fmt.Printf("Collecting feed every %ss", time_between_reqs)
	ticker := time.NewTicker(duration)
	for ; ; <-ticker.C {
		err := handlerScrapeFeeds(s, user)
		if err != nil {
			return err
		}
	}
}

func handlerAddFeed(s *State, cmd Command, user database.User) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("not enough arguments")
	}
	ctx := context.Background()
	userID := user.ID
	userArgs := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.args[0],
		Url:       cmd.args[1],
		UserID:    userID,
	}
	feed, err := s.db.CreateFeed(ctx, userArgs)
	if err != nil {
		return err
	}

	// Create a new command with just the URL for the follow operation
	followCmd := Command{
		name: "follow",
		args: []string{cmd.args[1]}, // Just the URL
	}
	err = handlerFollow(s, followCmd, user)
	if err != nil {
		return err
	}
	fmt.Println(feed)
	return nil
}

func handlerFeeds(s *State, cmd Command) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeeds(ctx)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		creator, err := s.db.GetUserByID(ctx, feed.UserID)
		if err != nil {
			return err
		}
		fmt.Println(feed.Name, feed.Url, creator)
	}
	return nil
}

func handlerFollow(s *State, cmd Command, user database.User) error {
	ctx := context.Background()
	feed, err := s.db.GetFeedByURL(ctx, cmd.args[0])
	if err != nil {
		return err
	}

	followArgs := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	}
	followRecord, err := s.db.CreateFeedFollow(ctx, followArgs)
	if err != nil {
		return err
	}
	fmt.Println(followRecord)
	return nil
}

func handlerFollowsForUser(s *State, cmd Command, user database.User) error {
	ctx := context.Background()
	feeds, err := s.db.GetFeedFollowsForUser(ctx, user.ID)
	if err != nil {
		return err
	}
	for _, feed := range feeds {
		fmt.Println(feed.FeedName)
	}

	return nil

}

func handlerUnfollow(s *State, cmd Command, user database.User) error {
	ctx := context.Background()
	feed, err := s.db.GetFeedByURL(ctx, cmd.args[0])
	if err != nil {
		return err
	}

	params := database.UnfollowParams{
		FeedID: feed.ID,
		UserID: user.ID,
	}
	err = s.db.Unfollow(ctx, params)
	if err != nil {
		return err
	}
	return nil
}

func handlerScrapeFeeds(s *State, user database.User) error {
	ctx := context.Background()
	nextFeed, err := s.db.GetNextFeedToFetch(ctx, user.ID)
	if err != nil {
		return err
	}
	markParams := database.MarkFeedFetchedParams{
		ID:     nextFeed.FeedID,
		UserID: user.ID,
	}
	err = s.db.MarkFeedFetched(ctx, markParams)
	if err != nil {
		return err
	}
	feedURL, err := s.db.GetFeedByID(ctx, nextFeed.FeedID)
	if err != nil {
		return err
	}
	rssFeed, err := fetchFeed(ctx, feedURL.Url)
	if err != nil {
		return err
	}
	for _, item := range rssFeed.Channel.Item {
		var desc sql.NullString
		var pubDate sql.NullString

		if item.Description != "" {
			desc = sql.NullString{
				String: item.Description,
				Valid:  true,
			}
		} else {
			desc = sql.NullString{
				String: "No description",
				Valid:  false,
			}
		}

		if item.PubDate != "" {
			pubDate = sql.NullString{
				String: item.PubDate,
				Valid:  true,
			}
		} else {
			pubDate = sql.NullString{
				String: "",
				Valid:  false,
			}
		}

		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: desc,
			PublishedAt: pubDate,
			FeedID:      nextFeed.FeedID,
		}
		err := s.db.CreatePost(ctx, postParams)
		if err != nil {
			return err
		}
	}
	return nil
}

func handlerBrowse(s *State, cmd Command, user database.User) error {
	ctx := context.Background()
	limit := 2
	if len(cmd.args) > 0 && cmd.args[0] != "" {
		n, err := strconv.Atoi(cmd.args[0])
		if err != nil {
			return fmt.Errorf("invalid limit: %w", err)
		}
		limit = n
	}
	posts, err := s.db.GetPostsForUser(ctx, database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	})
	if err != nil {
		return err
	}

	for _, post := range posts {
		published := "Unknown date"
		if post.PublishedAt.Valid {
			published = post.PublishedAt.String
		}

		desc := ""
		if post.Description.Valid {
			desc = post.Description.String
		}
		fmt.Printf("Published %v \n", published)
		fmt.Printf("--- %s ---\n", post.Title)
		fmt.Printf("    %s\n", desc)
		fmt.Printf("Link: %s\n", post.Url)
		fmt.Println("=====================================")

	}

	return nil
}
