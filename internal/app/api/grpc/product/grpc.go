package product

import (
	"demogo/config"
	"demogo/internal/app/service/product"
	productPb "demogo/internal/proto/product"
)

type Service struct {
	productPb.UnimplementedProductServiceServer
	Conf       *config.Config
	productSvc product.Service
}

func New(conf *config.Config, productSvc product.Service) Service {
	return Service{
		Conf:       conf,
		productSvc: productSvc,
	}
}
