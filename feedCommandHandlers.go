package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/ffb6c1/bootdev_aggregator/internal/database"
	"github.com/ffb6c1/bootdev_aggregator/internal/interaction"
	"github.com/google/uuid"
	"github.com/lib/pq"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.arguments) < 1 {
		fmt.Println("agg: Insufficient arguments. Please provide time between requests.")
		os.Exit(1)
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.arguments[0])
	if err != nil {
		fmt.Println("agg: Invalid time unit. Please use a supported time unit - s, m, or h")
		os.Exit(1)
	}
	fmt.Println("Collecting feeds every", timeBetweenReqs)
	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		helperScrapeFeeds(s)
	}
}

func helperScrapeFeeds(s *state) {
	user, err := s.db.GetUser(context.Background(), s.config.CurrentUserName)
	if err != nil {
		fmt.Println("Could not find current user, please log in.")
		os.Exit(1)
	}
	nextFeedURL, err := s.db.GetNextFeed(context.Background(), user.ID)
	if err != nil {
		fmt.Println("Could not get next feed to fetch:", err)
		os.Exit(1)
	}
	fmt.Println(nextFeedURL)
	markFetchedParams := database.MarkFetchedByURLParams{
		Url:       nextFeedURL,
		UpdatedAt: time.Now(),
	}
	if err := s.db.MarkFetchedByURL(context.Background(), markFetchedParams); err != nil {
		fmt.Println("Could not update feed:", err)
		os.Exit(1)
	}
	feed, err := interaction.FetchFeed(context.Background(), nextFeedURL)
	if err != nil {
		fmt.Println("Could not fetch feed:", err)
		os.Exit(1)
	}
	for _, item := range feed.Channel.Item {
		pubDate, err := helperParseDate(item.PubDate)
		if err != nil {
			fmt.Println("Post not added:", err)
			continue
		}
		feedId, err := s.db.GetFeedIDByURL(context.Background(), nextFeedURL)
		if err != nil {
			fmt.Println("Post not added: could not get feed id -", err)
			continue
		}
		postParams := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			Title:       item.Title,
			Url:         item.Link,
			Description: helperMakeNullString(&item.Description),
			PublishedAt: pubDate,
			FeedID:      feedId,
		}
		if err := s.db.CreatePost(context.Background(), postParams); err != nil {
			if pqErr, ok := err.(*pq.Error); ok {
				if pqErr.Code == "23505" {
					continue
				} else {
					fmt.Println("Post not added:", err)
					fmt.Println(pqErr.Code)
					continue
				}
			} else {
				fmt.Println("Post not added:", err)
				continue
			}
		}
		fmt.Println("Post added: ", item.Title)
	}
}

func helperParseDate(date string) (time.Time, error) {
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC3339,
		time.RFC850,
		time.RFC822Z,
		time.RFC822,
		time.RubyDate,
	}
	for _, format := range formats {
		parsedTime, err := time.Parse(format, date)
		if err == nil {
			return parsedTime, nil
		}
	}
	err := fmt.Errorf("Cannot parse time.")
	return time.Now(), err
}

func helperMakeNullString(item *string) sql.NullString {
	if item == nil {
		return sql.NullString{
			String: "",
			Valid:  false,
		}
	}
	return sql.NullString{
		String: *item,
		Valid:  true,
	}
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

func handlerBrowse(s *state, cmd command, user database.User) error {
	limit := 0
	if len(cmd.arguments) < 1 {
		limit = 2
	} else {
		var err error
		limit, err = strconv.Atoi(cmd.arguments[0])
		if err != nil {
			fmt.Println("Invalid number for limit, defaulting to 2.")
			limit = 2
		}
	}
	params := database.GetPostsForUserParams{
		UserID: user.ID,
		Limit:  int32(limit),
	}
	posts, err := s.db.GetPostsForUser(context.Background(), params)
	if err != nil {
		fmt.Println("Error retrieving posts:", err)
		os.Exit(1)
	}
	for _, post := range posts {
		fmt.Println("-----------")
		fmt.Println(post.Title)
		fmt.Println(post.Description.String)
	}
	return nil
}
