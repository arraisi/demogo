package product

import (
	"context"
	"demogo/internal/domain/product"
)

func (s *service) SaveProduct(ctx context.Context, request product.Product) (product.Product, error) {
	if request.ID == 0 {
		return s.productRepo.InsertProduct(ctx, request)
	}

	return s.productRepo.UpdateProduct(ctx, request)
}
