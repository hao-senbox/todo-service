package repair

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RepairRepository interface {
	CreateRepair(ctx context.Context, repair *Repair) error
	GetRepairs(ctx context.Context) ([]*Repair, error)
	GetJobCount(ctx context.Context, organizationID string) (int, error)
	GetRepairByID(ctx context.Context, id primitive.ObjectID) (*Repair, error)
	UpdateRepair(ctx context.Context, id primitive.ObjectID, repair *Repair) error
	DeleteRepair(ctx context.Context, id primitive.ObjectID) error
}

type repairRepository struct {
	repairCollection *mongo.Collection
}

func NewRepairRepository(repairCollection *mongo.Collection) RepairRepository {
	return &repairRepository{
		repairCollection: repairCollection,
	}
}

func (r *repairRepository) CreateRepair(ctx context.Context, repair *Repair) error {
	_, err := r.repairCollection.InsertOne(ctx, repair)
	if err != nil {
		return err
	}
	return nil
}

func (r *repairRepository) GetRepairs(ctx context.Context) ([]*Repair, error) {

	cursor, err := r.repairCollection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var repairs []*Repair

	for cursor.Next(ctx) {
		var repair Repair
		if err := cursor.Decode(&repair); err != nil {
			return nil, err
		}
		repairs = append(repairs, &repair)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return repairs, nil
}

func (r *repairRepository) GetJobCount(ctx context.Context, organizationID string) (int, error) {
	count, err := r.repairCollection.CountDocuments(ctx, bson.M{"organization_id": organizationID})
	if err != nil {
		return 0, err
	}
	return int(count), nil
}

func (r *repairRepository) GetRepairByID(ctx context.Context, id primitive.ObjectID) (*Repair, error) {
	var repair Repair
	err := r.repairCollection.FindOne(ctx, bson.M{"_id": id}).Decode(&repair)
	if err != nil {
		return nil, err
	}
	return &repair, nil
}

func (r *repairRepository) UpdateRepair(ctx context.Context, id primitive.ObjectID, repair *Repair) error {
	_, err := r.repairCollection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": repair})
	if err != nil {
		return err
	}
	return nil
}

func (r *repairRepository) DeleteRepair(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.repairCollection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}