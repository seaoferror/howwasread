package controller

import (
	"backend/chat/internal/dto"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

func messagingRouter(c *Controller) {
	c.Router(POST, "/chat/like", c.sendLike)
	c.Router(GET, "/chat/messaging/connect", c.connectMessaging)
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

func (c *Controller) connectMessaging(w http.ResponseWriter, r *http.Request) {
	memberIdRaw := r.Header.Get("X-User-Id")
	memberId, err := uuid.Parse(memberIdRaw)
	if err != nil {
		slog.Error("fail to parse member id from raw string",
			"err", err,
			"memberIdRaw", memberIdRaw)
		handleError(w, errors.New("fail to parse"))
		return
	}
	slog.Info("try to make connection", "memberId", memberIdRaw)
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		//OriginPatterns:     []string{"example.com"},
		InsecureSkipVerify: true,
	})
	if err != nil {
		slog.Error("fail to accept connection", err)
		handleError(w, errors.New("fail to accept ws connection"))
		return
	}
	if c.connections == nil {
		c.connections = make(map[uuid.UUID]*websocket.Conn)
	}
	c.connections[memberId] = conn
	slog.Info("success to make connection",
		"number of current connection", len(c.connections))

	defer func() {
		destroy := context.Background()
		conn.Close(websocket.StatusNormalClosure, "")
		delete(c.connections, memberId)
		c.service.RemoveServerIP(destroy, memberId)
		slog.Info("success to close connection",
			"number of current connection", len(c.connections))
	}()

	ip, _ := getPodIP()

	init := context.Background()
	err = c.service.SetServerIP(init, memberId, ip)
	if err != nil {
		handleWebsocketError(init, conn, err)
		return
	}
	//TODO: maybe get all recent chat which aren't fetched yet

	for {
		ctx := context.Background()
		msgType, data, err := conn.Read(ctx)
		if websocket.CloseStatus(err) != -1 {
			slog.Error("Connection closed",
				"err", err)
			return
		}
		if msgType != websocket.MessageText {
			slog.Error("incorrect payload type",
				"msgType", msgType,
				"data", data)
			handleWebsocketError(ctx, conn, errors.New("incorrect payload type"))
			return
		}
		if err != nil {
			slog.Error("read error",
				"err", err)
			handleWebsocketError(ctx, conn, errors.New("read error"))
			return
		}
		var req dto.MessagingRequest
		err = json.Unmarshal(data, &req)
		if err != nil {
			handleWebsocketError(ctx, conn, errors.New("fail to unmarshal"))
			return
		}
		if req.ToIdType == "personal" {
			to, ok := c.connections[req.ToId]
			if ok {
				resp := dto.MessagingResponse{
					FromId:      memberId,
					ContentType: req.ContentType,
					Content:     req.Content,
				}
				payload, err := json.Marshal(resp)
				if err != nil {
					slog.Error("fail to marshal",
						"err", err)
					handleWebsocketError(ctx, conn, errors.New("fail to marshal"))
					return
				}
				err = to.Write(ctx, websocket.MessageText, payload)
				if err != nil {
					slog.Error("fail to write payload",
						"err", err)
					return
				}
				continue
			}
			err = c.service.PublishPersonalMessaging(req.ToId, memberId, req.ContentType, req.Content)
			if err != nil {
				handleWebsocketError(ctx, conn, errors.New("fail to publish"))
				return
			}
		}
	}
}

func (c *Controller) RelayMessaging(ctx context.Context, toIds []uuid.UUID, roomId, fromId uuid.UUID, contentType, content string) error {
	res := dto.MessagingResponse{
		FromId:      fromId,
		ContentType: contentType,
		Content:     content,
	}
	if roomId != uuid.Nil {
		res.RoomId = roomId
	}
	resRaw, err := json.Marshal(res)
	if err != nil {
		slog.Error("fail to marshal MessagingResponse",
			"err", err)
		return err
	}

	for _, toId := range toIds {
		func() {
			err = c.connections[toId].Write(ctx, websocket.MessageText, resRaw)
			if err != nil {
				slog.Error("fail to write payload",
					"err", err,
				)
				//TODO: notification fcm?
				return
			}
		}()
	}
	return nil
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
