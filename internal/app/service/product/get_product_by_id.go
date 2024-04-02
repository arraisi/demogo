package product

import (
	"context"
	"demogo/internal/domain/product"
)

func (s *service) GetProductByID(ctx context.Context, ID int64) (product.Product, error) {
	return s.productRepo.GetProductByID(ctx, ID)
}
