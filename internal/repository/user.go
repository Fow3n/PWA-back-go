package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"pwa/internal/models"
)

type UserRepository struct {
	Collection *mongo.Collection
}

func (r *UserRepository) CreateUser(ctx context.Context, user models.User) (*mongo.InsertOneResult, error) {
	return r.Collection.InsertOne(ctx, user)
}

func (r *UserRepository) FindUserByID(ctx context.Context, id string) (models.User, error) {
	var user models.User
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	err := r.Collection.FindOne(ctx, filter).Decode(&user)
	return user, err
}

func (r *UserRepository) UpdateUser(ctx context.Context, id string, user models.User) (*mongo.UpdateResult, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": user}
	return r.Collection.UpdateOne(ctx, filter, update)
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	return r.Collection.DeleteOne(ctx, filter)
}
