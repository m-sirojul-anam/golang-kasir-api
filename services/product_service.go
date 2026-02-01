package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type ProductService struct {
	productRepo *repositories.ProductRepository
}

func NewProductService(productRepo *repositories.ProductRepository) *ProductService {
	return &ProductService{productRepo: productRepo}
}

func (s *ProductService) GetAllProducts() ([]models.Product, error) {
	return s.productRepo.GetAll()
}

func (s *ProductService) CreateProduct(product *models.Product) error {
	return s.productRepo.Create(product)
}

func (s *ProductService) GetProductByID(id int) (*models.Product, error) {
	return s.productRepo.GetByID(id)
}

func (s *ProductService) UpdateProduct(product *models.Product) error {
	return s.productRepo.Update(product)
}

func (s *ProductService) DeleteProduct(id int) error {
	return s.productRepo.Delete(id)
}
