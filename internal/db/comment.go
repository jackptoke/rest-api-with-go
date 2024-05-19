package db

import (
	"context"
	"database/sql"
	"errors"
	"fibre/internal/comment"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"log"
)

type CommentRow struct {
	ID     string
	Slug   sql.NullString
	Body   sql.NullString
	Author sql.NullString
}

func convertCommentRowToComment(row CommentRow) comment.Comment {
	return comment.Comment{
		ID:     row.ID,
		Slug:   row.Slug.String,
		Body:   row.Body.String,
		Author: row.Author.String,
	}
}

func (d *Database) GetComment(
	ctx context.Context,
	uuid string) (comment.Comment, error) {

	var cmtRow CommentRow
	row := d.Client.QueryRowContext(
		ctx,
		`SELECT id, slug, body, author
				FROM comments
				WHERE id = $1;`, uuid,
	)

	err := row.Scan(&cmtRow.ID, &cmtRow.Slug, &cmtRow.Body, &cmtRow.Author)
	if errors.Is(err, sql.ErrNoRows) {
		return comment.Comment{}, fmt.Errorf("error fetching the comment by uuid %v", uuid)
	}

	return convertCommentRowToComment(cmtRow), nil
}

func (d *Database) GetComments(
	ctx context.Context) ([]comment.Comment, error) {

	rows, err := d.Client.QueryContext(
		ctx,
		`SELECT id, slug, body, author FROM comments;`)
	if err != nil {
		return []comment.Comment{}, fmt.Errorf("error fetching the comments: %v", err)
	}

	defer rows.Next()
	var cmtRow CommentRow
	var comments []comment.Comment

	for rows.Next() {
		if err := rows.Scan(&cmtRow.ID, &cmtRow.Slug, &cmtRow.Body, &cmtRow.Author); err != nil {
			log.Fatal(err)
		}
		comments = append(comments, convertCommentRowToComment(cmtRow))
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	return comments, nil
}

func (d *Database) PostComment(
	ctx context.Context,
	cmt comment.Comment) (comment.Comment, error) {
	fmt.Println("Creating a new comment")
	cmt.ID = uuid.NewV4().String()
	query := `INSERT INTO comments(ID, Slug, Body, Author) VALUES (:id, :slug, :body, :author);`
	postRow := CommentRow{
		ID:     cmt.ID,
		Slug:   sql.NullString{String: cmt.Slug, Valid: true},
		Body:   sql.NullString{String: cmt.Body, Valid: true},
		Author: sql.NullString{String: cmt.Author, Valid: true},
	}

	rows, err := d.Client.NamedQueryContext(ctx, query, postRow)
	var newCmt CommentRow
	for rows.Next() {
		if err := rows.Scan(&newCmt.ID, &newCmt.Slug, &newCmt.Body, &newCmt.Author); err != nil {
			log.Fatal("Failed to save the comment. Error: ", err)
		}
	}

	if err != nil {
		return comment.Comment{}, fmt.Errorf("error inserting comment: %v", err)
	}
	if err := rows.Close(); err != nil {
		return comment.Comment{}, fmt.Errorf("error closing rows: %v", err)
	}
	return convertCommentRowToComment(newCmt), nil
}

func (d *Database) DeleteComment(ctx context.Context, id string) error {
	_, err := d.Client.ExecContext(
		ctx,
		`DELETE FROM comments WHERE id = $1;`, id)
	if err != nil {
		return fmt.Errorf("error deleting comment: %v", err)
	}

	return nil
}

func (d *Database) UpdateComment(ctx context.Context, id string, cmt comment.Comment) (comment.Comment, error) {
	cmtRow := CommentRow{
		ID:     id,
		Slug:   sql.NullString{String: cmt.Slug, Valid: true},
		Body:   sql.NullString{String: cmt.Body, Valid: true},
		Author: sql.NullString{String: cmt.Author, Valid: true},
	}

	query := `UPDATE comments SET 
                    slug = :slug,
                    author = :author,
                    body = :body
                    WHERE id = :id;`

	fmt.Println(query)

	row, err := d.Client.NamedQueryContext(ctx, query, cmtRow)

	if err != nil {
		return comment.Comment{}, fmt.Errorf("error updating comment: %v", err)
	}

	if err := row.Close(); err != nil {
		return comment.Comment{}, fmt.Errorf("error closing row: %v", err)
	}

	return convertCommentRowToComment(cmtRow), nil
}
