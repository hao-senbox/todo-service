package shop

import (
	"context"
	"fmt"
	"time"
	"todo-service/internal/uploader"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShopService interface {
	CreateShop(ctx context.Context, req CreateShopRequest, ownerID string) (*string, error)
	GetMyShop(ctx context.Context, ownerID string) (*Shop, error)
	GetShopByID(ctx context.Context, shopID string) (*Shop, error)
	UpdateShop(ctx context.Context, shopID string, req UpdateShopRequest, ownerID string) error
	DeleteShop(ctx context.Context, shopID string, ownerID string) error

	CreateProduct(ctx context.Context, req CreateProductRequest, shopID string, ownerID string) (*string, error)
	GetProductsByShop(ctx context.Context, shopID string, ownerID string) ([]Product, error)
	UpdateProduct(ctx context.Context, productID string, req UpdateProductRequest, ownerID string) error
	DeleteProduct(ctx context.Context, productID string, ownerID string) error

	AddRepairItem(ctx context.Context, repairID primitive.ObjectID, req AddRepairItemRequest) (*RepairItem, error)
	GetRepairItems(ctx context.Context, repairID primitive.ObjectID) (*RepairItemsSummary, error)
	RemoveRepairItem(ctx context.Context, repairItemID primitive.ObjectID) error
}

type shopService struct {
	ShopRepo    ShopRepository
	UploaderSvc uploader.ImageService
}

func NewShopService(shopRepo ShopRepository, uploaderSvc uploader.ImageService) ShopService {
	return &shopService{
		ShopRepo:    shopRepo,
		UploaderSvc: uploaderSvc,
	}
}

func (s *shopService) CreateShop(ctx context.Context, req CreateShopRequest, ownerID string) (*string, error) {
	
	if req.Name == "" {
		return nil, fmt.Errorf("name is required")
	}

	shop, err := s.ShopRepo.GetMyShop(ctx, ownerID)
	if err != nil {
		return nil, err
	}

	if shop != nil {
		return nil, fmt.Errorf("shop already exists")
	}

	id := primitive.NewObjectID()

	shopNew := &Shop{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		OwnerID:     ownerID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.ShopRepo.CreateShop(ctx, shopNew)
	if err != nil {
		return nil, err
	}

	idParse := id.Hex()

	return &idParse, nil
}

func (s *shopService) GetMyShop(ctx context.Context, ownerID string) (*Shop, error) {
	shop, err := s.ShopRepo.GetMyShop(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	return shop, nil
}

func (s *shopService) GetShopByID(ctx context.Context, shopID string) (*Shop, error) {
	
	shopIDObj, err := primitive.ObjectIDFromHex(shopID)
	if err != nil {
		return nil, err
	}

	shop, err := s.ShopRepo.GetShopByID(ctx, shopIDObj)
	if err != nil {
		return nil, err
	}

	return shop, nil
}

func (s *shopService) UpdateShop(ctx context.Context, shopID string, req UpdateShopRequest, ownerID string) error {
	
	shopIDObj, err := primitive.ObjectIDFromHex(shopID)
	if err != nil {
		return err
	}

	shop, err := s.ShopRepo.GetShopByID(ctx, shopIDObj)
	if err != nil {
		return err
	}
	if shop == nil {
		return fmt.Errorf("shop not found")
	}
	if shop.OwnerID != ownerID {
		return fmt.Errorf("you are not the owner of this shop")
	}

	if req.Name != nil {
		shop.Name = *req.Name
	}
	if req.Description != nil {
		shop.Description = *req.Description
	}

	shop.UpdatedAt = time.Now()

	err = s.ShopRepo.UpdateShop(ctx, shopIDObj, shop)
	if err != nil {
		return err
	}

	return nil
}

func (s *shopService) DeleteShop(ctx context.Context, shopID string, ownerID string) error {
	shopIDObj, err := primitive.ObjectIDFromHex(shopID)
	if err != nil {
		return err
	}

	shop, err := s.ShopRepo.GetShopByID(ctx, shopIDObj)
	if err != nil {
		return err
	}
	if shop == nil {
		return fmt.Errorf("shop not found")
	}

	if shop.OwnerID != ownerID {
		return fmt.Errorf("you are not the owner of this shop")
	}

	return s.ShopRepo.DeleteShop(ctx, shopIDObj)
}

func (s *shopService) CreateProduct(ctx context.Context, req CreateProductRequest, shopID string, ownerID string) (*string, error) {
	
	shopIDObj, err := primitive.ObjectIDFromHex(shopID)
	if err != nil {
		return nil, err
	}

	shop, err := s.ShopRepo.GetShopByID(ctx, shopIDObj)
	if err != nil {
		return nil, err
	}

	if shop == nil {
		return nil, fmt.Errorf("shop not found")
	}

	if shop.OwnerID != ownerID {
		return nil, fmt.Errorf("you are not the owner of this shop")
	}

	id := primitive.NewObjectID()



	product := &Product{
		ID:          id,
		ShopID:      shopID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Category:    req.Category,
		ImageURL:    req.ImageURL,
		Stock:       req.Stock,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = s.ShopRepo.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	idParse := id.Hex()

	return &idParse, nil
}

func (s *shopService) GetProductsByShop(ctx context.Context, shopID string, ownerID string) ([]Product, error) {
	
	shopIDObj, err := primitive.ObjectIDFromHex(shopID)
	if err != nil {
		return nil, err
	}

	shop, err := s.ShopRepo.GetShopByID(ctx, shopIDObj)
	if err != nil {
		return nil, err
	}

	if shop == nil {
		return nil, fmt.Errorf("shop not found")
	}

	if shop.OwnerID != ownerID {
		return nil, fmt.Errorf("you are not the owner of this shop")
	}

	return s.ShopRepo.GetProductsByShop(ctx, shopID)

}

func (s *shopService) UpdateProduct(ctx context.Context, productID string, req UpdateProductRequest, ownerID string) error {
	
	productIDObj, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return err
	}

	existingProduct, err := s.ShopRepo.GetProductByID(ctx, productIDObj)
	if err != nil {
		return err
	}

	if req.Name != nil {
		existingProduct.Name = *req.Name
	}

	if req.Description != nil {
		existingProduct.Description = *req.Description
	}

	if req.Price != nil {
		existingProduct.Price = *req.Price
	}

	if req.Category != nil {
		existingProduct.Category = *req.Category
	}

	if req.ImageURL != nil {
		existingProduct.ImageURL = *req.ImageURL
	}

	if req.Stock != nil {
		existingProduct.Stock = *req.Stock
	}

	if req.IsActive != nil {
		existingProduct.IsActive = *req.IsActive
	}

	existingProduct.UpdatedAt = time.Now()

	return s.ShopRepo.UpdateProduct(ctx, productIDObj, existingProduct)

}

func (s *shopService) DeleteProduct(ctx context.Context, productID string, ownerID string) error {
	
	productIDObj, err := primitive.ObjectIDFromHex(productID)
	if err != nil {
		return err
	}

	product, err := s.ShopRepo.GetProductByID(ctx, productIDObj)
	if err != nil {
		return err
	}

	if product == nil {
		return fmt.Errorf("product not found")
	}

	return s.ShopRepo.DeleteProduct(ctx, productIDObj)
}

func (s *shopService) AddRepairItem(ctx context.Context, repairID primitive.ObjectID, req AddRepairItemRequest) (*RepairItem, error) {
	// Get product details to get current price
	product, err := s.ShopRepo.GetProductByID(ctx, req.ProductID)
	if err != nil {
		return nil, fmt.Errorf("product not found: %v", err)
	}

	if !product.IsActive {
		return nil, fmt.Errorf("product is not active")
	}

	if product.Stock < req.Quantity {
		return nil, fmt.Errorf("insufficient stock: available %d, requested %d", product.Stock, req.Quantity)
	}

	repairItem := &RepairItem{
		ID:        primitive.NewObjectID(),
		RepairID:  repairID,
		ProductID: req.ProductID,
		Quantity:  req.Quantity,
		Price:     product.Price, // Store price at time of adding
		AddedAt:   time.Now(),
	}

	err = s.ShopRepo.CreateRepairItem(ctx, repairItem)
	if err != nil {
		return nil, err
	}

	return repairItem, nil
}

func (s *shopService) GetRepairItems(ctx context.Context, repairID primitive.ObjectID) (*RepairItemsSummary, error) {
	repairItems, err := s.ShopRepo.GetRepairItems(ctx, repairID)
	if err != nil {
		return nil, err
	}

	var itemResponses []RepairItemResponse
	var totalCost float64

	for _, item := range repairItems {
		product, err := s.ShopRepo.GetProductByID(ctx, item.ProductID)
		if err != nil {
			continue // Skip if product not found
		}

		subTotal := item.Price * float64(item.Quantity)
		totalCost += subTotal

		itemResponse := RepairItemResponse{
			ID: item.ID,
			Product: ProductResponse{
				ID:          product.ID,
				ShopID:      product.ShopID,
				Name:        product.Name,
				Description: product.Description,
				Price:       product.Price,
				Category:    product.Category,
				ImageURL:    product.ImageURL,
				Stock:       product.Stock,
				IsActive:    product.IsActive,
				CreatedAt:   product.CreatedAt,
				UpdatedAt:   product.UpdatedAt,
			},
			Quantity: item.Quantity,
			Price:    item.Price,
			SubTotal: subTotal,
			AddedAt:  item.AddedAt,
		}
		itemResponses = append(itemResponses, itemResponse)
	}

	return &RepairItemsSummary{
		Items:     itemResponses,
		TotalCost: totalCost,
		ItemCount: len(itemResponses),
	}, nil
}

func (s *shopService) RemoveRepairItem(ctx context.Context, repairItemID primitive.ObjectID) error {
	return s.ShopRepo.DeleteRepairItem(ctx, repairItemID)
}
