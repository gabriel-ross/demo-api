package server

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

// PostsClient is an interface for a client that can retrieve
// Post slices.
type PostsClient interface {
	GetPostsWithTag(string) ([]byte, error)
}

// HatchwayClientWithCache is a client for interacting with the hatchway
// api.
type hatchwayClientWithCache struct {
	url          string
	cache        map[string]cachedData
	cacheTimeout time.Duration
}

type cachedData struct {
	data    []byte
	timeout time.Time
}

// GetPostsWithTag gets posts with a tag. If the tag is cached and the
// data is not stale data is returned from the cache. Otherwise data
// is requested from the primary server and returned. The cache is updated
// with this fresh data.
func (h hatchwayClientWithCache) GetPostsWithTag(tag string) ([]byte, error) {

	reqURL := fmt.Sprintf("%s/posts?tag=%s", h.url, tag)

	cachedEntry, exists := h.cache[reqURL]
	if exists && cachedEntry.timeout.Before(time.Now()) {
		return cachedEntry.data, nil
	}

	res, err := http.Get(reqURL)
	if err != nil {
		return []byte{}, err
	}
	data, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	h.cache[reqURL] = cachedData{
		data:    data,
		timeout: time.Now().Add(h.cacheTimeout),
	}
	return data, nil
}
