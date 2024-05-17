package comment

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrFetchingComment = errors.New("error fetching comment")
	ErrNotImplemented  = errors.New("not implemented")
	ErrPostingComment  = errors.New("error posting comment")
	ErrUpdatingComment = errors.New("error updating comment")
	ErrDeletingComment = errors.New("error deleting comment")
	ErrNoCommentFound  = errors.New("no comment found")
)

// Comment - a representation of the comment
// structure for our service
type Comment struct {
	ID     string `json:"id"`
	Slug   string `json:"slug"`
	Body   string `json:"body"`
	Author string `json:"author"`
}

type CmtStore interface {
	GetComment(context.Context, string) (Comment, error)
	GetComments(context.Context) ([]Comment, error)
	PostComment(context.Context, Comment) (Comment, error)
	UpdateComment(context.Context, string, Comment) (Comment, error)
	DeleteComment(context.Context, string) error
}

// Service - is the struct on which all our logi will be built on top of
type Service struct {
	Store CmtStore
}

// NewService - returns a pointer to a new service
func NewService(store CmtStore) *Service {
	return &Service{
		Store: store,
	}
}

func (s *Service) GetComment(ctx context.Context, id string) (Comment, error) {
	fmt.Println("Retrieving a comment ", id)
	cmt, err := s.Store.GetComment(ctx, id)
	if err != nil {
		fmt.Println("Error retrieving comment: ", err)
		return Comment{}, ErrFetchingComment
	}
	return cmt, nil
}

func (s *Service) GetComments(ctx context.Context) ([]Comment, error) {
	fmt.Println("Retrieving comments")
	cmts, err := s.Store.GetComments(ctx)
	if err != nil {
		fmt.Println("Error retrieving comments: ", err)
		return []Comment{}, ErrFetchingComment
	}
	return cmts, nil
}

func (s *Service) UpdateComment(ctx context.Context, id string, cmt Comment) (Comment, error) {
	cmt, err := s.Store.UpdateComment(ctx, id, cmt)
	if err != nil {
		fmt.Println("Error updating comment: ", err)
		return Comment{}, ErrUpdatingComment
	}
	return cmt, nil
}

func (s *Service) DeleteComment(ctx context.Context, id string) error {
	err := s.Store.DeleteComment(ctx, id)
	if err != nil {
		fmt.Println("Error deleting comment: ", err)
		return ErrDeletingComment
	}
	return nil
}

func (s *Service) PostComment(ctx context.Context, cmt Comment) (Comment, error) {
	insertedComment, err := s.Store.PostComment(ctx, cmt)
	if err != nil {
		return Comment{}, ErrPostingComment
	}
	return insertedComment, nil
}
