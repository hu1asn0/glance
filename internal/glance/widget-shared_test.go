package glance

import (
	"testing"
	"time"
)

func TestCalculateEngagement(t *testing.T) {
	t.Run("basic engagement scoring", func(t *testing.T) {
		posts := forumPostList{
			{CommentCount: 100, Score: 200, TimePosted: time.Now()},
			{CommentCount: 50, Score: 100, TimePosted: time.Now()},
			{CommentCount: 10, Score: 20, TimePosted: time.Now()},
		}

		posts.calculateEngagement()

		// Post with highest comments and score should have highest engagement
		if posts[0].Engagement <= posts[1].Engagement {
			t.Errorf("post 0 (100 comments, 200 score) should have higher engagement than post 1 (50, 100): %v vs %v",
				posts[0].Engagement, posts[1].Engagement)
		}
		if posts[1].Engagement <= posts[2].Engagement {
			t.Errorf("post 1 should have higher engagement than post 2: %v vs %v",
				posts[1].Engagement, posts[2].Engagement)
		}
	})

	t.Run("depreciation for old posts", func(t *testing.T) {
		now := time.Now()
		posts := forumPostList{
			{CommentCount: 100, Score: 200, TimePosted: now},
			{CommentCount: 100, Score: 200, TimePosted: now.Add(-12 * time.Hour)},
		}

		posts.calculateEngagement()

		// Same stats but the older post should have lower engagement due to depreciation
		if posts[1].Engagement >= posts[0].Engagement {
			t.Errorf("old post should have lower engagement due to depreciation: fresh=%v old=%v",
				posts[0].Engagement, posts[1].Engagement)
		}
	})

	t.Run("recent posts not depreciated", func(t *testing.T) {
		now := time.Now()
		posts := forumPostList{
			{CommentCount: 100, Score: 200, TimePosted: now},
			{CommentCount: 100, Score: 200, TimePosted: now.Add(-3 * time.Hour)},
		}

		posts.calculateEngagement()

		// Both are within the depreciatePostsOlderThanHours window, should be equal
		if posts[0].Engagement != posts[1].Engagement {
			t.Errorf("posts within %dh window should have equal engagement: %v vs %v",
				depreciatePostsOlderThanHours, posts[0].Engagement, posts[1].Engagement)
		}
	})

	t.Run("single post", func(t *testing.T) {
		posts := forumPostList{
			{CommentCount: 50, Score: 100, TimePosted: time.Now()},
		}

		posts.calculateEngagement()

		// Single post: both ratios = 1.0, engagement = 1.0
		if posts[0].Engagement != 1.0 {
			t.Errorf("single post engagement should be 1.0, got %v", posts[0].Engagement)
		}
	})
}

func TestSortByEngagement(t *testing.T) {
	posts := forumPostList{
		{Title: "low", Engagement: 0.5},
		{Title: "high", Engagement: 2.0},
		{Title: "mid", Engagement: 1.0},
	}

	posts.sortByEngagement()

	if posts[0].Title != "high" || posts[1].Title != "mid" || posts[2].Title != "low" {
		t.Errorf("expected [high, mid, low], got [%s, %s, %s]",
			posts[0].Title, posts[1].Title, posts[2].Title)
	}
}
