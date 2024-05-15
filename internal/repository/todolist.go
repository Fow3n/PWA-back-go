package repository

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
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

func (r *TodoListRepository) FindTodoListsByUserID(ctx context.Context, userID string) ([]models.TodoList, error) {
	var todoLists []models.TodoList
	filter := bson.M{"owner": userID}
	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var todoList models.TodoList
		if err := cursor.Decode(&todoList); err != nil {
			return nil, err
		}
		todoLists = append(todoLists, todoList)
	}
	return todoLists, nil
}

func (r *TodoListRepository) FindTodoListsByChannelID(ctx context.Context, channelID string) ([]models.TodoList, error) {
	var todoLists []models.TodoList
	objID, _ := primitive.ObjectIDFromHex(channelID)
	filter := bson.M{"channelId": objID}
	cursor, err := r.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {
			fmt.Println(err)
		}
	}(cursor, ctx)

	for cursor.Next(ctx) {
		var todoList models.TodoList
		if err := cursor.Decode(&todoList); err != nil {
			return nil, err
		}
		todoLists = append(todoLists, todoList)
	}
	return todoLists, nil
}

func (r *TodoListRepository) AddTaskToList(ctx context.Context, todoListID string, task models.Task) error {
	tid, _ := primitive.ObjectIDFromHex(todoListID)
	filter := bson.M{"_id": tid}
	update := bson.M{"$push": bson.M{"tasks": task}}
	_, err := r.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *TodoListRepository) UpdateTask(ctx context.Context, todoListID string, taskID string, task models.Task) error {
	tid, _ := primitive.ObjectIDFromHex(todoListID)
	tkID, _ := primitive.ObjectIDFromHex(taskID)
	filter := bson.M{"_id": tid, "tasks._id": tkID}
	update := bson.M{"$set": bson.M{"tasks.$": task}}
	_, err := r.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *TodoListRepository) DeleteTask(ctx context.Context, todoListID string, taskID string) error {
	tid, _ := primitive.ObjectIDFromHex(todoListID)
	tkID, _ := primitive.ObjectIDFromHex(taskID)
	filter := bson.M{"_id": tid}
	update := bson.M{"$pull": bson.M{"tasks": bson.M{"_id": tkID}}}
	_, err := r.Collection.UpdateOne(ctx, filter, update)
	return err
}

func (r *TodoListRepository) GetTaskByID(c *gin.Context, todoListID string, taskID string) (models.Task, error) {
	var task models.Task
	tid, _ := primitive.ObjectIDFromHex(todoListID)
	tkID, _ := primitive.ObjectIDFromHex(taskID)
	filter := bson.M{"_id": tid, "tasks._id": tkID}
	err := r.Collection.FindOne(c, filter).Decode(&task)
	return task, err
}
