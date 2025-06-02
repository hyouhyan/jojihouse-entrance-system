package router

import (
	"jojihouse-entrance-system/api/handler"

	"github.com/gin-gonic/gin"
)

func SetupPaymentRoutes(router *gin.Engine, handler *handler.PaymentHandler) {
	group := router.Group("/payment")
	{
		group.GET("", handler.GetAllPaymentLogs)
		group.GET("/monthly", handler.GetMonthlyPaymentLogs)
	}
}
