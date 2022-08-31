package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"

	"github.com/go-chi/chi"
)

func router(client PostsClient) chi.Router {
	r := chi.NewRouter()

	r.Get("/ping", handleGetPing)
	r.Get("/posts", handleGetPosts(client, []queryParam{tagQuery, sortByQuery, sortDirectionQuery}))

	return r
}

func handleGetPing(w http.ResponseWriter, r *http.Request) {
	asJson, err := json.MarshalIndent(Ping{true}, "", "	")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(asJson)
}

func handleGetPosts(client PostsClient, expectedQueries []queryParam) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		queries, err := extractAndValidateQueries(expectedQueries, r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}

		// make http requests for each tag
		tags := strings.Split(queries["tags"], ",")
		posts := []Post{}
		var wg sync.WaitGroup
		datac := make(chan []byte)
		requestsDone := make(chan bool)
		processingDone := make(chan bool)

		// goroutine responsible for unmarshaling api responses and adding them
		// to data slice
		go func() {
			for {
				select {
				case data := <-datac:
					var responseData APIResponse
					err = json.Unmarshal(data, &responseData)
					if err != nil {
						log.Println("Error unmarshaling response from primary server")
					}

					// append posts from api response
					posts = append(posts, responseData.Posts...)
				case <-requestsDone:
					close(processingDone)
					return
				}
			}
		}()

		// send http request for each tag concurrently
		for _, tag := range tags {
			wg.Add(1)

			go func(tag string) {
				data, err := client.GetPostsWithTag(tag)
				if err != nil {
					log.Println("Error retrieving data from primary server")
				}
				datac <- data
				wg.Done()
			}(tag)
		}

		// wait on requesting goroutine and processing goroutine to finish
		wg.Wait()
		close(requestsDone)
		<-processingDone

		SortPostsUnique(&posts, queries["sortBy"], queries["direction"])
		asJson, err := json.MarshalIndent(APIResponse{posts}, "", "	")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write(asJson)
	}
}

// queryParam is a struct for defining query parameters permitted
// on an endpoint.
type queryParam struct {
	Key             string
	Required        bool
	PermittedValues []string // empty means anything permitted
	DefaultValue    string
}

func (q *queryParam) isPermittedValue(val string) bool {
	if len(q.PermittedValues) == 0 {
		return true
	}
	for _, acceptableVal := range q.PermittedValues {
		if val == acceptableVal {
			return true
		}
	}
	return false
}

func (q *queryParam) extractAndValidate(r *http.Request) (string, error) {
	val := r.URL.Query().Get(q.Key)
	if val == "" {
		if q.Required {
			return "", fmt.Errorf("required parameter not provided: %s", q.Key)
		} else {
			return q.DefaultValue, nil
		}
	} else if q.isPermittedValue(val) {
		return val, nil
	} else {
		return "", fmt.Errorf("invalid value for parameter %s: %s", q.Key, val)
	}
}

var tagQuery = queryParam{
	Key:             "tags",
	Required:        true,
	PermittedValues: []string{},
	DefaultValue:    "",
}

var sortByQuery = queryParam{
	Key:             "sortBy",
	Required:        false,
	PermittedValues: []string{"id", "reads", "likes", "popularity"},
	DefaultValue:    "id",
}

var sortDirectionQuery = queryParam{
	Key:             "direction",
	Required:        false,
	PermittedValues: []string{"asc", "desc"},
	DefaultValue:    "asc",
}

// extractAndValidateQueries takes in a slice of expected queries as well
// as an http request and extracts and validates the parameters from the
// request. Any expected, but not required parameters not found in the
// request will be set to their default value.
func extractAndValidateQueries(expectedQueries []queryParam, r *http.Request) (map[string]string, error) {
	currentQueries := map[string]string{}

	for _, param := range expectedQueries {
		val, err := param.extractAndValidate(r)
		if err != nil {
			return map[string]string{}, err
		}
		currentQueries[param.Key] = val
	}

	return currentQueries, nil
}

func writeError(w http.ResponseWriter, statusCode int, errMsg string) {
	asJson, err := json.MarshalIndent(Error{errMsg}, "", "	")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(statusCode)
	w.Write(asJson)
}
