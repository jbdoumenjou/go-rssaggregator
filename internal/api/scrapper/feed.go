package scrapper

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jbdoumenjou/go-rssaggregator/internal/database"
)

// RSSFeed represents the structure of an RSS feed.
type RSSFeed struct {
	XMLName xml.Name       `xml:"rss"`
	Channel RSSFeedChannel `xml:"channel"`
}

// RSSFeedChannel represents the structure of an RSS feed channel.
type RSSFeedChannel struct {
	Title       string        `xml:"title"`
	Description string        `xml:"description"`
	Language    string        `xml:"language"`
	Items       []RSSFeedItem `xml:"item"`
}

// RSSFeedItem represents the structure of an RSS feed item.
type RSSFeedItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// FeedStore represents a feedRepository for managing feed data.
type FeedStore interface {
	GetNextFeedsToFetch(ctx context.Context, limit int32) ([]database.Feed, error)
	MarkFeedFetched(ctx context.Context, id uuid.UUID) error
}

type PostRepository interface {
	CreatePost(ctx context.Context, arg database.CreatePostParams) (database.Post, error)
}

// FeedFetcher represents a feed fetcher.
type FeedFetcher struct {
	feedRepository FeedStore
	postRepository PostRepository
	interval       time.Duration
	limit          int32
}

// NewFeedFetcher returns a new feed fetcher.
// It fetches feeds from the feedRepository at the given interval and limits the number of feeds to fetch.
func NewFeedFetcher(feedRepository FeedStore, postRepository PostRepository, limit int32, interval time.Duration) *FeedFetcher {
	return &FeedFetcher{
		feedRepository: feedRepository,
		postRepository: postRepository,
		interval:       interval,
		limit:          limit,
	}
}

// Start starts the feed fetcher.
func (f *FeedFetcher) Start(ctx context.Context) {
	ticker := time.NewTicker(f.interval)
	defer ticker.Stop()

	// Fetch the feeds immediately when starting
	err := f.processFeeds(ctx, f.limit)
	if err != nil {
		log.Printf("error fetching feeds: %v", err)
	}

	for {
		select {
		case <-ticker.C:
			if err := f.processFeeds(ctx, f.limit); err != nil {
				log.Printf("error fetching feeds: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

// processFeeds fetches the next feeds to fetch and marks them as fetched.
func (f *FeedFetcher) processFeeds(ctx context.Context, limit int32) error {
	feeds, err := f.feedRepository.GetNextFeedsToFetch(ctx, limit)
	if err != nil {
		return fmt.Errorf("error getting next feeds to fetch: %w", err)
	}

	// Use a wait group to wait for all feeds to be processed
	var wg sync.WaitGroup

	// Fetch and process all the feeds concurrently
	for _, feed := range feeds {
		wg.Add(1)
		go func(feed database.Feed) {
			defer wg.Done()

			rssFeed, err := fetchRSSFeed(feed.Url)
			if err != nil {
				log.Printf("error fetching rss feed: %v", err)
				return
			}
			log.Printf("Process rss feed: " + rssFeed.Channel.Title)
			ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
			defer cancel()
			rssLayout := "Mon, 02 Jan 2006 15:04:05 -0700"

			for _, item := range rssFeed.Channel.Items {

				pubDate, err := time.Parse(rssLayout, item.PubDate)
				if err != nil {
					log.Printf("error parsing timestamp: %v", err)
					continue
				}

				post, err := f.postRepository.CreatePost(ctx, database.CreatePostParams{
					Title:       item.Title,
					Url:         item.Link,
					Description: item.Description,
					PublishedAt: pubDate,
					FeedID: uuid.NullUUID{
						UUID:  feed.ID,
						Valid: true,
					},
				})
				if err != nil {
					log.Printf("error creating post: %v", err)
					return
				}
				fmt.Println("Create post: " + post.Title)
			}

			if err := f.feedRepository.MarkFeedFetched(ctx, feed.ID); err != nil {
				log.Printf("error marking feed fetched: %v", err)
				return
			}
		}(feed)
	}

	// Wait for all feeds to be processed before moving on to the next iteration
	wg.Wait()

	return nil
}

// fetchRSSFeed fetches data from an RSS feed URL and returns the parsed data in a Go struct.
func fetchRSSFeed(feedURL string) (*RSSFeed, error) {
	// Fetch the RSS feed from the URL
	response, err := http.Get(feedURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Check if the response has a valid XML Content-Type
	contentType := response.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/xml") && !strings.HasPrefix(contentType, "text/xml") {
		return nil, fmt.Errorf("unexpected Content-Type: %s", contentType)
	}

	// Read the response body
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	rssFeed, err := parseFeed(body)
	if err != nil {
		return nil, fmt.Errorf("error parsing feed: %w", err)
	}

	return rssFeed, nil
}

// parseFeed parses the XML data and returns the parsed RSS feed.
func parseFeed(data []byte) (*RSSFeed, error) {
	var rssFeed RSSFeed

	if err := xml.Unmarshal(data, &rssFeed); err != nil {
		return nil, fmt.Errorf("error unmarshalling rss feed: %w", err)
	}

	return &rssFeed, nil
}
