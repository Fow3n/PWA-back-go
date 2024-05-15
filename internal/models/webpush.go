package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WebPushSubscription struct {
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	UserID         primitive.ObjectID `bson:"userId" json:"userId"`
	Endpoint       string             `bson:"endpoint"`
	Keys           map[string]string  `bson:"keys"`
	ExpirationTime *int64             `bson:"expirationTime,omitempty" json:"expirationTime,omitempty"`
}
