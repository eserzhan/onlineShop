package service

import (
	"github.com/eserzhan/onlineShop/internal/domain"
	"github.com/eserzhan/onlineShop/internal/repository"
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