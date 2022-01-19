package blog

import (
	"time"
)

type Post struct {
	Date   time.Time
	Author string
	Title  string
	Thumb  string
	Body   string
}

type Posts []*Post

var posts Posts
