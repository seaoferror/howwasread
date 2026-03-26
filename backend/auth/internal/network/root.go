package network

import (
	"backend/auth/internal/service"

	"github.com/gin-gonic/gin"
)

type Network struct {
	service *service.Service
	engine  *gin.Engine
}

func NewNetwork(s *service.Service) *Network {
	n := &Network{
		service: s,
		engine:  gin.New(),
	}

	setGin(n.engine)

	idRouter(n)
	emailRouter(n)
	tokenRouter(n)
	smsRouter(n)

	return n
}

func (n *Network) Start() error {
	return n.engine.Run(":8081")
}
