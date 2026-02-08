package repositories

import (
	"database/sql"
	"fmt"
	"kasir-api/models"
	"strings"
	"time"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (repo *TransactionRepository) CreateTransaction(items []models.CheckoutItem) (*models.Transaction, error) {
	tx, err := repo.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	totalAmount := 0
	details := make([]models.TransactionDetail, 0)

	for _, item := range items {
		var productPrice, stock int
		var productName string

		err := tx.QueryRow(`
			SELECT
				name, price, stock
			FROM
				products
			WHERE
				id = $1
		`, item.ProductID).Scan(&productName, &productPrice, &stock)
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("Product ID %d not found", item.ProductID)
		}
		if err != nil {
			return nil, err
		}

		subTotal := productPrice * item.Quantity
		totalAmount += subTotal

		_, err = tx.Exec(`
			UPDATE
				products
			SET
				stock = stock - $1
			WHERE
				id = $2
		`, item.Quantity, item.ProductID)
		if err != nil {
			return nil, err
		}

		details = append(details, models.TransactionDetail{
			ProductID:   item.ProductID,
			ProductName: productName,
			Quantity:    item.Quantity,
			Subtotal:    subTotal,
		})
	}

	var transactionID int
	err = tx.QueryRow(`
		INSERT INTO
			transactions
			(total_amount)
		VALUES
			($1)
		RETURNING id
	`, totalAmount).Scan(&transactionID)
	if err != nil {
		return nil, err
	}

	values := []string{}
	args := []interface{}{}
	for i, detail := range details {
		pos := i * 4
		values = append(values, fmt.Sprintf("($%d, $%d, $%d, $%d)",
			pos+1, pos+2, pos+3, pos+4))

		args = append(args, transactionID, detail.ProductID,
			detail.Quantity, detail.Subtotal)
	}

	query := fmt.Sprintf(`
		INSERT INTO transaction_details 
			(transaction_id, product_id, quantity, subtotal) 
		VALUES %s
	`, strings.Join(values, ","))

	start := time.Now()
	_, err = tx.Exec(query, args...)
	fmt.Printf("Batch insert took: %v\n", time.Since(start))

	if err != nil {
		return nil, fmt.Errorf("batch insert failed: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return &models.Transaction{
		ID:          transactionID,
		TotalAmount: totalAmount,
		Details:     details,
	}, nil
}
