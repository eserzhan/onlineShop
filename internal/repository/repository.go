package repository

import (
	//"database/sql"


	"github.com/jmoiron/sqlx"
	"github.com/eserzhan/onlineShop/internal/domain"
)

type Users interface {
	CreateUser(user domain.User) (int, error)
	GetByCredentials(email, password string) (int, error)
	GetByRefreshToken(refreshToken string) (int, error)
	SetSession(userID string, session domain.Session) error
	AddToCart(quantity int, userId, productId string) error 
	CreateOrder(userId, delivery, payment string) (int, error)
	GetOrders(userId string) ([]domain.Order, error)
	GetOrderById(userId, orderId string) (domain.Order, error)
	GetCart(userId string) (domain.Cart, error)
}

type Admins interface {
	GetByCredentials(email, password string) (int, error)
	SetSession(userID string, session domain.Session) error
	GetByRefreshToken(refreshToken string) (int, error)
	CreateProduct(product domain.Product) (int, error)
	DeleteProduct(id string) error
	UpdateProduct(product domain.UpdateProduct, id string) error 
}

type General interface {
	GetProduct() ([]domain.Product, error)
	GetProductById(id string) (domain.Product, error)
}

type Repository struct {
	Users
	Admins
	General
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Users: NewUsersRepository(db),
		Admins: NewAdminsRepository(db),
		General: NewGeneralRepository(db),
	}
}