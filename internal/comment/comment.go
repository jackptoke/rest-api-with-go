package comment

import (
	"context"
	"errors"
	"fmt"
)

var (
	ErrFetchingComment = errors.New("error fetching comment")
	ErrNotImplemented  = errors.New("not implemented")
)

// Store - this interface defines all the methods
// that our service needs to operate
type Store interface {
	GetComment(context.Context, string) (Comment, error)
}

// Comment - a representation of the comment
// structure for our service
type Comment struct {
	ID     string
	Slug   string
	Body   string
	Author string
}

// Service - is the struct on which all our logi will be built on top of
type Service struct {
	Store Store
}

// NewService - returns a pointer to a new service
func NewService(store Store) *Service {
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

func (s *Service) UpdateComment(ctx context.Context, cmt Comment) error {
	return ErrNotImplemented
}

func (s *Service) DeleteComment(ctx context.Context, id string) error {
	return ErrNotImplemented
}

func (s *Service) CreateComment(ctx context.Context, cmt Comment) (Comment, error) {
	return Comment{}, ErrNotImplemented
}
