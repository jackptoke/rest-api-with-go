//go:build e2e

package tests

import (
	"fmt"
	resty "github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHealthCheckEndpoint(t *testing.T) {
	client := resty.New()
	resp, err := client.R().Get("http://localhost:8080/alive")
	assert.NoError(t, err)
	fmt.Println("Response Status:", resp.Status())
	assert.Equal(t, 200, resp.StatusCode())
}
