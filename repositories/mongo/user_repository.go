package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	db *mongo.Database
}

type User struct {
	ID       string `bson:"_id,omitempty" json:"id"`
	Username string `bson:"username" json:"username"`
	Email    string `bson:"email" json:"email"`
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user User) error {
	collection := r.db.Collection("users")
	_, err := collection.InsertOne(ctx, user)
	return err
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (User, error) {
	var user User
	collection := r.db.Collection("users")
	err := collection.FindOne(ctx, map[string]string{"_id": id}).Decode(&user)
	return user, err
}
