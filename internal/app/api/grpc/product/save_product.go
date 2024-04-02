package product

import (
	"context"
	productPb "github.com/arraisi/demogo-proto/golang/pb/product"
	"github.com/arraisi/demogo/internal/domain/product"
)

func (s *Service) SaveProduct(ctx context.Context, request *productPb.SaveProductRequest) (*productPb.SaveProductResponse, error) {
	saveProduct, err := s.productSvc.SaveProduct(ctx, product.Product{
		Name:  request.Data.Name,
		Price: request.Data.Price,
		ID:    request.Data.Id,
	})
	if err != nil {
		return &productPb.SaveProductResponse{}, err
	}

	return &productPb.SaveProductResponse{
		Data: &productPb.Product{
			Id:    saveProduct.ID,
			Name:  saveProduct.Name,
			Price: saveProduct.Price,
		}, Message: []string{"success save product"},
	}, nil
}
