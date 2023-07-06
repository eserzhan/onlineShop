package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/eserzhan/onlineShop/internal/domain"
)

// const (
// 	ProductTable = "Products"
// )
type GeneralRepository struct {
	db *sqlx.DB
}

func(r *GeneralRepository) GetProduct() ([]domain.Product, error){
	var Product []domain.Product

	query := fmt.Sprintf("SELECT name, description, price, quantity FROM %s ", productsTable)

	err := r.db.Select(&Product, query)

	if err != nil {
		return []domain.Product{}, err
	}

	return Product, nil
}

func(r *GeneralRepository) GetProductById(id string) (domain.Product, error){
	var good domain.Product

	query := fmt.Sprintf("SELECT name, description, price, quantity FROM %s WHERE id = $1", productsTable)

	err := r.db.Get(&good, query, id)

	if err != nil {
		return domain.Product{}, err
	}

	return good, nil
}

func NewGeneralRepository(db *sqlx.DB) *GeneralRepository{
	return &GeneralRepository{db: db }
}

