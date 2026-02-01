package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type CategoryService struct {
	categoryRepo *repositories.CategoryRepository
}

func NewCategoryService(categoryRepo *repositories.CategoryRepository) *CategoryService {
	return &CategoryService{categoryRepo: categoryRepo}
}

func (s *CategoryService) GetAllCategories() ([]models.Category, error) {
	return s.categoryRepo.GetAll()
}

func (s *CategoryService) CreateCategory(category *models.Category) error {
	return s.categoryRepo.Create(category)
}

func (s *CategoryService) GetCategoryByID(id int) (*models.Category, error) {
	return s.categoryRepo.GetByID(id)
}

func (s *CategoryService) UpdateCategory(category *models.Category) error {
	return s.categoryRepo.Update(category)
}

func (s *CategoryService) DeleteCategory(id int) error {
	return s.categoryRepo.Delete(id)
}
