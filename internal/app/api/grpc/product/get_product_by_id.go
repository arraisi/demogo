package product

import (
	"context"
	productPb "demogo/internal/proto/product"
)

func (s *Service) GetProductByID(ctx context.Context, req *productPb.GetProductIDRequest) (*productPb.GetProductIDResponse, error) {
	data, err := s.productSvc.GetProductByID(ctx, req.Id)
	if err != nil {
		return &productPb.GetProductIDResponse{
			Message: []string{err.Error()},
		}, err
	}

	return &productPb.GetProductIDResponse{
		Data: &productPb.Product{
			Id:    data.ID,
			Name:  data.Name,
			Price: data.Price,
		},
		Message: []string{"success get product by id"},
	}, nil
}
