package customer

import (
	"github.com/gin-gonic/gin"
	"qrapi-prd/g/x/web"
	"qrapi-prd/middleware"
	"qrapi-prd/o/order"
	"qrapi-prd/o/product"
)

type CustomerServer struct {
	*gin.RouterGroup
	web.JsonRender
}

func NewCustomerServer(parent *gin.RouterGroup, name string) *CustomerServer {
	var s = CustomerServer{
		RouterGroup: parent.Group(name),
	}
	s.Use(middleware.MustBeCustomer)
	s.POST("product/create", s.createProduct)
	s.GET("product/list", s.getProducts)
	s.GET("product/delete", s.deleteProduct)
	s.POST("product/update", s.updateProduct)
	s.POST("order/create", s.createOrder)
	s.POST("order/update", s.updateOrder)
	s.GET("order/delete", s.deleteOrder)
	s.GET("order/list", s.getOrders)
	return &s
}

func (s *CustomerServer) createProduct(c *gin.Context) {
	var product *product.Product
	c.BindJSON(&product)
	web.AssertNil(product.Create())
	s.SendData(c, product)
}

func (s *CustomerServer) updateProduct(c *gin.Context) {
	var product *product.Product
	c.BindJSON(&product)
	web.AssertNil(product.Create())
	s.SendData(c, product)
}

func (s *CustomerServer) deleteProduct(c *gin.Context) {
	web.AssertNil(product.DeleteByID(c.Query("id")))
	s.Success(c)
}

func (s *CustomerServer) getProducts(c *gin.Context) {
	var customerID = c.MustGet("user_id")
	var products, err = product.GetProductsByCustomer(customerID.(string))
	web.AssertNil(err)
	s.SendData(c, products)
}

func (s *CustomerServer) createOrder(c *gin.Context) {
	var order *order.Order
	web.AssertNil(c.BindJSON(&order))
	web.AssertNil(order.Create())
	s.SendData(c, order)
}

func (s *CustomerServer) updateOrder(c *gin.Context) {
	var order *order.Order
	web.AssertNil(c.BindJSON(&order))
	web.AssertNil(order.Create())
	s.SendData(c, order)
}

func (s *CustomerServer) deleteOrder(c *gin.Context) {
	web.AssertNil(order.DeleteByID(c.Query("id")))
	s.Success(c)
}

func (s *CustomerServer) getOrders(c *gin.Context) {
	var customerID = c.MustGet("user_id")
	var orders, err = order.GetOrdersByCustomer(customerID.(string))
	web.AssertNil(err)
	s.SendData(c, orders)
}
