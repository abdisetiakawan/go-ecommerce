package usecase

import (
	"context"

	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/helper"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/model/converter"
	repo "github.com/abdisetiakawan/go-ecommerce/internal/repository/interfaces"
	"github.com/abdisetiakawan/go-ecommerce/internal/usecase/interfaces"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type ProductUseCase struct {
	db          *gorm.DB
	val         *validator.Validate
	productRepo repo.ProductRepository
	storeRepo   repo.StoreRepository
	uuid        *helper.UUIDHelper
}

func NewProductUseCase(db *gorm.DB, validate *validator.Validate, productRepo repo.ProductRepository, storeRepo repo.StoreRepository, uuid *helper.UUIDHelper) interfaces.ProductUseCase {
	return &ProductUseCase{
		db:          db,
		val:         validate,
		productRepo: productRepo,
		storeRepo: storeRepo,
		uuid:        uuid,
	}
}


// CreateProduct creates a new product and returns the newly created product in the response.
// It first validates the request body and checks if the user is a seller.
// If the user is not a seller, it returns a 403 error.
// If the request body is invalid, it returns a 400 error.
// If the product cannot be created, it returns a 500 error.
func (u *ProductUseCase) CreateProduct(ctx context.Context, request *model.RegisterProduct) (*model.ProductResponse, error) {
	if err := helper.ValidateStruct(u.val, request); err != nil {
		return nil, err
	}
	storeID, err := u.storeRepo.GetStoreIDByUserID(request.AuthID)
	if err != nil {
		return nil, err
	}
	product := &entity.Product{
		ProductUUID: u.uuid.Generate(),
		StoreID: storeID,
		ProductName: request.ProductName,
		Description: request.Description,
		Price: request.Price,
		Stock: request.Stock,
		Category: request.Category,
	}
	if err := u.productRepo.CreateProduct(product); err != nil {
		return nil, err
	}

	return converter.ProductToResponse(product), nil
}


// GetProducts retrieves a list of products based on the provided request filters.
// It validates the request, queries the product repository, and converts the results
// into a list of product responses.
//
// Parameters:
//   - ctx: Context for the request, allowing for cancellation and timeouts.
//   - request: Pointer to GetProductsRequest containing filters like search terms,
//     category, price range, pagination, etc.
//
// Returns:
//   - A slice of ProductResponse containing the product details.
//   - An int64 representing the total number of products that match the query.
//   - An error, if any occurs during validation or data retrieval.

func (u *ProductUseCase) GetProducts(ctx context.Context, request *model.GetProductsRequest) ([]model.ProductResponse, int64, error) {
	if err := helper.ValidateStruct(u.val, request); err != nil {
		return nil, 0, err
	}
	products, total, err := u.productRepo.GetProducts(request)
	if err != nil {
		return nil, 0, err
	}
	responses := make([]model.ProductResponse, len(products))
	for i, product := range products {
		responses[i] = *converter.ProductsToResponse(&product)
	}
	return responses, total, nil
}

// GetProductById retrieves a product by its UUID and the user ID of the owner of the product.
// 
// Parameters:
// 
// 	* ctx: context.Context - Context for the request.
// 	* request: *model.GetProductRequest - Request containing user ID and product UUID.
// 
// Returns:
// 
// 	* 200 OK: model.ProductResponse if product is retrieved successfully.
// 
// Errors:
// 
// 	* Propagates error from use case layer if retrieval fails.
func (u *ProductUseCase) GetProductById(ctx context.Context, request *model.GetProductRequest) (*model.ProductResponse, error) {
	if err := helper.ValidateStruct(u.val, request); err != nil {
		return nil, err
	}
	response, err := u.productRepo.GetProductById(request.UserID, request.ProductUUID)
	if err != nil {
		return nil, err
	}
	return converter.ProductToResponse(response), nil
}

// UpdateProduct updates a product by its UUID and the user ID of the owner of the product.
// It first checks if the product exists and if the user is the owner of the product.
// If the product does not exist or the user is not the owner, it returns an error.
// If the product exists and the user is the owner, it updates the product and returns the updated product.
func (u *ProductUseCase) UpdateProduct(ctx context.Context, request *model.UpdateProduct) (*model.ProductResponse, error) {
	if err := helper.ValidateStruct(u.val, request); err != nil {
		return nil, err
	}
	product, err := u.productRepo.GetProductById(request.UserID, request.ProductName)
	if err != nil {
		return nil, err
	}
	if request.ProductName != "" {
		product.ProductName = request.ProductName
	}
	if request.Description != "" {
		product.Description = request.Description
	}
	if request.Price != 0 {
		product.Price = request.Price
	}
	if request.Stock != 0 {
		product.Stock = request.Stock
	}
	if request.Category != "" {
		product.Category = request.Category
	}
	if err := u.productRepo.UpdateProduct(product); err != nil {
		return nil, err
	}
	return converter.ProductToResponse(product), nil
}

// DeleteProduct deletes a product by its UUID and the user ID of the owner of the product.
// It first checks if the product exists and if the user is the owner of the product.
// If the product does not exist or the user is not the owner, it returns an error.
// If the product exists and the user is the owner, it deletes the product and returns nil.
func (u *ProductUseCase) DeleteProduct(ctx context.Context, request *model.DeleteProductRequest) error {
	if err := helper.ValidateStruct(u.val, request); err != nil {
		return err
	}
	product, err := u.productRepo.GetProductById(request.UserID, request.ProductUUID)
	if err != nil {
		return err
	}
	if err := u.productRepo.DeleteProduct(product); err != nil {
		return err
	}
	return nil
}