package controller

import (
	"backend/chat/internal/dto"
	"encoding/json"
	"errors"
	"net"
	"net/http"

	"github.com/google/uuid"
)

func messagingRouter(c *Controller) {
	c.Router(POST, "/chat/like", c.sendLike)
	c.Router(GET, "/chat/messaging/connect", c.connectMessaging)
	c.Router(GET, "/chat/messaging/recent", c.getRecentMessages)
}

func (c *Controller) sendLike(w http.ResponseWriter, r *http.Request) {
	memberIdRaw := r.Header.Get("X-User-Id")
	memberId, err := uuid.Parse(memberIdRaw)
	if err != nil {
		handleError(w, errors.New("fail to parse"))
		return
	}
	var req dto.SendLikeRequest
	err = json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		handleError(w, errors.New("fail to parse"))
		return
	}
	err = c.service.SendLike(r.Context(), memberId, req.ToId)
	if err != nil {
		handleError(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (c *Controller) getRecentMessages(w http.ResponseWriter, r *http.Request) {

}

// getPodIp will replace with k8s configmap pod ip
func getPodIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), nil
		}
	}
	return "", errors.New("IP not found")
}
