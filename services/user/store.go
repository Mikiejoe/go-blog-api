package user

import (
	"context"
	"fmt"

	"github.com/Mikiejoe/go-blog-api/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	collection *mongo.Collection
}

func NewStore(collection *mongo.Collection) *Store {
	return &Store{
		collection: collection,
	}
}

func (s Store) CreateUser(user types.User) (string, error) {
	result, err := s.collection.InsertOne(context.Background(), user)
	if err != nil {
		return "", err
	}
	return result.InsertedID.(primitive.ObjectID).Hex(), nil
}

func (s Store) GetUserByEmail(email string) (types.User, error) {

	var user types.User
	filter := bson.D{{Key: "email", Value: email}}
	err := s.collection.FindOne(context.Background(), filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return types.User{}, mongo.ErrNoDocuments
	} else if err != nil {
		return types.User{}, err
	}

	return user, nil
}

func (s Store) GetUsers() ([]types.User, error) {

	var users []types.User
	var user types.User
	filter := bson.D{{}}
	cur, err := s.collection.Find(context.Background(),filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()){
		if err:= cur.Decode(&user);err!=nil{
			continue
		}
		users = append(users, user)
	}
	return users,nil
}
func (s Store) GetUserByName(username string) (types.User, error) {

	var user types.User
	filter := bson.D{{Key: "username", Value: username}}
	err := s.collection.FindOne(context.Background(), filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return types.User{}, mongo.ErrNoDocuments
	} else if err != nil {
		return types.User{}, err
	}
	fmt.Printf("user %#v", user)

	return user, nil
}

func (s Store) GetUserByID(id string) (types.User, error) {
	docId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return types.User{}, err
	}
	var user types.User
	filter := bson.D{{Key: "_id", Value: docId}}
	err = s.collection.FindOne(context.Background(), filter).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return types.User{}, fmt.Errorf("user not found")
	}

	return user, nil
}
