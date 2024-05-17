package http

import (
	"context"
	"encoding/json"
	"fibre/internal/comment"
	"github.com/go-playground/validator/v10"
	_ "github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2/log"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"net/http"
)

type CommentService interface {
	PostComment(ctx context.Context, comment comment.Comment) (comment.Comment, error)
	GetComment(ctx context.Context, id string) (comment.Comment, error)
	GetComments(ctx context.Context) ([]comment.Comment, error)
	UpdateComment(ctx context.Context, id string, comment comment.Comment) (comment.Comment, error)
	DeleteComment(ctx context.Context, id string) error
}

type Response struct {
	Message string `json:"message"`
}

// PostCommentRequest JSON validation
type PostCommentRequest struct {
	Slug   string `json:"slug" validate:"required"`
	Author string `json:"author" validate:"required"`
	Body   string `json:"body" validate:"required"`
}

// convertPostCommentRequestToCommentRequest converts the validated payload to an internal one
func convertPostCommentRequestToCommentRequest(req PostCommentRequest) comment.Comment {
	return comment.Comment{
		Slug:   req.Slug,
		Author: req.Author,
		Body:   req.Body,
	}
}

func (h *Handler) PostComment(w http.ResponseWriter, r *http.Request) {
	var postCmt PostCommentRequest
	if err := json.NewDecoder(r.Body).Decode(&postCmt); err != nil {
		return
	}

	validate := validator.New()
	err := validate.Struct(postCmt)

	if err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	cmt := convertPostCommentRequestToCommentRequest(postCmt)
	savedCmt, err := h.Service.PostComment(r.Context(), cmt)
	if err != nil {
		log.Error(errors.Wrap(err, "error posting comment"))
	}

	if err := json.NewEncoder(w).Encode(savedCmt); err != nil {
		panic(err)
	}

}

func (h *Handler) GetComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cmt, err := h.Service.GetComment(r.Context(), id)

	if err != nil {
		log.Error(errors.Wrap(err, "error getting comment"))
	}

	if err := json.NewEncoder(w).Encode(cmt); err != nil {
		panic(err)
	}
}

func (h *Handler) GetComments(w http.ResponseWriter, r *http.Request) {
	cmts, err := h.Service.GetComments(r.Context())
	if err != nil {
		log.Error(errors.Wrap(err, "error getting comments"))
	}
	if err := json.NewEncoder(w).Encode(cmts); err != nil {
		panic(err)
	}
}

func (h *Handler) UpdateComment(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	commentID := vars["id"]

	if commentID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	cmt, err := h.Service.GetComment(r.Context(), commentID)
	if err != nil {
		log.Error(errors.Wrap(err, "error finding the comment"))
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&cmt); err != nil {
		return
	}

	cmt, err = h.Service.UpdateComment(r.Context(), commentID, cmt)
	if err != nil {
		log.Error(errors.Wrap(err, "error updating comment"))
		return
	}

	if err := json.NewEncoder(w).Encode(cmt); err != nil {
		panic(err)
	}
}

func (h *Handler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	commentID := vars["id"]

	if commentID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Service.DeleteComment(r.Context(), commentID); err != nil {
		log.Error(errors.Wrap(err, "error deleting comment"))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(Response{
		Message: "comment deleted",
	}); err != nil {
		log.Error(errors.Wrap(err, "error encoding message"))
		panic(err)
	}

}
