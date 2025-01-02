package blog

import (
	"context"

	"github.com/Mikiejoe/go-blog-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	collection *mongo.Collection
}

func NewStore(collection *mongo.Collection) *Store {
	return &Store{collection: collection}
}

func (s *Store) CreateBlog(blog types.Blog) (string, error) {
	result, err := s.collection.InsertOne(context.Background(), blog)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s *Store) GetBlogs() ([]types.Blog, error) {
	filter := bson.D{}
	var blogs []types.Blog
	var blog types.Blog
	cur, err := s.collection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		if err := cur.Decode(&blog); err != nil {
			continue
		}
		blogs = append(blogs, blog)

	}

	return blogs, nil
}
func (s *Store) GetBlogByID(id string) (types.Blog, error) {
	var blog types.Blog
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return blog, err
	}
	filter := bson.D{{Key: "_id", Value: docId}}
	err = s.collection.FindOne(context.Background(), filter).Decode(&blog)
	if err == mongo.ErrNoDocuments {
		return types.Blog{}, mongo.ErrNoDocuments
	} else if err != nil {
		return types.Blog{}, err
	}

	return blog, nil
}
func (s *Store) UpdateBlog(string) error {
	return nil
}
func (s *Store) DeleteBlog(string) error {
	return nil
}
