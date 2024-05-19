//go:build integration
// +build integration

package db

import (
	"context"
	"fibre/internal/comment"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommentDatabase(t *testing.T) {

	t.Run("test create comment", func(t *testing.T) {
		db, err := NewTestDatabase()
		assert.NoError(t, err)

		cmt, err := db.PostComment(context.Background(), comment.Comment{
			Slug:   "My New Comment",
			Body:   "This is a test comment",
			Author: "jack",
		})

		assert.NoError(t, err)

		fmt.Printf("Comment: %+v, %+v\n", cmt.ID, cmt.Slug)
	})

	t.Run("test retrieving a comment", func(t *testing.T) {
		db, err := NewTestDatabase()

		assert.NoError(t, err)

		id := "88d3872b-1c12-4cbd-8e2f-ebf2b0f06cd6"

		newCmt, err := db.GetComment(context.Background(), id)
		assert.NoError(t, err)
		assert.Equal(t, id, newCmt.ID)

		fmt.Printf("Comment: %+v = %+v\n", newCmt.ID, newCmt.Slug)
	})

	t.Run("test updating a comment", func(t *testing.T) {
		db, err := NewTestDatabase()

		assert.NoError(t, err)

		id := "88d3872b-1c12-4cbd-8e2f-ebf2b0f06cd6"

		cmt := comment.Comment{
			ID:   id,
			Slug: "My New Comment",
			Body: "This is a test comment",
		}

		newCmt, err := db.UpdateComment(context.Background(), id, cmt)

		assert.NoError(t, err)
		assert.Equal(t, newCmt.Body, cmt.Body)
	})

	t.Run("test deleting a comment", func(t *testing.T) {
		db, err := NewTestDatabase()

		assert.NoError(t, err)

		id := "88d3872b-1c12-4cbd-8e2f-ebf2b0f06cd6"

		err = db.DeleteComment(context.Background(), id)

		assert.NoError(t, err)

		_, err = db.GetComment(context.Background(), id)
		assert.Error(t, err)
	})
}
