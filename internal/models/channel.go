package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Channel struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string             `bson:"name" json:"name"`
	Members   []string           `bson:"members" json:"members"`
	Password  string             `bson:"password,omitempty"`
	CreatedAt time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt time.Time          `bson:"updatedAt" json:"updatedAt"`
}
