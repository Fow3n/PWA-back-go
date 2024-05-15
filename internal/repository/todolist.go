package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"pwa/internal/models"
)

type TodoListRepository struct {
	Collection *mongo.Collection
}

func (r *TodoListRepository) CreateTodoList(ctx context.Context, todoList models.TodoList) (*mongo.InsertOneResult, error) {
	return r.Collection.InsertOne(ctx, todoList)
}

func (r *TodoListRepository) FindTodoListByID(ctx context.Context, id string) (models.TodoList, error) {
	var todoList models.TodoList
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	err := r.Collection.FindOne(ctx, filter).Decode(&todoList)
	return todoList, err
}

func (r *TodoListRepository) UpdateTodoList(ctx context.Context, id string, todoList models.TodoList) (*mongo.UpdateResult, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	update := bson.M{"$set": todoList}
	return r.Collection.UpdateOne(ctx, filter, update)
}

func (r *TodoListRepository) DeleteTodoList(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	objID, _ := primitive.ObjectIDFromHex(id)
	filter := bson.M{"_id": objID}
	return r.Collection.DeleteOne(ctx, filter)
}
