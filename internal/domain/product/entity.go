package product

type Product struct {
	ID    int64   `db:"id" json:"id"`
	Name  string  `db:"name" json:"name"`
	Price float64 `db:"price" json:"price"`
}

func (p Product) TableName() string {
	return "products"
}

func (p Product) TableNameAlias() string {
	return "products p"
}
