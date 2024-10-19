package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ffb6c1/bootdev_aggregator/internal/database"
	"github.com/google/uuid"
)

func handlerFollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		fmt.Println("handlerFollow: Insufficient arguments.")
		os.Exit(1)
	}

	feedID, err := s.db.GetFeedIDByURL(context.Background(), cmd.arguments[0])
	if err != nil {
		fmt.Println("handlerFollow: feed not in feeds - use add to add a feed before following.")
	}

	newFeedFollow := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		FeedID:    feedID,
	}

	feedFollow, err := s.db.CreateFeedFollow(context.Background(), newFeedFollow)
	if err != nil {
		fmt.Println("handlerFollow:", err)
		os.Exit(1)
	}
	fmt.Printf("%s - %s\n", feedFollow.FeedName, feedFollow.UserName)
	return nil
}

func handlerFollowing(s *state, _ command, user database.User) error {

	following, err := s.db.GetFollowsByUserID(context.Background(), user.ID)
	if err != nil {
		fmt.Println("handlerFollowing:", err)
		os.Exit(1)
	}

	fmt.Println(" -- following:")
	for _, item := range following {
		fmt.Println(item.FeedName)
	}
	return nil
}

func handlerUnfollow(s *state, cmd command, user database.User) error {
	if len(cmd.arguments) < 1 {
		fmt.Println("Please specify the URL of the feed you want to unfollow.")
	}
	feed, err := s.db.GetFeedIDByURL(context.Background(), cmd.arguments[0])
	if err != nil {
		fmt.Println("Could not find feed at provided url:", err)
		os.Exit(1)
	}
	unfollowParams := database.UnfollowParams{
		UserID: user.ID,
		FeedID: feed,
	}

	if err := s.db.Unfollow(context.Background(), unfollowParams); err != nil {
		fmt.Println("Could not unfollow:", err)
		os.Exit(1)
	}
	fmt.Println("Unfollowed feed at ", cmd.arguments[0])
	return nil
}
