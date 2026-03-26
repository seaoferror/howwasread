package network

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func idRouter(n *Network) {
	n.Router(GET, "/my-id", n.getMyId)
}

func (n *Network) getMyId(c *gin.Context) {
	id := c.Request.Header.Get("X-User-Id")
	c.JSON(http.StatusOK, map[string]string{"id": id})
}
