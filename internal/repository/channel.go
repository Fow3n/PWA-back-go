package repository

import (
	"context"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"pwa/internal/models"
)

type ChannelRepository struct {
	Collection *mongo.Collection
}

func (r *ChannelRepository) CreateChannel(ctx context.Context, channel models.Channel) (*mongo.InsertOneResult, error) {
	return r.Collection.InsertOne(ctx, channel)
}

func (r *ChannelRepository) FindChannelsByUserID(ctx context.Context, userID string) ([]models.Channel, error) {
	var channels []models.Channel
	filter := bson.M{"members": userID}
	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			// log this error
		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var channel models.Channel
		if err := cursor.Decode(&channel); err != nil {
			return nil, err
		}
		channels = append(channels, channel)
	}
	return channels, nil
}

func (r *ChannelRepository) FindChannelByID(ctx context.Context, id string) (models.Channel, error) {
	var channel models.Channel
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return channel, fmt.Errorf("invalid id format: %w", err)
	}

	filter := bson.M{"_id": objID}
	err = r.Collection.FindOne(ctx, filter).Decode(&channel)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return channel, fmt.Errorf("no channel found with ID: %s", id)
		}
		return channel, err
	}
	return channel, nil
}

func (r *ChannelRepository) UpdateChannel(ctx context.Context, id string, channel models.Channel) (*mongo.UpdateResult, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": channel}
	return r.Collection.UpdateOne(ctx, filter, update)
}

func (r *ChannelRepository) DeleteChannel(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	return r.Collection.DeleteOne(ctx, filter)
}

func (r *ChannelRepository) GetChannelMembers(ctx context.Context, channelID string) ([]string, error) {
	var channel models.Channel
	objID, err := primitive.ObjectIDFromHex(channelID)
	if err != nil {
		return nil, err
	}
	err = r.Collection.FindOne(ctx, bson.M{"_id": objID}).Decode(&channel)
	if err != nil {
		return nil, err
	}
	return channel.Members, nil
}
