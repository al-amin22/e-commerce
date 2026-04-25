package service

import (
	"context"
	"strings"

	"ecommerce-backend/internal/domain"
	"ecommerce-backend/internal/repository"
)

type ProductService interface {
	Create(ctx context.Context, product *domain.Product) error
	List(ctx context.Context) ([]domain.Product, error)
}

type productService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) Create(ctx context.Context, product *domain.Product) error {
	product.Name = strings.TrimSpace(product.Name)
	if product.Stock < 0 {
		product.Stock = 0
	}
	return s.repo.Create(ctx, product)
}

func (s *productService) List(ctx context.Context) ([]domain.Product, error) {
	return s.repo.List(ctx)
}
