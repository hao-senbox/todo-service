package todo

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TodoRepository interface {
	GetAllTodo(ctx context.Context, status, name, teacher, student, staff string) ([]*Todo, error)
	GetTodoByID(ctx context.Context, todoID primitive.ObjectID) (*Todo, error)
	CreateTodo(ctx context.Context, todo *Todo) (*string, error)
	UpdateTodo(ctx context.Context, todo *Todo) error
	DeleteTodo(ctx context.Context, todoID primitive.ObjectID) error
	// Join Todo
	GetTodoByQRCode(ctx context.Context, qrCode string) (*Todo, error)
	JoinTodo(ctx context.Context, todoID primitive.ObjectID, userID, typeUser string) error
	AddUsers(ctx context.Context, todoID primitive.ObjectID, userIDs []string, typeUser string) error
	GetMyTodo(ctx context.Context, userID string) ([]*Todo, error)
}

type todoRepository struct {
	todoCollection *mongo.Collection
}

func NewTodoRepository(todoCollection *mongo.Collection) TodoRepository {
	return &todoRepository{
		todoCollection: todoCollection,
	}
}

func (r *todoRepository) GetAllTodo(ctx context.Context, status, name, teacher, student, staff string) ([]*Todo, error) {

	var todos []*Todo

	filter := bson.M{}

	if status != "" {
		filter["status"] = status
	}

	if name != "" {
		filter["name"] = name
	}

	if teacher != "" {
		filter["task_users.teachers"] = teacher
	}

	if student != "" {
		filter["task_users.students"] = student
	}

	if staff != "" {
		filter["task_users.staff"] = staff
	}

	cursor, err := r.todoCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(ctx) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	return todos, nil

}

func (r *todoRepository) GetTodoByID(ctx context.Context, todoID primitive.ObjectID) (*Todo, error) {

	var todo Todo

	err := r.todoCollection.FindOne(ctx, bson.M{"_id": todoID}).Decode(&todo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &todo, nil

}

func (r *todoRepository) CreateTodo(ctx context.Context, todo *Todo) (*string, error) {

	result, err := r.todoCollection.InsertOne(ctx, todo)
	if err != nil {
		return nil, err
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, fmt.Errorf("failed to convert inserted ID to ObjectID")
	}

	id := objectID.Hex()

	return &id, nil
}

func (r *todoRepository) UpdateTodo(ctx context.Context, todo *Todo) error {
	_, err := r.todoCollection.UpdateOne(ctx, bson.M{"_id": todo.ID}, bson.M{"$set": todo})
	return err
}

func (r *todoRepository) DeleteTodo(ctx context.Context, todoID primitive.ObjectID) error {
	_, err := r.todoCollection.DeleteOne(ctx, bson.M{"_id": todoID})
	return err
}

func (r *todoRepository) GetTodoByQRCode(ctx context.Context, qrCode string) (*Todo, error) {

	var todo Todo

	err := r.todoCollection.FindOne(ctx, bson.M{"qrcode": qrCode}).Decode(&todo)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &todo, nil

}

func (r *todoRepository) JoinTodo(ctx context.Context, todoID primitive.ObjectID, userID, typeUser string) error {

	filter := bson.M{
		"_id": todoID,
	}

	var filed string
	switch typeUser {
	case "student":
		filed = "task_users.students"
	case "teacher":
		filed = "task_users.teachers"
	case "staff":
		filed = "task_users.staff"
	default:
		return fmt.Errorf("type user not found")
	}

	update := bson.M{
		"$addToSet": bson.M{
			filed: userID,
		},
	}

	_, err := r.todoCollection.UpdateOne(ctx, filter, update)
	return err
}

func (r *todoRepository) AddUsers(ctx context.Context, todoID primitive.ObjectID, userIDs []string, typeUser string) error {

	filter := bson.M{
		"_id": todoID,
	}

	var field string
	switch typeUser {
	case "student":
		field = "task_users.students"
	case "teacher":
		field = "task_users.teachers"
	case "staff":
		field = "task_users.staff"
	default:
		return fmt.Errorf("type user not found")
	}

	update := bson.M{
		"$set": bson.M{
			field: userIDs, 
		},
	}

	_, err := r.todoCollection.UpdateOne(ctx, filter, update)
	return err
}

func (r *todoRepository) GetMyTodo(ctx context.Context, userID string) ([]*Todo, error) {

	var todos []*Todo

	filter := bson.M{
		"$or": []bson.M{
			{"created_by": userID},
			{"task_users.teachers": bson.M{"$in": []string{userID}}},
			{"task_users.students": bson.M{"$in": []string{userID}}},
			{"task_users.staff": bson.M{"$in": []string{userID}}},
		},
	}

	cursor, err := r.todoCollection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return nil, err
		}
		todos = append(todos, &todo)
	}

	return todos, nil

}
