package service

import (
	"time"

	"github.com/yervsil/onlineShop/internal/domain"
	"github.com/yervsil/onlineShop/internal/repository"
	"github.com/yervsil/onlineShop/pkg/auth"
)

type UserSignUpInput struct {
	Name     string
	Email    string
	Phone    string
	Password string
}

type UserSignInInput struct {
	Email    string
	Password string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}

type Product struct {	
	Name string 
	Description string 
	Price float32 
	Quantity int 
}

type Admins interface {
	SignIn(user UserSignInInput) (Tokens, error)
	RefreshTokens(token string) (Tokens, error)
	CreateProduct(product Product) (int, error)
	DeleteProduct(id string) error 
	ChangeProduct(product domain.UpdateProduct, id string) error 
}

type Users interface {
	SignUp(user UserSignUpInput) (int, error)
	SignIn(user UserSignInInput) (Tokens, error)
	RefreshTokens(token string) (Tokens, error)
	AddToCart(quantity int, userID, productId string) error
	CreateOrder(userId, delivery, payment string) (int, error)
	GetOrders(userId string) ([]domain.Order, error)
	GetOrderById(userId, orderId string) (domain.Order, error)
	GetCart(userId string) (domain.Cart, error)
}

type General interface {
	GetProduct() ([]domain.Product, error)
	GetProductById(id string) (domain.Product, error)
}

type Services struct {
	Admins
	Users
	General
}

func NewService(deps Deps) *Services {
	return &Services{
		Admins: NewAdminsService(deps.Repos.Admins, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
		Users: NewUsersService(deps.Repos.Users, deps.TokenManager, deps.AccessTokenTTL, deps.RefreshTokenTTL),
		General: NewGeneralService(deps.Repos.General),
}
}

type Deps struct {
	Repos                  *repository.Repository
	TokenManager           auth.TokenManager
	AccessTokenTTL 		   time.Duration
	RefreshTokenTTL        time.Duration
}