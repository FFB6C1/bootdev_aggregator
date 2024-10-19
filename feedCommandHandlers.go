package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ffb6c1/bootdev_aggregator/internal/database"
	"github.com/ffb6c1/bootdev_aggregator/internal/interaction"
	"github.com/google/uuid"
)

func handlerAgg(_ *state, _ command) error {
	feed, err := interaction.FetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		fmt.Println("handlerAgg err:", err)
		os.Exit(1)
	}
	fmt.Println(feed)
	return nil
}

func handlerAddFeed(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 2 {
		return fmt.Errorf("handlerAddFeed: insufficient arguments")
	}

	feed := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.arguments[0],
		Url:       cmd.arguments[1],
		UserID:    user.ID,
	}

	feedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    user.ID,
	}

	if err := s.db.AddFeed(context.Background(), feed); err != nil {
		return fmt.Errorf("handlerAddFeed: %w", err)
	}

	if _, err := s.db.CreateFeedFollow(context.Background(), feedFollow); err != nil {
		return fmt.Errorf("handlerAddFeed: %w", err)
	}

	fmt.Println(feed)
	return nil
}

func handlerFeeds(s *state, _ command) error {
	users := helperMapUsersUUIDtoNames(s)
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("handlerFeeds: %w", err)
	}
	for _, feed := range feeds {
		fmt.Printf("%s - %s - %s/n", feed.Name, feed.Url, users[feed.UserID])
	}
	return nil
}
