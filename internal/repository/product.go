package repository

import (
	"github.com/abdisetiakawan/go-ecommerce/internal/entity"
	"github.com/abdisetiakawan/go-ecommerce/internal/model"
	"github.com/abdisetiakawan/go-ecommerce/internal/repository/interfaces"
	"gorm.io/gorm"
)

type ProductRepository struct {
	DB *gorm.DB
}

func NewProductRepository(DB *gorm.DB) interfaces.ProductRepository {
	return &ProductRepository{DB}
}

func (r *ProductRepository) CreateProduct(product *entity.Product) error {
	return r.DB.Create(product).Error
}

func (r *ProductRepository) GetProducts(request *model.GetProductsRequest) ([]entity.Product, int64, error) {
	var products []entity.Product
	var total int64

	query := r.DB.Model(&entity.Product{}).
		Preload("Store", func(db *gorm.DB) *gorm.DB {
			return r.DB.Where("user_id = ?", request.UserID)
		})

	if request.Search != "" {
		searchTerm := "%" + request.Search + "%"
		query = query.Where("product_name LIKE ? OR description LIKE ?", searchTerm, searchTerm)
	}
	if request.Category != "" {
		query = query.Where("category = ?", request.Category)
	}
	if request.PriceMin > 0 {
		query = query.Where("price >= ?", request.PriceMin)
	}
	if request.PriceMax > 0 {
		query = query.Where("price <= ?", request.PriceMax)
	}
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (request.Page - 1) * request.Limit
	if err := query.Limit(request.Limit).Offset(offset).Find(&products).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, model.ErrNotFound
		}
		return nil, 0, err
	}

	return products, total, nil
}

func (r *ProductRepository) GetProductById(userID uint, productUUID string) (*entity.Product, error) {
	var product entity.Product
	if err := r.DB.Model(&entity.Product{}).Preload("Store", func(db *gorm.DB) *gorm.DB {
		return r.DB.Where("user_id = ?", userID)
	}).Take(&product, "product_uuid = ?", productUUID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, model.ErrNotFound
		}
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) UpdateProduct(product *entity.Product) error {
	return r.DB.Save(product).Error
}

func (r *ProductRepository) DeleteProduct(product *entity.Product) error {
	return r.DB.Delete(product).Error
}

func (r *ProductRepository) FindProductByUUID(productUUID string) (entity.Product, error) {
	var product entity.Product
	err := r.DB.Where("product_uuid = ?", productUUID).First(&product).Error
	return product, err
}

func (r *ProductRepository) FindProductByID(productID uint) (entity.Product, error) {
	var product entity.Product
	err := r.DB.Where("id = ?", productID).First(&product).Error
	return product, err
}