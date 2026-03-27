package network

import (
	"net/http"

	"backend/auth/internal/constant"
	"backend/auth/internal/dto"

	"github.com/gin-gonic/gin"
)

func smsRouter(n *Network) {
	n.Router(POST, "/sms/otp/send", n.sendSMSOTP)
	n.Router(POST, "/sms/otp/verify", n.verifySMSOTP)
}

func (n *Network) sendSMSOTP(c *gin.Context) {
	var req dto.SendSMSOTPRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	result, err := n.service.SendSMSOTP(req.SessionId, req.PhoneNumber)
	if err != nil {
		c.JSON(getStatusCode(err), err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (n *Network) verifySMSOTP(c *gin.Context) {
	var req dto.VerifySMSOTPRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	result, rt, err := n.service.VerifySMSOTP(req.SessionId, req.VerificationId, req.OTP)
	if err != nil {
		c.JSON(getStatusCode(err), err.Error())
		return
	}
	if rt != "" {
		c.SetCookie("refresh_token",
			rt,
			constant.RefreshTokenTTL,
			"",
			"",
			false,
			true,
		)
	}
	c.JSON(http.StatusOK, result)
}
