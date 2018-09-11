package product

import (
	"gopkg.in/mgo.v2/bson"
	"qrapi-prd/g/x/web"
	"qrapi-prd/x/logger"
	"qrapi-prd/x/mongodb"
	"qrapi-prd/x/validator"
)

var productLog = logger.NewLogger("tbl_product")
var productTable = mongodb.NewTable("product", "prd")

type Product struct {
	mongodb.Model `bson:",inline"`
	Name          string   `bson:"name" json:"name"`
	Gallery       []string `bson:"gallery" json:"gallery"`
	Description   string   `bson:"description" json:"description"`
	CustomerID    string   `bson:"customer_id" json:"customer_id"`
}

func (product *Product) Create() error {
	err := validator.Struct(product)
	if err != nil {
		productLog.Error(err)
		return web.WrapBadRequest(err, "")
	}
	return productTable.Create(product)
}

func GetProductsByCustomer(customerID string) ([]*Product, error) {
	var products []*Product
	var err = productTable.FindWhere(bson.M{"customer_id": customerID}, &products)
	return products, err
}

func GetProductByID(id string) (*Product, error) {
	var product *Product
	var err = productTable.FindID(id, &product)
	return product, err
}
