package public

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"qrapi-prd/common"
	"qrapi-prd/g/x/web"
	"qrapi-prd/o/admin"
	"qrapi-prd/o/customer"
	"qrapi-prd/o/order"
	"qrapi-prd/o/product"
	sHistory "qrapi-prd/o/scan_history"
	"qrapi-prd/x/fcm"
	"qrapi-prd/x/security"
	"strings"
)

type PublicServer struct {
	*gin.RouterGroup
	web.JsonRender
}

func NewPublicServer(parent *gin.RouterGroup, name string) *PublicServer {
	var s = PublicServer{
		RouterGroup: parent.Group(name),
	}
	s.GET("qrcode/scan", s.scanQrcode)
	s.GET("marketing/scan", s.scanMarketing)
	s.GET("product/scan", s.scanProduct)
	s.GET("order/detail", s.getOrder)
	s.POST("register", s.register)
	return &s
}

func (s *PublicServer) scanMarketing(c *gin.Context) {
	var order, err = order.GetOrderByID(c.Query("order_id"))
	web.AssertNil(err)
	var scanHistory = &sHistory.ScanHistory{
		OrderID: order.ID,
		URL:     order.URL,
	}
	scanHistory.SetID(c.Query("order_id"))
	web.AssertNil(scanHistory.Create())
	s.Success(c)
}

var errNotValidCode = errors.New("Không tìm thấy thông tin sản phẩm")

func (s *PublicServer) scanProduct(c *gin.Context) {
	var id = c.Query("id")
	var code = c.Query("code")
	var orderID = c.Query("order_id")
	// var lat = c.Query("lat")
	// var lng = c.Query("lng")
	if code != "" {
		id = id + code
	}
	decrypted, err := security.Decrypt([]byte(common.CIPHER_KEY), id)
	if err != nil {
		web.AssertNil(errNotValidCode)
	}
	//get customer, product ,scan history Info
	customerID, productID := getCustomerProductID(decrypted)
	customer, err := customer.GetCustomerByID(customerID)
	if err != nil {
		web.AssertNil(errNotValidCode)
	}
	if customer == nil {
		web.AssertNil(errNotValidCode)
	}
	product, err := product.GetProductByID(productID)
	if err != nil || product == nil {
		web.AssertNil(errNotValidCode)
	}
	//write scan history
	{
		var insertedScanHistory = &sHistory.ScanHistory{
			OrderID:   orderID,
			ProductID: productID,
		}
		//check first scan
		insertedScanHistory.SetID(id)
		insertedScanHistory.Create()
	}
	scanHistory := sHistory.GetByID(id)
	s.SendData(c, map[string]interface{}{
		"product":   product,
		"customer":  customer,
		"scan_info": scanHistory,
	})
}

func (s *PublicServer) scanQrcode(c *gin.Context) {
	var token = c.Query("token")
	var err, str = fcm.SendToOne(token, fcm.FmcMessage{Title: "Hello", Body: "Anh"})
	fmt.Println(err)
	fmt.Println(str)
	s.Success(c)
}

func (s *PublicServer) register(c *gin.Context) {
	var admin *admin.Admin
	web.AssertNil(c.BindJSON(&admin))
	web.AssertNil(admin.CreateAccount())
	s.Success(c)
}

func (s *PublicServer) getOrder(c *gin.Context) {
	var order, err = order.GetOrderByID(c.Query("order_id"))
	web.AssertNil(err)
	s.SendData(c, order)
}

//CTwHkYUk2AdRqEkLbfa6_rlK6_-x6h78ySbTNs1k4eDu9Zj1DtWlRg==

func getCustomerProductID(gid string) (string, string) {
	return strings.Split(gid, "$$")[0], strings.Split(gid, "$$")[1]
}
