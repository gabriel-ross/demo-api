package server

import (
	"sort"
)

// SortPostsUnique sorts a slice of posts in-place and removes
// duplicates.
func SortPostsUnique(inp *[]Post, sortBy, sortDirection string) {

	posts := *inp

	if len(posts) <= 1 {
		return
	}

	var less func(i, j int) bool
	switch sortBy {
	case "id":
		less = func(i, j int) bool { return posts[i].Id < posts[j].Id }
	case "reads":
		less = func(i, j int) bool { return posts[i].Reads < posts[j].Reads }
	case "likes":
		less = func(i, j int) bool { return posts[i].Likes < posts[j].Likes }
	case "popularity":
		less = func(i, j int) bool { return posts[i].Popularity < posts[j].Popularity }
	}

	if sortDirection == "desc" {
		sort.Slice(posts, func(i, j int) bool { return less(j, i) })
	} else {
		sort.Slice(posts, less)
	}

	i := 0
	for j := 1; j < len(posts); j++ {
		if posts[i].Id != posts[j].Id {
			i++
			posts[i] = posts[j]
		}
	}
	i++
	*inp = posts[:i]
}
