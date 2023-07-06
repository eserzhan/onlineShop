package service

import (
	"strconv"
	"time"

	"github.com/eserzhan/onlineShop/internal/domain"
	"github.com/eserzhan/onlineShop/internal/repository"
	"github.com/eserzhan/onlineShop/pkg/auth"
)

type AdminsService struct {
	repo            repository.Admins
	tokenManager    auth.TokenManager
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

func NewAdminsService(repo repository.Admins, tokenManager auth.TokenManager, accessTokenTTL, refreshTokenTTL time.Duration) *AdminsService {
	return &AdminsService{repo: repo, tokenManager: tokenManager, accessTokenTTL: accessTokenTTL, refreshTokenTTL: refreshTokenTTL}
}

func (s *AdminsService) SignIn(user UserSignInInput) (Tokens, error) {
	id, err := s.repo.GetByCredentials(user.Email, user.Password)
	if err != nil {
		return Tokens{}, err 
	}

	return s.createSession(strconv.Itoa(id))
}

func (s *AdminsService) createSession(userId string) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userId, "admin", s.accessTokenTTL)
	if err != nil {
		return res, err
	}

	res.RefreshToken, err = s.tokenManager.NewRefreshToken()
	if err != nil {
		return res, err
	}

	session := domain.Session{
		RefreshToken: res.RefreshToken,
		ExpiresAt:    time.Now().Add(s.refreshTokenTTL),
	}

	err = s.repo.SetSession(userId, session)

	return res, err
}

func(s *AdminsService) RefreshTokens(refreshToken string) (Tokens, error) {
	id, err := s.repo.GetByRefreshToken(refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(strconv.Itoa(id))
}

func(s *AdminsService) CreateProduct(product Product) (int, error) {
	item := domain.Product{
		Name:         	 product.Name,
		Description:     product.Description,
		Price:        	 product.Price,
		Quantity:        product.Quantity,
	}

	return s.repo.CreateProduct(item)
}

func(s *AdminsService) ChangeProduct(product domain.UpdateProduct, id string) error {
	return s.repo.UpdateProduct(product, id)
}

func(s *AdminsService) DeleteProduct(id string) error {
	return s.repo.DeleteProduct(id)
}