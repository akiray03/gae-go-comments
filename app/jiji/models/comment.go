package models

import (
	"google.golang.org/appengine/datastore"
	"time"
)

const KIND_COMMENT = "Comment"

type Comment struct {
	Id     int64          `datastore:"-" goon:"id"           json:"id"`
	_kind  string         `datastore:"-" goon:"kind,Comment" json:"_kind"`
	Parent *datastore.Key `datastore:"-" goon:"parent"       json:"parent"`
	Author string         `datastore:"author"       json:"author"`
	Body   string         `datastore:"body,noindex" json:"body"`
	CreatedAt time.Time   `datastore:"created_at"   json:"created_at"`
	UpdatedAt time.Time   `datastore:"updated_at"   json:"updated_at"`
}

func (comment *Comment) Kind() string  {
	return comment._kind
}

func NewComment() Comment {
	comment := Comment{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	return comment
}
