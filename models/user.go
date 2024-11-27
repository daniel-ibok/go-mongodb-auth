package models

import (
	"context"
	"go-mongodb-auth/database"
	"go-mongodb-auth/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	UUID      string             `bson:"uuid" json:"uuid"`
	Email     string             `bson:"email" json:"email"`
	Username  string             `bson:"username" json:"username"`
	Password  string             `bson:"password" json:"-"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

type RegisterUser struct {
	Email    string `json:"email" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginUser struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func NewUser(user RegisterUser) error {
	collection := database.GetDBCollection()
	_, err := collection.InsertOne(context.TODO(), &User{
		ID:        primitive.NewObjectID(),
		UUID:      utils.GetUUID(),
		Username:  user.Username,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	return err
}

func GetUserByEmail(email string) (*User, error) {
	collection := database.GetDBCollection()

	var user *User
	filter := bson.D{primitive.E{Key: "email", Value: email}}
	err := collection.FindOne(context.TODO(), filter).Decode(&user)
	if err != nil {
		return &User{}, err
	}
	return user, nil
}

func CheckUser(email string) bool {
	user, _ := GetUserByEmail(email)
	return !user.ID.IsZero()
}
