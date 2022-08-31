package server

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockPostsClient struct{}

func (c mockPostsClient) GetPostsWithTag(tag string) ([]byte, error) {
	switch tag {
	case "error":
		return nil, errors.New("test error")
	default:
		return []byte{}, nil
	}
}

func setup() *httptest.Server {
	return httptest.NewServer(router(mockPostsClient{}))
}

func TestHandlePing(t *testing.T) {
	ts := setup()
	defer ts.Close()

	req := httptest.NewRequest("GET", fmt.Sprintf("%s/ping", ts.URL), nil)
	w := httptest.NewRecorder()

	ts.Config.Handler.ServeHTTP(w, req)
	defer w.Result().Body.Close()

	actual, err := io.ReadAll(w.Result().Body)
	assert.Nil(t, err)

	expected := `{
	"success": true
}`
	assert.Equal(t, http.StatusOK, w.Code)
	assert.JSONEq(t, expected, string(actual))
}

func TestHandlePosts(t *testing.T) {
	ts := setup()
	defer ts.Close()

	testCases := []struct {
		endpoint     string
		expectedCode int
	}{
		{"posts", 400},
		{"posts?tags=science&sortBy=jibberish", 400},
		{"posts?tags=science&direction=sideways", 400},
		{"posts?tags=science", 200},
		{"posts?tags=science&sortBy=id&direction=desc", 200},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("endpoint: %s", tc.endpoint), func(t *testing.T) {
			req := httptest.NewRequest("GET", fmt.Sprintf("%s/%s", ts.URL, tc.endpoint), nil)
			w := httptest.NewRecorder()
			ts.Config.Handler.ServeHTTP(w, req)

			actual := w.Code

			assert.Equal(t, tc.expectedCode, actual)
		})
	}
}
