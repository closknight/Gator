package main

import (
	"Gator/internal/database"
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

func handlerAgg(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage: agg <time_between_reqs>")
	}
	timeBetweenRequests, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("error: not a valid time between requests")
	}

	log.Printf("Collecting Feeds every %s", timeBetweenRequests)

	ticker := time.NewTicker(timeBetweenRequests)

	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func scrapeFeeds(s *state) {
	feed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("Couldn't find new feeds to fetch")
		return
	}
	log.Println("Fetching feed...")
	scrapeFeed(s.db, feed)
}

func scrapeFeed(db *database.Queries, feed database.Feed) {
	_, err := db.MarkFeedFetched(context.Background(), feed.ID)
	if err != nil {
		log.Printf("could not mark feed %s as read: %v", feed.Name, err)
		return
	}
	feedData, err := fetchFeed(context.Background(), feed.Url)
	if err != nil {
		log.Printf("could not collect feed %s: %v", feed.Name, err)
		return
	}
	for _, item := range feedData.Channel.Item {
		fmt.Printf(" * %s\n", item.Title)
		published_at := sql.NullTime{}
		pubDate, err := time.Parse(time.RFC822, item.PubDate)
		if err == nil {
			published_at = sql.NullTime{Valid: true, Time: pubDate}
		}

		params := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now().UTC(),
			UpdatedAt:   time.Now().UTC(),
			Title:       item.Title,
			Url:         item.Link,
			Description: sql.NullString{String: item.Description, Valid: true},
			PublishedAt: published_at,
			FeedID:      feed.ID,
		}

		_, err = db.CreatePost(context.Background(), params)
		if err != nil {
			log.Printf("error logging in post %s: %v", item.Title, err)
		}

	}
	log.Printf("Feed %s collected, %v posts found", feed.Name, len(feedData.Channel.Item))
}
