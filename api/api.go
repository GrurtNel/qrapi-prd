package api

import (
	"qrapi-prd/api/admin"
	"qrapi-prd/api/auth"
	"qrapi-prd/api/customer"
	"qrapi-prd/api/guest"
	"qrapi-prd/api/public"

	"github.com/gin-gonic/gin"
)

func NewApiServer(root *gin.RouterGroup) {
	admin.NewAdminServer(root, "admin")
	auth.NewAuthServer(root, "auth")
	public.NewPublicServer(root, "public")
	guest.NewGuestServer(root, "guest")
	customer.NewCustomerServer(root, "customer")
}
