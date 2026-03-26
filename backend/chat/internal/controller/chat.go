package controller

import (
	"net/http"

	"github.com/google/uuid"
)

func chatRouter(c *Controller) {
	c.Router(POST, "/chat/like", c.sendLike)
}

func (c *Controller) sendLike(w http.ResponseWriter, r *http.Request) {
	memberIdRaw := r.Header.Get("X-User-Id")
	memberId, err := uuid.Parse(memberIdRaw)
	if err != nil {
		handleParseError(w)
		return
	}

}
