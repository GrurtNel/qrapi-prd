package guest

import (
	"github.com/gin-gonic/gin"
	"qrapi-prd/g/x/web"
	"qrapi-prd/o/auth"
	"qrapi-prd/o/customer"
)

type GuestServer struct {
	*gin.RouterGroup
	web.JsonRender
}

func NewGuestServer(parent *gin.RouterGroup, name string) *GuestServer {
	var s = GuestServer{
		RouterGroup: parent.Group(name),
	}
	s.POST("register", s.register)
	return &s
}

func (s *GuestServer) register(c *gin.Context) {
	var u *customer.Customer
	c.BindJSON(&u)
	u.Role = customer.CUSTOMER
	web.AssertNil(u.Create())
	var auth = auth.Create(u.ID, u.Role)
	s.SendData(c, map[string]interface{}{
		"token":     auth.ID,
		"user_info": u,
	})
}
