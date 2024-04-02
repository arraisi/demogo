package product

import (
	"context"
	"github.com/arraisi/demogo/internal/domain/product"
)

type Service interface {
	GetProductByID(ctx context.Context, ID int64) (product.Product, error)
	SaveProduct(ctx context.Context, request product.Product) (product.Product, error)
}

type service struct {
	productRepo product.Repository
}

func New(productRepo product.Repository) Service {
	return &service{
		productRepo: productRepo,
	}
}
