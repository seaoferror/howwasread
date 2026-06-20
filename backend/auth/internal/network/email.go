package network

import (
	"backend/auth/internal/constant"
	"backend/auth/internal/dto"
	"net/http"

	"github.com/gin-gonic/gin"
)

func emailRouter(n *Network) {
	n.Router(POST, "/email/create", n.createMemberByEmail)
	n.Router(POST, "/email/login", n.loginWithEmail)
	n.Router(GET, "/email/check", n.checkEmail)
	n.Router(POST, "/email/otp/verify", n.verifyEmailOTP)
	n.Router(POST, "/email/apple", n.signInWithApple)
	n.Router(POST, "/email/google", n.signInWithGoogle)
}

func (n *Network) createMemberByEmail(c *gin.Context) {
	var req dto.SignInWithEmailRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	result, err := n.service.CreateMemberByEmail(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.JSON(getStatusCode(err), err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (n *Network) loginWithEmail(c *gin.Context) {
	var req dto.SignInWithEmailRequest

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	result, rt, err := n.service.LoginWithEmail(req.Email, req.Password)
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

func (n *Network) checkEmail(c *gin.Context) {
	var req dto.CheckEmailRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	ctx := c.Request.Context()
	ok, err := n.service.CheckEmailUsability(ctx, req.Email)
	if err != nil {
		c.JSON(getStatusCode(err), err.Error())
		return
	}
	c.JSON(http.StatusOK, ok)
}

func (n *Network) verifyEmailOTP(c *gin.Context) {
	var req dto.VerifyEmailOTPRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	result, err := n.service.VerifyEmailOTP(req.OTP, req.VerificationId)
	if err != nil {
		c.JSON(getStatusCode(err), err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (n *Network) signInWithApple(c *gin.Context) {
	var req dto.SignInWithAppleRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	responseBody, rt, err := n.service.SignInWithApple(
		req.IdentityToken,
		req.IsFirstSignIn,
	)
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
	c.JSON(http.StatusOK, responseBody)
}

func (n *Network) signInWithGoogle(c *gin.Context) {
	var req dto.SignInWithGoogleRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	responseBody, rt, err := n.service.SignInWithGoogle(c.Request.Context(), req.IdToken)
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
	c.JSON(http.StatusOK, responseBody)
}
