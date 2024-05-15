package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Task struct {
	ID          primitive.ObjectID `bson:"_id" json:"id"`
	Title       string             `bson:"title" json:"title"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Completed   bool               `bson:"completed" json:"completed"`
	CreatedAt   time.Time          `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt" json:"updatedAt"`
	UpdatedBy   primitive.ObjectID `bson:"updatedBy" json:"updatedBy"`
}

type TodoList struct {
	ID          primitive.ObjectID  `bson:"_id" json:"id"`
	Title       string              `bson:"title" json:"title"`
	Description string              `bson:"description,omitempty" json:"description,omitempty"`
	Owner       primitive.ObjectID  `bson:"owner" json:"owner"`
	ChannelID   *primitive.ObjectID `bson:"channelId,omitempty" json:"channelId,omitempty"`
	Tasks       []Task              `bson:"tasks" json:"tasks"`
	CreatedAt   time.Time           `bson:"createdAt" json:"createdAt"`
	UpdatedAt   time.Time           `bson:"updatedAt" json:"updatedAt"`
}
