package mysqlProduct

import (
	"moori/entity"
	"moori/pkg/errormsg"
	"moori/pkg/richError"
	"moori/service/product"
	"moori/storage/mysql"
	"strings"
)

type ProductDb struct {
	conn *mysql.DB
}

func New(conn *mysql.DB) *ProductDb {
	return &ProductDb{conn: conn}
}

func (db *ProductDb) CreateProduct(product entity.Product) (entity.ProductID, error) {
	const op = "mysqlProduct.AddProduct"

	res := db.conn.Connect().Create(&product)
	if res.Error != nil {
		return 0, richError.New(op).
			SetMessage(errormsg.DataBaseInsertionError).
			SetCode(richError.UnexpectedCode).
			SetWrappedError(res.Error)
	}
	var images []entity.Image
	for _, img := range product.Images {
		images = append(images,
			entity.Image{
				ID:        0,
				Address:   img.Address,
				ProductId: product.ID,
			})
	}
	res = db.conn.Connect().CreateInBatches(&images, len(images))
	if res.Error != nil {
		return 0, richError.New(op).SetWrappedError(res.Error).
			SetCode(richError.UnexpectedCode).
			SetMessage(errormsg.DataBaseInsertionError)
	}

	if res.Error != nil {
		return 0, richError.New(op).SetWrappedError(res.Error).
			SetCode(richError.UnexpectedCode).
			SetMessage(errormsg.DataBaseInsertionError)
	}

	return product.ID, nil

}

func (db *ProductDb) AddBulky(product []entity.Product) error {
	const op = "mysqlProduct.AddBulky"
	tx := db.conn.Connect().Begin() // Start a transaction
	res := tx.CreateInBatches(&product, len(product))
	if res.Error != nil {
		return richError.New(op + ".CreateInBatches.product").
			SetWrappedError(res.Error).
			SetCode(richError.UnexpectedCode).
			SetMessage(errormsg.DataBaseInsertionError)
	}

	tx.Commit()
	return nil
}

func (db *ProductDb) Filter(req product.FilterProductsFields) ([]entity.Product, error) {
	const op = "mysqlProduct.FilterProducts"

	var products []entity.Product

	query := db.conn.Connect().Model(&entity.Product{}).Limit(50)

	if len(req.ID) > 0 {
		query = query.Where("id IN ?", req.ID)
	}

	// Filter by price range if provided
	if req.FromPrice > 0 && req.EndPrice > 0 {
		query = query.Where("current_price BETWEEN ? AND ?", req.FromPrice, req.EndPrice)
	} else if req.FromPrice > 0 {
		query = query.Where("current_price >= ?", req.FromPrice)
	} else if req.EndPrice > 0 {
		query = query.Where("current_price <= ?", req.EndPrice)
	} else if len(strings.Trim(req.Keyword, " ")) > 0 {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	query.Preload("Images")
	if err := query.Find(&products).Error; err != nil {
		return nil, richError.New(op).
			SetWrappedError(err).
			SetCode(richError.UnexpectedCode).
			SetMessage("Failed to filter products")
	}

	return products, nil
}
