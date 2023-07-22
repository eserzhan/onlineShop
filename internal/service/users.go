package service

import (
	"crypto/sha1"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/yervsil/onlineShop/internal/domain"
	"github.com/yervsil/onlineShop/internal/repository"
	"github.com/yervsil/onlineShop/pkg/auth"
)

const (
	salt = "3p4tm24t1fdvsdk,v.g,rlw,gs"
)

type UsersService struct {
	repo         		   repository.Users
	tokenManager 		   auth.TokenManager
	accessTokenTTL         time.Duration
	refreshTokenTTL        time.Duration
}

func NewUsersService(repo repository.Users, tokenManager auth.TokenManager, accessTokenTTL, refreshTokenTTL time.Duration) *UsersService {
	return &UsersService{
		repo:         repo,
		tokenManager: tokenManager,
		accessTokenTTL: accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}
}

func (s *UsersService) SignUp(input UserSignUpInput) (int, error) {
	passwordHash := hashPasswordWithSalt(input.Password)

	user := domain.User{
		Name:         input.Name,
		Password:     passwordHash,
		Phone:        input.Phone,
		Email:        input.Email,
		RegisteredAt: time.Now(),
		LastVisitAt:  time.Now(),
	}

	return s.repo.CreateUser(user)
}

func (s *UsersService) SignIn(input UserSignInInput) (Tokens, error) {
	passwordHash := hashPasswordWithSalt(input.Password)

	id, err := s.repo.GetByCredentials(input.Email, passwordHash)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(strconv.Itoa(id))
}

func (s *UsersService) createSession(userId string) (Tokens, error) {
	var (
		res Tokens
		err error
	)

	res.AccessToken, err = s.tokenManager.NewJWT(userId, "user", s.accessTokenTTL)
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

func(s *UsersService) RefreshTokens(refreshToken string) (Tokens, error) {
	id, err := s.repo.GetByRefreshToken(refreshToken)
	if err != nil {
		return Tokens{}, err
	}

	return s.createSession(strconv.Itoa(id))
}

func(s *UsersService) AddToCart(quantity int, userId,  productId string) error {
	return s.repo.AddToCart(quantity, userId, productId)
}

func(s *UsersService) CreateOrder(userId, delivery, payment string) (int, error) {
	return s.repo.CreateOrder(userId, delivery, payment)
}

func hashPasswordWithSalt(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password + salt))
	hashBytes := hash.Sum(nil)
	return hex.EncodeToString(hashBytes)
}

func(s *UsersService) GetOrders(userId string) ([]domain.Order, error) {
	return s.repo.GetOrders(userId)
}

func(s *UsersService) GetOrderById(userId, orderId string) (domain.Order, error) {
	return s.repo.GetOrderById(userId, orderId)
}

func(s *UsersService) GetCart(userId string) (domain.Cart, error) {
	return s.repo.GetCart(userId)
}