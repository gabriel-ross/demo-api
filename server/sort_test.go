package server

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestSortPosts(t *testing.T) {

	inp := genScrambled()
	testCases := []struct {
		sortBy   string
		expected []Post
	}{
		{"id", idSortedGolden},
		{"reads", readsSortedGolden},
		{"likes", likesSortedGolden},
		{"popularity", popularitySortedGolden},
	}
	testCaseDirection := []string{"asc", "desc"}

	for _, dir := range testCaseDirection {
		for _, tc := range testCases {
			t.Run(fmt.Sprintf("sortBy: %s direction: %s", tc.sortBy, dir), func(t *testing.T) {
				actual := []Post{}
				actual = append(actual, inp...)
				SortPostsUnique(&actual, tc.sortBy, dir)
				if dir == "desc" {
					reverse(actual)
				}
				assert.Equal(t, tc.expected, actual)
			})
		}
	}
}

func TestSortTrivial(t *testing.T) {
	inp := []Post{}
	SortPostsUnique(&inp, "id", "asc")
	inp = []Post{post1}
	SortPostsUnique(&inp, "id", "asc")

	assert.Equal(t, 1, 1)
}

var idSortedGolden = []Post{post1, post2, post3, post4, post5}
var readsSortedGolden = []Post{post3, post2, post5, post4, post1}
var likesSortedGolden = []Post{post1, post4, post2, post3, post5}
var popularitySortedGolden = []Post{post4, post2, post5, post3, post1}

func reverse[T any](input []T) {
	for i, j := 0, len(input)-1; i < j; i, j = i+1, j-1 {
		input[i], input[j] = input[j], input[i]
	}
}

func genScrambled() []Post {
	scrambled := []Post{post1, post2, post3, post4, post5, post1, post2, post3, post4, post5}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(scrambled), func(i, j int) { scrambled[i], scrambled[j] = scrambled[j], scrambled[i] })
	return scrambled
}

var post1 = Post{
	Id:         1,
	Author:     "Spiderman",
	AuthorId:   1,
	Likes:      1,
	Popularity: .9,
	Reads:      1000,
	Tags:       []string{"marvel"},
}

var post2 = Post{
	Id:         2,
	Author:     "Batman",
	AuthorId:   500,
	Likes:      71,
	Popularity: .3,
	Reads:      232,
	Tags:       []string{"crime"},
}

var post3 = Post{
	Id:         3,
	Author:     "Superman",
	AuthorId:   9,
	Likes:      542,
	Popularity: .6,
	Reads:      197,
	Tags:       []string{"space"},
}

var post4 = Post{
	Id:         4,
	Author:     "Kratos",
	AuthorId:   1,
	Likes:      5,
	Popularity: .2,
	Reads:      922,
	Tags:       []string{"mythology"},
}

var post5 = Post{
	Id:         5,
	Author:     "Ivy",
	AuthorId:   1,
	Likes:      654,
	Popularity: .35,
	Reads:      876,
	Tags:       []string{"nature"},
}
