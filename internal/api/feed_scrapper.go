package api

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strings"
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
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// FetchRSSFeed fetches data from an RSS feed URL and returns the parsed data in a Go struct.
func FetchRSSFeed(feedURL string) (*RSSFeed, error) {
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
