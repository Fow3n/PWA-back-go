package repository

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"pwa/internal/models"
)

type UserRepository struct {
	Collection *mongo.Collection
}

func (r *UserRepository) CreateUser(ctx context.Context, user models.User) (*mongo.InsertOneResult, error) {
	result, err := r.Collection.InsertOne(ctx, user)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return result, nil
}

func (r *UserRepository) FindUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	cursor, err := r.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := cursor.Close(ctx); err != nil {
			fmt.Println(err)
		}
	}()

	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (r *UserRepository) FindUserByIdentifier(ctx context.Context, identifier string) (models.User, error) {
	var user models.User
	filter := bson.M{"$or": []bson.M{{"username": identifier}, {"email": identifier}}}
	err := r.Collection.FindOne(ctx, filter).Decode(&user)
	return user, err
}

func (r *UserRepository) UpdateUser(ctx context.Context, filter bson.M, update bson.M) (*mongo.UpdateResult, error) {
	return r.Collection.UpdateOne(ctx, filter, update)
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	return r.Collection.DeleteOne(ctx, filter)
}
