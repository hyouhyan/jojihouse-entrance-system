package handler

import (
	"jojihouse-entrance-system/internal/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	adminManagementService *service.AdminManagementService
}

func NewPaymentHandler(
	adminManagementService *service.AdminManagementService,
) *PaymentHandler {
	return &PaymentHandler{
		adminManagementService: adminManagementService,
	}
}

func (h *PaymentHandler) GetAllPaymentLogs(c *gin.Context) {
	lastID := c.Query("last_id") // クエリパラメータから lastID を取得
	limitStr := c.Query("limit") // クエリパラメータから limit を取得

	// デフォルトの取得件数を設定（limit が指定されていなければ 10）
	limit := int64(10)
	if limitStr != "" {
		parsedLimit, err := strconv.Atoi(limitStr)
		if err == nil && parsedLimit > 0 {
			limit = int64(parsedLimit)
		}
	}

	paymentLogs, err := h.adminManagementService.GetAllPaymentLogs(lastID, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get payment log"})
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"payment_logs": paymentLogs})
}
