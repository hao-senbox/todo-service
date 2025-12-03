package task

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TaskRepository interface {
	CreateTask(ctx context.Context, task *Task) error
	GetTasks(ctx context.Context, role string, status string) ([]*Task, error)
	GetTaskById(ctx context.Context, id primitive.ObjectID) (*Task, error)
	UpdateTask(ctx context.Context, id primitive.ObjectID, task *Task) error
	DeleteTask(ctx context.Context, id primitive.ObjectID) error
	GetMyTask(ctx context.Context, userID string) ([]*Task, error)
}

type taskRepository struct {
	taskCollection *mongo.Collection
}

func NewTaskRepository(taskCollection *mongo.Collection) TaskRepository {
	return &taskRepository{
		taskCollection: taskCollection,
	}
}

func (r *taskRepository) CreateTask(ctx context.Context, task *Task) error {
	_, err := r.taskCollection.InsertOne(ctx, task)
	if err != nil {
		return err
	}
	return nil
}

func (r *taskRepository) GetTasks(ctx context.Context, role string, status string) ([]*Task, error) {
	filter := bson.M{}

	if role != "" {
		filter["group.role"] = role
	}

	if status != "" {
		filter["group.status"] = status
	}

	opts := options.Find().SetSort(bson.M{"created_at": -1})

	cursor, err := r.taskCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*Task

	for cursor.Next(ctx) {
		var task Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	return tasks, nil
}

func (r *taskRepository) GetTaskById(ctx context.Context, id primitive.ObjectID) (*Task, error) {
	var task Task
	err := r.taskCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&task)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (r *taskRepository) UpdateTask(ctx context.Context, id primitive.ObjectID, task *Task) error {
	_, err := r.taskCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": task})
	if err != nil {
		return err
	}
	return nil
}

func (r *taskRepository) DeleteTask(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.taskCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (r *taskRepository) GetMyTask(ctx context.Context, userID string) ([]*Task, error) {

	filter := bson.M{
		"$or": []bson.M{
			{"leader.user_id": userID},
		},
	}

	opts := options.Find().SetSort(bson.M{"created_at": -1})

	cursor, err := r.taskCollection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*Task
	for cursor.Next(ctx) {
		var task Task
		if err := cursor.Decode(&task); err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}
	return tasks, nil
}
