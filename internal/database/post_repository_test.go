package database

import (
	"context"
	"sort"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jbdoumenjou/go-rssaggregator/internal/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostRepository_CreatePost(t *testing.T) {
	postRepository := NewPostRepository(testDB)
	feed := CreateRandomFeed(t)

	params := CreatePostParams{
		Title:       generator.RandomString(10),
		Url:         generator.RandomURL(5),
		Description: generator.RandomString(50),
		PublishedAt: time.Now().UTC().Add(-time.Hour * 24),
		FeedID: uuid.NullUUID{
			UUID:  feed.ID,
			Valid: true,
		},
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	post, err := postRepository.CreatePost(ctx, params)
	require.NoError(t, err)
	require.NotEmpty(t, post)
	assert.Equal(t, feed.ID, post.FeedID.UUID)
	assert.Equal(t, params.Url, post.Url)
	assert.Equal(t, params.Title, post.Title)
	assert.Equal(t, params.PublishedAt.UTC().Round(time.Microsecond), post.PublishedAt.UTC())
}

func TestPostRepository_GetPostsByUser(t *testing.T) {
	postRepository := NewPostRepository(testDB)
	feed := CreateRandomFeed(t)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var createdPosts []Post

	for i := 0; i < 11; i++ {
		params := CreatePostParams{
			Title:       generator.RandomString(10),
			Url:         generator.RandomURL(5),
			Description: generator.RandomString(50),
			PublishedAt: time.Now().UTC().Round(time.Microsecond),
			FeedID: uuid.NullUUID{
				UUID:  feed.ID,
				Valid: true,
			},
		}

		post, err := postRepository.CreatePost(ctx, params)
		require.NoError(t, err)
		require.NotEmpty(t, post)
		require.NotEmpty(t, post.ID)
		createdPosts = append(createdPosts, post)
	}
	sort.Slice(createdPosts, func(i, j int) bool {
		return createdPosts[i].ID > createdPosts[j].ID
	})

	posts, err := postRepository.GetPostsByUser(ctx, feed.UserID, 10)
	sort.Slice(posts, func(i, j int) bool {
		return posts[i].ID > posts[j].ID
	})

	require.NoError(t, err)
	require.Len(t, posts, 10)
	for i, post := range posts {
		assert.Equal(t, createdPosts[i].ID, post.ID)
		assert.Equal(t, feed.ID, post.FeedID.UUID)
	}
}
