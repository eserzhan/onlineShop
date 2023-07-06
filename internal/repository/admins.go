package repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/eserzhan/onlineShop/internal/domain"
)

const (
	adminsTable = "Admins"
)

type AdminsRepository struct {
	db *sqlx.DB
}

func (r *AdminsRepository) GetByCredentials(email, password string) (int, error) {
	var id int
	query := fmt.Sprintf("SELECT id FROM %s where email = $1 and password = $2", adminsTable)

	err := r.db.Get(&id, query, email, password)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AdminsRepository) SetSession(userID string, session domain.Session) error {
	query := fmt.Sprintf(`
		UPDATE %s
		SET session = $1, last_login = $2
		WHERE id = $3
	`, adminsTable)
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(query, sessionJSON, time.Now(), userID)
	return err
}

func (r *AdminsRepository) GetByRefreshToken(refreshToken string) (int, error) {

	var id int
	query := fmt.Sprintf("SELECT id FROM %s WHERE session ->> 'refreshToken' = $1 AND (session ->> 'expiresAt')::timestamp > NOW() LIMIT 1", adminsTable)

	err := r.db.Get(&id, query, refreshToken)

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (r *AdminsRepository) CreateProduct(product domain.Product) (int, error){
	query := fmt.Sprintf("INSERT INTO %s (name, description, price, quantity) VALUES ($1, $2, $3, $4) RETURNING id", productsTable)

	var productID int
	err := r.db.QueryRow(query, product.Name, product.Description, product.Price, product.Quantity).Scan(&productID)
	if err != nil {
		return -1, err
	}

	return productID, nil 
}

func(r *AdminsRepository) DeleteProduct(id string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", productsTable)

	_, err := r.db.Exec(query, id)

	return err 
}

func(r *AdminsRepository) UpdateProduct(product domain.UpdateProduct, id string) error  {
	s := ""
	count := 1
	m := make([]interface{}, 0, 6)

	if product.Name != nil {
		s += fmt.Sprintf("name = $%d,", count)
		count += 1
		m = append(m, product.Name)
	}

	if product.Description != nil {
		s += fmt.Sprintf("description = $%d,", count)
		count += 1
		m = append(m, product.Description)
	}

	if product.Price != nil {
		s += fmt.Sprintf("price = $%d,", count)
		count += 1
		m = append(m, product.Price)
	}

	if product.Quantity != nil {
		s += fmt.Sprintf("quantity = $%d", count)
		count += 1
		m = append(m, product.Quantity)
	}

	if string(s[len(s) - 1]) == "," {
		s = s[:len(s) - 1]
	}

	m = append(m, id)
	query := fmt.Sprintf("UPDATE %s SET %s WHERE id = $%d", productsTable, s, count)

	_, err := r.db.Exec(query, m...)
	if err != nil {
		return err 
	}

	return nil
}

func NewAdminsRepository(db *sqlx.DB) *AdminsRepository {
	return &AdminsRepository{db: db}
}
