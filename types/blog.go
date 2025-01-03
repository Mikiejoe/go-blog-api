package types

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Blog struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId    primitive.ObjectID `bson:"userid,omitempty" json:"author"`
	Title     string             `bson:"title,omitempty" json:"title"`
	Content   string             `bson:"content,omitempty" json:"content"`
	CreatedAt time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt,omitempty" json:"updateddAt"`
}

type BlogInterface interface {
	CreateBlog(Blog) (string, error)
	GetBlogs() ([]Blog, error)
	GetBlogByID(string) (Blog, error)
	UpdateBlog(string, Blog) error
	DeleteBlog(string) error
}

type BlogPayload struct {
	Title   string `bson:"title,omitempty" json:"title"`
	Content string `bson:"content,omitempty" json:"content"`
}
