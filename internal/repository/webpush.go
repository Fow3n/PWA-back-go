package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"pwa/internal/models"
)

type WebPushRepository struct {
	Collection *mongo.Collection
}

func (r *WebPushRepository) CreateWebPushSubscription(ctx context.Context, webPushSubscription models.WebPushSubscription) (*mongo.InsertOneResult, error) {
	return r.Collection.InsertOne(ctx, webPushSubscription)
}

func (r *WebPushRepository) FindWebPushSubscriptionByID(ctx context.Context, id string) (models.WebPushSubscription, error) {
	var webPushSubscription models.WebPushSubscription
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return webPushSubscription, err
	}
	filter := bson.M{"_id": objID}
	err = r.Collection.FindOne(ctx, filter).Decode(&webPushSubscription)
	return webPushSubscription, err
}

func (r *WebPushRepository) UpdateWebPushSubscription(ctx context.Context, id string, webPushSubscription models.WebPushSubscription) (*mongo.UpdateResult, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": webPushSubscription}
	return r.Collection.UpdateOne(ctx, filter, update)
}

func (r *WebPushRepository) DeleteWebPushSubscription(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := bson.M{"_id": objID}
	return r.Collection.DeleteOne(ctx, filter)
}

func (r *WebPushRepository) FindByUserID(ctx context.Context, userID primitive.ObjectID) ([]models.WebPushSubscription, error) {
	var subscriptions []models.WebPushSubscription
	filter := bson.M{"userID": userID}
	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	if err = cursor.All(ctx, &subscriptions); err != nil {
		return nil, err
	}
	return subscriptions, nil
}
