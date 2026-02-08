package services

import (
	"kasir-api/models"
	"kasir-api/repositories"
)

type TransactionService struct {
	transactionRepo *repositories.TransactionRepository
}

func NewTransactionService(transactionRepo *repositories.TransactionRepository) *TransactionService {
	return &TransactionService{transactionRepo: transactionRepo}
}

func (s *TransactionService) Checkout(items []models.CheckoutItem) (*models.Transaction, error) {
	return s.transactionRepo.CreateTransaction(items)
}
