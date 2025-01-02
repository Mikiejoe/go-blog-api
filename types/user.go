package types

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Username  string `bson:"username,omitempty" json:"username"`
	Email     string `bson:"email,omitempty" json:"email"`
	Firstname string `bson:"firstname,omitempty" json:"firstname"`
	Lastname  string `bson:"lastname,omitempty" json:"lastname"`
	Location  string `bson:"location,omitempty" json:"location"`
	Password  string `bson:"password,omitempty" json:"password"`
}

type RegisterUserPayload struct {
	Username  string `bson:"username,omitempty" json:"username"`
	Email     string `bson:"email,omitempty" json:"email"`
	Firstname string `bson:"firstname,omitempty" json:"firstname"`
	Lastname  string `bson:"lastname,omitempty" json:"lastname"`
	Location  string `bson:"location,omitempty" json:"location"`
	Password  string `bson:"password,omitempty" json:"password"`
}

type LoginUserPayload struct {
	Password string `bson:"password,omitempty" json:"password"`
	Username string `bson:"username,omitempty" json:"username"`
}

type UserInTerface interface {
	CreateUser(User) (string,error)
	GetUserByID(string) (User, error)
	GetUserByEmail(string) (User, error)
	GetUserByName(string) (User, error)
	GetUsers() ([]User, error)
}
