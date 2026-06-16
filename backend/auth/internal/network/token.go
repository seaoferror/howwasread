package network

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func tokenRouter(n *Network) {
	n.Router(POST, "/refresh-token", n.refreshToken)
	n.Router(POST, "/account/logout", n.logout)
	n.Router(DELETE, "/account/delete", n.deleteAccount)
}

func (n *Network) refreshToken(c *gin.Context) {
	rt, err := c.Request.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	result, err := n.service.GenerateAccessToken(rt.Value)
	if err != nil {
		c.JSON(getStatusCode(err), err.Error())
		return
	}
	c.JSON(http.StatusOK, result)
}

func (n *Network) logout(c *gin.Context) {
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"",
		"",
		false,
		true,
	)
	c.Status(http.StatusOK)
}

func (n *Network) deleteAccount(c *gin.Context) {
	rt, err := c.Request.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, err.Error())
		return
	}
	err = n.service.DeleteAccount(rt.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	c.Status(http.StatusOK)
}
