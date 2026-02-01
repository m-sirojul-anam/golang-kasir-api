package services

import (
	"errors"
	"kasir-api/models"
	"kasir-api/repositories"
)

type ProductService struct {
	productRepo  *repositories.ProductRepository
	categoryRepo *repositories.CategoryRepository
}

func NewProductService(productRepo *repositories.ProductRepository, categoryRepo *repositories.CategoryRepository) *ProductService {
	return &ProductService{productRepo: productRepo, categoryRepo: categoryRepo}
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return s.productRepo.GetAll()
}

func (s *ProductService) CreateProduct(product *models.CreateProductRequest) error {

	if product.CategoryID != nil {
		category, err := s.categoryRepo.GetByID(*product.CategoryID)
		if err != nil || category == nil {
			return errors.New("Invalid category ID.")
		}
	}

	return s.productRepo.Create(product)
}

func (s *ProductService) GetProductByID(id int) (*models.Product, error) {
	return s.productRepo.GetByID(id)
}

func (s *ProductService) UpdateProduct(product *models.CreateProductRequest) error {

	if product.CategoryID != nil {
		category, err := s.categoryRepo.GetByID(*product.CategoryID)
		if err != nil || category == nil {
			return errors.New("Invalid category ID.")
		}
	}

	return s.productRepo.Update(product)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.productRepo.Delete(id)
}
