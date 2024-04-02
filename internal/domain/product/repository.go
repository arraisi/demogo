package product

import (
	"context"
	"github.com/arraisi/demogo/pkg/logger"
	"github.com/arraisi/demogo/pkg/safesql"
	"github.com/elgris/sqrl"
)

type Repository interface {
	GetProductByID(ctx context.Context, ID int64) (Product, error)
	InsertProduct(ctx context.Context, request Product) (Product, error)
	UpdateProduct(ctx context.Context, request Product) (Product, error)
}

type repository struct {
	MasterDB safesql.IDatabase
	SlaveDB  safesql.IDatabase
}

func (r *repository) SaveProduct(ctx context.Context, request Product) (Product, error) {

	return request, nil
}

func New(masterDB safesql.IDatabase, slaveDB safesql.IDatabase) Repository {
	return &repository{
		SlaveDB:  slaveDB,
		MasterDB: masterDB,
	}
}

func (r *repository) GetProductByID(ctx context.Context, ID int64) (resp Product, err error) {
	columns := []string{
		"p.id",
		"p.name",
		"p.price",
	}
	query, args, err := sqrl.Select(columns...).From(resp.TableNameAlias()).Where(sqrl.Eq{"p.id": ID}).ToSql()
	if err != nil {
		return
	}

	query = r.SlaveDB.Rebind(query)

	if err = r.SlaveDB.GetContext(ctx, &resp, query, args...); err != nil {
		logger.Log.Errorf("[GetProductByID] error while run query: %v", err)
		return resp, err
	}
	return resp, nil
}

func (r *repository) InsertProduct(ctx context.Context, request Product) (Product, error) {
	sql, args, err := sqrl.Insert(request.TableName()).Columns("name", "price", "created_by", "updated_by").
		Values(request.Name, request.Price, "admin", "admin").Returning("id").
		ToSql()
	if err != nil {
		return Product{}, err
	}

	query := r.MasterDB.Rebind(sql)
	err = r.MasterDB.QueryRowxContext(ctx, query, args...).Scan(&request.ID)
	if err != nil {
		logger.Log.Errorf("[InsertProduct] error while run query: %v", err)
		return Product{}, err
	}

	return request, nil
}

func (r *repository) UpdateProduct(ctx context.Context, request Product) (Product, error) {
	sql, args, err := sqrl.Update(request.TableName()).
		Set("name", request.Name).
		Set("price", request.Price).
		Set("updated_by", "user"). // dummy
		Where(sqrl.Eq{"id": request.ID}).ToSql()
	if err != nil {
		return Product{}, err
	}

	query := r.MasterDB.Rebind(sql)
	_, err = r.MasterDB.ExecContext(ctx, query, args...)
	if err != nil {
		logger.Log.Errorf("[UpdateProduct] error while run query: %v", err)
		return Product{}, err
	}

	return request, nil
}
