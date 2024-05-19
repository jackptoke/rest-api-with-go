//go:build e2e

package tests

import (
	"fmt"
	resty "github.com/go-resty/resty/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func createToken() string {
	token := jwt.New(jwt.SigningMethodHS256)
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println(err)
	}
	return tokenString
}

func TestPostComment(t *testing.T) {
	t.Run("can post comment", func(t *testing.T) {
		client := resty.New()
		tokenString := createToken()
		resp, err := client.R().
			SetHeader("Authorization", "bearer "+tokenString).
			SetBody(`{"Slug": "/", "author": "Jack", "body": "This is an example of a comment"}`).
			Post("http://localhost:8080/api/v1/comment")
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
	})

	t.Run("can post comment", func(t *testing.T) {
		client := resty.New()
		resp, err := client.R().
			SetBody(`{"Slug": "/", "author": "Jack", "body": "This is an example of a comment"}`).
			Post("http://localhost:8080/api/v1/comment")
		assert.Error(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode())
	})
}
