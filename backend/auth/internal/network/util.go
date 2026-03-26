package network

import (
	"backend/auth/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HTTPMethod int

const (
	GET HTTPMethod = iota
	POST
	DELETE
	PUT
)

func getStatusCode(err error) int {
	switch err {
	case service.ErrInternalServer:
		return http.StatusInternalServerError

	case service.ErrSignInWithApple:
		return http.StatusUnauthorized
	case service.ErrCheckEmail:
		return http.StatusBadRequest
	case service.ErrSignUpWithEmail:
		return http.StatusBadRequest
	case service.ErrLoginWithEmail:
		return http.StatusUnauthorized
	case service.ErrSendEmailOTP:
		return http.StatusBadRequest
	case service.ErrVerifyEmailOTP:
		return http.StatusUnauthorized
	case service.ErrSendSMSOTP:
		return http.StatusBadRequest
	case service.ErrPhoneNumberAlreadyLinked:
		return http.StatusBadRequest
	case service.ErrVerifySMSOTP:
		return http.StatusUnauthorized
	case service.ErrGenerateToken:
		return http.StatusUnauthorized
	}
	return http.StatusBadRequest
}

func setGin(engine *gin.Engine) {
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
}

func (n *Network) Router(httpMethod HTTPMethod, path string, handler ...gin.HandlerFunc) {
	e := n.engine.Group("/auth")

	switch httpMethod {
	case GET:
		e.GET(path, handler...)
	case POST:
		e.POST(path, handler...)
	case PUT:
		e.PUT(path, handler...)
	case DELETE:
		e.DELETE(path, handler...)

	default:
		panic("This HTTP method is not registered")
	}
}
