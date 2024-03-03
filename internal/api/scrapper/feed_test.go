package scrapper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_parseFeed(t *testing.T) {
	content, err := os.ReadFile("testdata/feed.xml")
	require.NoError(t, err)
	require.NotEmpty(t, content)

	rssFeed, err := parseFeed(content)
	require.NoError(t, err)
	require.Equal(t, "Boot.dev Blog", rssFeed.Channel.Title)
	require.Equal(t, "Recent content on Boot.dev Blog", rssFeed.Channel.Description)
	require.Equal(t, "en-us", rssFeed.Channel.Language)

	require.Len(t, rssFeed.Channel.Items, 2)
	item := rssFeed.Channel.Items[1]
	assert.Equal(t, "The Boot.dev Beat. February 2024", item.Title)

	assert.Equal(t, "Wed, 31 Jan 2024 00:00:00 +0000", item.PubDate)
	assert.Equal(t, "https://blog.boot.dev/news/bootdev-beat-2024-02/", item.Link)
	assert.Equal(t, `609,179. That&rsquo;s the number of lessons you crazy folks have completed on Boot.dev in the last 30 days.`, item.Description)
}
