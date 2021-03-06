package admin

import (
	"bytes"
	"encoding/csv"
	"github.com/gin-gonic/gin"
	"qrapi-prd/common"
	"qrapi-prd/g/x/web"
	"qrapi-prd/middleware"
	"qrapi-prd/o/admin"
	"qrapi-prd/o/auth"
	"qrapi-prd/o/customer"
	"qrapi-prd/o/order"
	"qrapi-prd/x/security"
	"strconv"
)

type AdminServer struct {
	*gin.RouterGroup
	web.JsonRender
}

func NewAdminServer(parent *gin.RouterGroup, name string) *AdminServer {
	var s = AdminServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("auth/login", s.login)
	s.Use(middleware.MustBeAdmin)
	//product api
	s.POST("product/create", s.createProduct)
	s.GET("product/list", s.getProducts)
	s.GET("product/delete", s.deleteProduct)
	s.POST("product/update", s.updateProduct)
	//order api
	s.GET("order/delivery", s.deliveryOrder)
	s.GET("order/list", s.getOrders)
	s.POST("order/create", s.createOrder)
	s.POST("order/update", s.updateOrder)
	s.GET("order/delete", s.deleteOrder)
	s.GET("order/generate", s.generateCSV)
	//customer api
	s.GET("customer/list", s.getCustomers)
	return &s
}

func (s *AdminServer) login(c *gin.Context) {
	var loginInfo = struct {
		Phone    string
		Password string
	}{}
	c.BindJSON(&loginInfo)
	user, err := admin.Login(loginInfo.Phone, loginInfo.Password)
	web.AssertNil(err)
	var auth = auth.Create(user.ID, "admin")
	s.SendData(c, map[string]interface{}{
		"token":     auth.ID,
		"user_info": user,
	})
}

func (s *AdminServer) getCustomers(c *gin.Context) {
	var users, err = customer.GetCustomers()
	web.AssertNil(err)
	s.SendData(c, users)
}

func (s *AdminServer) generateCSV(c *gin.Context) {
	var quantity, _ = strconv.Atoi(c.Query("quantity"))
	var orderID = c.Query("order_id")
	var order, err = order.GetOrderByID(orderID)
	var endpointCheck = "http://app.qrcode-united.com/app/#/product/scan?type=" + order.Type + "&order_id=" + orderID + "&id="
	web.AssertNil(err)
	record := []string{"Link sản phẩm", "Mã thẻ cào"}
	b := &bytes.Buffer{}
	wr := csv.NewWriter(b)
	if order.Type == common.QRCOODE_MARKETING {
		for i := 0; i < quantity; i++ {
			record = []string{order.URL, ""}
			wr.Write(record)
		}
	} else if order.Type == common.QRCOODE_TYPE1 {
		for i := 0; i < quantity; i++ {
			var encrypted, _ = security.Encrypt([]byte(common.CIPHER_KEY), order.CustomerID+"$$"+order.ProductID)
			record = []string{endpointCheck + encrypted, encrypted[len(encrypted)-6 : len(encrypted)]}
			wr.Write(record)
		}
	} else {
		for i := 0; i < quantity; i++ {
			var encrypted, _ = security.Encrypt([]byte(common.CIPHER_KEY), order.CustomerID+"$$"+order.ProductID)
			record = []string{endpointCheck + encrypted[0:len(encrypted)-6], encrypted[len(encrypted)-6 : len(encrypted)]}
			wr.Write(record)
		}
	}
	wr.Flush()                                        // writes the csv writer data to  the buffered data io writer(b(bytes.buffer))
	c.Writer.Header().Set("Content-Type", "text/csv") // setting the content type header to text/csv
	c.Writer.Header().Set("Content-Type", "text/csv")
	c.Writer.Header().Set("Content-Disposition", "attachment;filename="+order.Name+".csv")
	c.Writer.Write(b.Bytes())
}
