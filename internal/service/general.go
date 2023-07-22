package service

import (
	"github.com/yervsil/onlineShop/internal/domain"
	"github.com/yervsil/onlineShop/internal/repository"
)

type GeneralService struct {
	repo repository.General
}

func(s *GeneralService) GetProduct() ([]domain.Product, error){
	return s.repo.GetProduct()
}

func(s *GeneralService) GetProductById(id string) (domain.Product, error){
	return s.repo.GetProductById(id)
}

func NewGeneralService(repo repository.General) *GeneralService {
	return &GeneralService{repo: repo}
}