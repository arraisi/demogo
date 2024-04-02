package product

import (
	productPb "github.com/arraisi/demogo-proto/golang/pb/product"
	"github.com/arraisi/demogo/config"
	"github.com/arraisi/demogo/internal/app/service/product"
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
