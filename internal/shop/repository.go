package shop

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShopRepository interface {
	CreateShop(ctx context.Context, shop *Shop) error
	GetMyShop(ctx context.Context, ownerID string) (*Shop, error)
	GetShopByID(ctx context.Context, shopID primitive.ObjectID) (*Shop, error)
	UpdateShop(ctx context.Context, shopID primitive.ObjectID, shop *Shop) error
	DeleteShop(ctx context.Context, shopID primitive.ObjectID) error

	CreateProduct(ctx context.Context, product *Product) error
	GetProductsByShop(ctx context.Context, shopID string) ([]Product, error)
	GetProductByID(ctx context.Context, productID primitive.ObjectID) (*Product, error)
	UpdateProduct(ctx context.Context, productID primitive.ObjectID, product *Product) error
	DeleteProduct(ctx context.Context, productID primitive.ObjectID) error

	CreateRepairItem(ctx context.Context, repairItem *RepairItem) error
	GetRepairItems(ctx context.Context, repairID primitive.ObjectID) ([]RepairItem, error)
	GetRepairItemTotal(ctx context.Context, repairID primitive.ObjectID) (float64, error)
	DeleteRepairItem(ctx context.Context, repairItemID primitive.ObjectID) error
}

type shopRepository struct {
	productCollection    *mongo.Collection
	repairItemCollection *mongo.Collection
	shopCollection       *mongo.Collection
}

func NewShopRepository(productCollection, repairItemCollection, shopCollection *mongo.Collection) ShopRepository {
	return &shopRepository{
		productCollection:    productCollection,
		repairItemCollection: repairItemCollection,
		shopCollection:       shopCollection,
	}
}

func (r *shopRepository) CreateShop(ctx context.Context, shop *Shop) error {
	_, err := r.shopCollection.InsertOne(ctx, shop)
	return err
}

func (r *shopRepository) GetMyShop(ctx context.Context, ownerID string) (*Shop, error) {
	var shop Shop
	err := r.shopCollection.FindOne(ctx, bson.M{"owner_id": ownerID}).Decode(&shop)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &shop, nil
}

func (r *shopRepository) GetShopByID(ctx context.Context, shopID primitive.ObjectID) (*Shop, error) {
	var shop Shop
	err := r.shopCollection.FindOne(ctx, bson.M{"_id": shopID}).Decode(&shop)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &shop, nil
}

func (r *shopRepository) UpdateShop(ctx context.Context, shopID primitive.ObjectID, shop *Shop) error {
	_, err := r.shopCollection.UpdateOne(ctx, bson.M{"_id": shopID}, bson.M{"$set": shop})
	return err
}

func (r *shopRepository) DeleteShop(ctx context.Context, shopID primitive.ObjectID) error {
	_, err := r.shopCollection.DeleteOne(ctx, bson.M{"_id": shopID})
	return err
}

func (r *shopRepository) CreateProduct(ctx context.Context, product *Product) error {
	_, err := r.productCollection.InsertOne(ctx, product)
	return err
}

func (r *shopRepository) GetProductsByShop(ctx context.Context, shopID string) ([]Product, error) {
	cursor, err := r.productCollection.Find(ctx, bson.M{"shop_id": shopID, "is_active": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []Product
	for cursor.Next(ctx) {
		var product Product
		if err := cursor.Decode(&product); err != nil {
			return nil, err
		}
		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *shopRepository) GetProductByID(ctx context.Context, productID primitive.ObjectID) (*Product, error) {
	var product Product
	err := r.productCollection.FindOne(ctx, bson.M{"_id": productID}).Decode(&product)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *shopRepository) UpdateProduct(ctx context.Context, productID primitive.ObjectID, product *Product) error {
	_, err := r.productCollection.UpdateOne(ctx, bson.M{"_id": productID}, bson.M{"$set": product})
	return err
}

func (r *shopRepository) DeleteProduct(ctx context.Context, productID primitive.ObjectID) error {
	_, err := r.productCollection.UpdateOne(ctx, bson.M{"_id": productID}, bson.M{"$set": bson.M{"is_active": false}})
	return err
}

func (r *shopRepository) CreateRepairItem(ctx context.Context, repairItem *RepairItem) error {
	_, err := r.repairItemCollection.InsertOne(ctx, repairItem)
	return err
}

func (r *shopRepository) GetRepairItems(ctx context.Context, repairID primitive.ObjectID) ([]RepairItem, error) {
	cursor, err := r.repairItemCollection.Find(ctx, bson.M{"repair_id": repairID})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var repairItems []RepairItem
	for cursor.Next(ctx) {
		var repairItem RepairItem
		if err := cursor.Decode(&repairItem); err != nil {
			return nil, err
		}
		repairItems = append(repairItems, repairItem)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return repairItems, nil
}

func (r *shopRepository) GetRepairItemTotal(ctx context.Context, repairID primitive.ObjectID) (float64, error) {
	repairItems, err := r.GetRepairItems(ctx, repairID)
	if err != nil {
		return 0, err
	}

	var total float64
	for _, item := range repairItems {
		total += item.Price * float64(item.Quantity)
	}

	return total, nil
}

func (r *shopRepository) DeleteRepairItem(ctx context.Context, repairItemID primitive.ObjectID) error {
	_, err := r.repairItemCollection.DeleteOne(ctx, bson.M{"_id": repairItemID})
	return err
}
