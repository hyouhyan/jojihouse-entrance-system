package handler

import (
	"fmt"
	"jojihouse-entrance-system/api/model/request"
	"jojihouse-entrance-system/internal/model"
	"jojihouse-entrance-system/internal/service"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type KaisukenHandler struct {
	userPortalService      *service.UserPortalService
	adminManagementService *service.AdminManagementService
}

func NewKaisukenHandler(
	userPortalService *service.UserPortalService,
	adminManagementService *service.AdminManagementService,
) *KaisukenHandler {
	return &KaisukenHandler{
		userPortalService:      userPortalService,
		adminManagementService: adminManagementService,
	}
}

// @Summary 回数券購入
// @Tags 回数券
// @Description ユーザーが回数券を購入します
// @Accept json
// @Produce json
// @Param kaisuken body request.BuyKaisuken true "回数券購入のリクエストデータ"
// @Success 200 {object} map[string]string
// @Router /kaisuken [post]
func (h *KaisukenHandler) BuyKaisuken(c *gin.Context) {
	var req request.BuyKaisuken
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		log.Print(err)
		return
	}

	description := fmt.Sprintf("回数券購入 %d回分 %d円", req.Count, req.Amount)
	err := h.adminManagementService.IncreaseRemainingEntries(req.UserID, req.Count, description, req.Receiver)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not increase remaining entries"})
		log.Print(err)
	}

	// 支払いログの作成
	paymentLog := &model.PaymentLog{
		UserID:      req.UserID,
		Amount:      req.Amount,
		Time:        time.Now(),
		Description: description,
		Payway:      req.Payway,
	}
	err = h.adminManagementService.CreatePaymentLog(paymentLog)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create payment log"})
		log.Print(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Success", "time": paymentLog.Time.Format(time.RFC3339)})
}
