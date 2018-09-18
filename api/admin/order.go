package admin

import (
	"github.com/gin-gonic/gin"
	"qrapi-prd/g/x/web"
	"qrapi-prd/o/order"
)

func (s *AdminServer) getOrders(c *gin.Context) {
	var customerID = c.Query("customer_id")
	var orders, err = order.GetOrdersByCustomer(customerID)
	web.AssertNil(err)
	s.SendData(c, orders)
}

func (s *AdminServer) deliveryOrder(c *gin.Context) {
	var orderID = c.Query("order_id")
	web.AssertNil(order.DeliveryOrder(orderID))
	s.Success(c)
}

func (s *AdminServer) createOrder(c *gin.Context) {
	var order *order.Order
	web.AssertNil(c.BindJSON(&order))
	web.AssertNil(order.Create())
	s.SendData(c, order)
}

func (s *AdminServer) updateOrder(c *gin.Context) {
	var order *order.Order
	web.AssertNil(c.BindJSON(&order))
	web.AssertNil(order.Create())
	s.SendData(c, order)
}

func (s *AdminServer) deleteOrder(c *gin.Context) {
	web.AssertNil(order.DeleteByID(c.Query("id")))
	s.Success(c)
}
