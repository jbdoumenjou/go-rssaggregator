package database

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/jbdoumenjou/go-rssaggregator/internal/generator"

	"github.com/lib/pq"

	"github.com/google/uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func CreateRandomFeed(t *testing.T) Feed {
	t.Helper()

	user := CreateRandomUser(t)
	require.NotEmpty(t, user)
	require.NotEmpty(t, user.ID)

	feedParams := CreateFeedParams{
		Name:   generator.RandomString(12),
		Url:    fmt.Sprintf("https://%s.%s", generator.RandomString(8), generator.RandomString(3)),
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	feed, err := testQueries.CreateFeed(ctx, feedParams)
	require.NoError(t, err)
	require.NotEmpty(t, feed)
	assert.Equal(t, feedParams.Name, feed.Name)
	assert.Equal(t, feedParams.Url, feed.Url)
	assert.Equal(t, feedParams.UserID, feed.UserID)
	require.NotEmpty(t, feed.CreatedAt)
	require.NotEmpty(t, feed.UpdatedAt)

	return feed
}

func TestQueries_CreateFeed(t *testing.T) {
	feed := CreateRandomFeed(t)
	assert.NotEmpty(t, feed)
}

func TestQueries_CreateFeed_BadUserID(t *testing.T) {
	feedParams := CreateFeedParams{
		Name:   generator.RandomString(12),
		Url:    fmt.Sprintf("https://%s.%s", generator.RandomString(8), generator.RandomString(3)),
		UserID: uuid.NullUUID{UUID: uuid.New(), Valid: true},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	feed, err := testQueries.CreateFeed(ctx, feedParams)
	assert.Error(t, err)
	assert.Empty(t, feed)
	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)
	require.True(t, ok)
	assert.Equal(t, "foreign_key_violation", pqErr.Code.Name())
}

func TestQueries_CreateFeed_ExistingURL(t *testing.T) {
	feed1 := CreateRandomFeed(t)
	require.NotEmpty(t, feed1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	feed2, err := testQueries.CreateFeed(ctx, CreateFeedParams{
		Name:   "name2",
		Url:    feed1.Url,
		UserID: feed1.UserID,
	})

	assert.Error(t, err)
	assert.Empty(t, feed2)
	var pqErr *pq.Error
	ok := errors.As(err, &pqErr)
	require.True(t, ok)
	assert.Equal(t, "unique_violation", pqErr.Code.Name())
}

func TestQueries_ListFeeds(t *testing.T) {
	// Use a separate container to avoid conflicts with the CreateFeed tests
	// TODO: find a more elegant way to do this
	container, err := NewPGContainer()
	require.NoError(t, err)
	defer container.Terminate(context.Background())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := container.DB()

	queries := New(db)
	user, err := queries.CreateUser(ctx, generator.RandomString(6))
	require.NoError(t, err)
	require.NotEmpty(t, user)

	var feeds []Feed
	for i := 0; i < 10; i++ {
		feed, err := queries.CreateFeed(ctx, CreateFeedParams{
			Name:   generator.RandomString(12),
			Url:    generator.RandomURL(6),
			UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
		})
		require.NoError(t, err)
		feeds = append(feeds, feed)
	}

	// Ensure that the feeds are sorted by the most recent first
	sort.Slice(feeds, func(i, j int) bool {
		return feeds[i].UpdatedAt.After(feeds[j].CreatedAt)
	})

	tests := []struct {
		name   string
		params ListFeedsParams
		want   []Feed
	}{
		{
			name:   "offset 0, limit 5",
			params: ListFeedsParams{Limit: 5, Offset: 0},
			want:   feeds[:5],
		},
		{
			name:   "offset 1, limit 5",
			params: ListFeedsParams{Limit: 5, Offset: 1},
			want:   feeds[1:6],
		},
		{
			name:   "offset 1, limit 5",
			params: ListFeedsParams{Limit: 5, Offset: 5},
			want:   feeds[5:10],
		},
		{
			name:   "Out of limit: offset 11, limit 5",
			params: ListFeedsParams{Limit: 5, Offset: 11},
			want:   []Feed{},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			list, err := queries.ListFeeds(ctx, test.params)
			require.NoError(t, err)
			assert.Len(t, list, len(test.want))

			for i, feed := range test.want {
				assert.Equal(t, feed.Name, list[i].Name)
				assert.Equal(t, feed.Url, list[i].Url)
				assert.Equal(t, feed.UserID, list[i].UserID)
				assert.Equal(t, feed.CreatedAt, list[i].CreatedAt)
				assert.Equal(t, feed.UpdatedAt, list[i].UpdatedAt)
			}
		})
	}
}

func TestQueries_GetNextFeedsToFetch(t *testing.T) {
	// Use a separate container to avoid conflicts with the CreateFeed tests
	// TODO: find a more elegant way to do this
	container, err := NewPGContainer()
	require.NoError(t, err)
	defer container.Terminate(context.Background())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	db := container.DB()

	queries := New(db)
	user, err := queries.CreateUser(ctx, generator.RandomString(6))

	var feeds []Feed
	// adds new feeds with last_fetched_at
	for i := 0; i < 3; i++ {
		var feed Feed
		query := `
		INSERT INTO feeds (name, url, user_id, created_at, updated_at, last_fetched_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, url, user_id, created_at, updated_at, last_fetched_at;
		`
		now := time.Now()

		err := db.QueryRowContext(ctx, query,
			generator.RandomString(12),
			generator.RandomURL(10),
			user.ID,
			now,
			now,
			now.Add(time.Duration(i)*time.Second),
		).Scan(
			&feed.ID,
			&feed.Name,
			&feed.Url,
			&feed.UserID,
			&feed.CreatedAt,
			&feed.UpdatedAt,
			&feed.LastFetchedAt)
		require.NoError(t, err)

		feeds = append(feeds, feed)
	}

	// Adds a new feed without last_fetched_at
	feed, err := queries.CreateFeed(ctx, CreateFeedParams{
		Name:   generator.RandomString(12),
		Url:    generator.RandomURL(10),
		UserID: uuid.NullUUID{UUID: user.ID, Valid: true},
	})
	require.NoError(t, err)

	// add before to keep the right order
	feeds = append([]Feed{feed}, feeds...)

	feedsToFetch, err := queries.GetNextFeedsToFetch(context.Background(), 4)
	require.NoError(t, err)

	assert.Len(t, feedsToFetch, 4)
	for i := 0; i < len(feedsToFetch); i++ {
		assert.Equal(t, feeds[i].ID, feedsToFetch[i].ID)
		// the last_fetch_at field is not updated at this time, should be done in a service.
		assert.Equal(t, feeds[i].LastFetchedAt, feedsToFetch[i].LastFetchedAt)
	}
}

func TestQueries_MarkFeedFetched(t *testing.T) {
	feed := CreateRandomFeed(t)
	require.Empty(t, feed.LastFetchedAt)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	start := time.Now()
	err := testQueries.MarkFeedFetched(ctx, feed.ID)
	require.NoError(t, err)

	query := `
	SELECT last_fetched_at
	FROM feeds
	WHERE id = $1;
	`

	var lastFetchedAt *time.Time
	err = testDB.QueryRowContext(ctx, query, feed.ID).Scan(&lastFetchedAt)
	require.NoError(t, err)
	require.NotEmpty(t, lastFetchedAt)
	require.True(t, lastFetchedAt.After(start))
}
