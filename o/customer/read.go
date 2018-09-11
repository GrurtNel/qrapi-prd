package customer

import (
	"gopkg.in/mgo.v2/bson"
	"qrapi-prd/g/x/web"
)

var ErrMismatchedHashAndPassword = "crypto/bcrypt: hashedPassword is not the hash of the given password"

func GetUsers() ([]*Customer, error) {
	var customer []*Customer
	err := customerTable.FindWhere(bson.M{
		"dtime": bson.M{
			"$ne": 0,
		},
	}, &customer)
	return customer, err
}

func Login(phone, pwd string) (*Customer, error) {
	var customer *Customer
	var query = bson.M{"phone": phone}
	err := customerTable.FindOne(query, &customer)
	if err != nil {
		if err.Error() == "not found" {
			return nil, web.BadRequest("Sai tên đăng nhập hoặc mật khẩu")
		}
		return nil, err
	}
	if err := Password(pwd).ComparePassword(customer.HashedPassword); err != nil {
		if err.Error() == ErrMismatchedHashAndPassword {
			customerLog.Error(err)
			return nil, web.BadRequest("Sai tên đăng nhập hoặc mật khẩu")
		}
		return nil, err
	}
	// if !customer.Activated {
	// 	return nil, web.BadRequest("Tài khoản chưa được kích hoạt vui lòng liên hệ để được kích hoạt")
	// }
	return customer, nil
}

func GetCustomerByID(id string) (*Customer, error) {
	var customer *Customer
	var err = customerTable.FindID(id, &customer)
	return customer, err
}

func GetCustomers() ([]*Customer, error) {
	var customer []*Customer
	var err = customerTable.FindAll(&customer)
	return customer, err
}
