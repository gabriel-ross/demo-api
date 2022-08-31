package server

// Ping is a struct defining the response structure of a server ping.
type Ping struct {
	Success bool `json:"success"`
}

// Error is a struct defining the response structure of an error.
type Error struct {
	Error string `json:"error"`
}

// APIResponse is a struct matching the expected structure
// of a response from the hatchways API.
type APIResponse struct {
	Posts []Post `json:"posts"`
}

// Post is a struct matching the expected structure of
// a post json object from the hatchways API.
type Post struct {
	Id         int      `json:"id"`
	Author     string   `json:"author"`
	AuthorId   int      `json:"authorid"`
	Likes      int      `json:"likes"`
	Popularity float64  `json:"popularity"`
	Reads      int      `json:"reads"`
	Tags       []string `json:"tags"`
}
