package controller

import (
	"backend/online/internal/dto"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"

	"github.com/coder/websocket"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (c *Controller) joinConversation(w http.ResponseWriter, r *http.Request) {
	slog.Info("try to make connection")
	memberIdRaw := r.Header.Get("X-User-Id")
	memberId, err := uuid.Parse(memberIdRaw)
	if err != nil {
		slog.Error("fail to parse member id from raw string",
			"err", err,
			"memberIdRaw", memberIdRaw)
		handleError(w, errors.New("fail to parse"))
		return
	}
	conversationIdRaw := r.URL.Query().Get("id")
	conversationId, err := bson.ObjectIDFromHex(conversationIdRaw)
	if err != nil {
		slog.Error("fail to parse conversation object id from raw string",
			"conversationIdRaw", conversationIdRaw)
		handleError(w, errors.New("fail to parse"))
		return
	}

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

	ip, _ := getPodIP()

	defer func() {
		destroy := context.Background()
		conn.Close(websocket.StatusNormalClosure, "")
		delete(c.connections, memberId)
		c.service.RemoveServerIP(destroy, memberId)
		c.service.RemoveParticipant(destroy, conversationId, memberId)
		slog.Info("success to close connection",
			"number of current connection", len(c.connections))
	}()

	init := context.Background()

	err = c.service.SetServerIP(init, memberId, ip)
	if err != nil {
		handleWebsocketError(init, conn, err)
		return
	}
	pids, err := c.service.GetParticipants(init, conversationId, memberId)
	if err != nil {
		handleWebsocketError(init, conn, err)
		return
	}
	if pids != nil {
		resp := dto.ConversationSignalResponse{FromIds: pids}
		payload, err := json.Marshal(resp)
		if err != nil {
			slog.Error("fail to marshal")
			handleWebsocketError(init, conn, err)
			return
		}
		err = conn.Write(init, websocket.MessageText, payload)
		if err != nil {
			slog.Error("fail to write payload",
				"err", err,
			)
			return
		}
	}

	err = c.service.AddParticipant(init, conversationId, memberId)
	if err != nil {
		handleWebsocketError(init, conn, err)
		return
	}
	for _, pid := range pids {
		resp := dto.ConversationSignalResponse{
			FromIds: []uuid.UUID{memberId},
		}
		payload, err := json.Marshal(resp)
		if err != nil {
			slog.Error("fail to marshal")
			handleWebsocketError(init, conn, errors.New("fail to get participant"))
			return
		}
		p, ok := c.connections[pid]
		if ok {
			err = p.Write(init, websocket.MessageText, payload)
			if err != nil {
				slog.Error("fail to write payload",
					"err", err,
				)
				return
			}
			continue
		}
		err = c.service.PublishConversationSignal(memberId, pid, []byte{})
		if err != nil {
			handleWebsocketError(init, conn, errors.New("fail to publish"))
			return
		}
	}
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
		var req dto.ConversationSignalRequest
		err = json.Unmarshal(data, &req)
		if err != nil {
			slog.Error("fail to unmarshalling data",
				"err", err)
			handleWebsocketError(ctx, conn, errors.New("incorrect data"))
			return
		}
		for _, toId := range req.ToIds {
			to, ok := c.connections[toId]
			if ok {
				resp := dto.ConversationSignalResponse{
					FromIds: []uuid.UUID{memberId},
					Signal:  req.Signal,
				}
				payload, err := json.Marshal(resp)
				if err != nil {
					slog.Error("fail to marshalling conversationSignal",
						"err", err)
					handleWebsocketError(ctx, conn, errors.New("something went wrong"))
					return
				}
				err = to.Write(ctx, websocket.MessageText, payload)
				if err != nil {
					slog.Error("fail to write payload",
						"err", err,
					)
					return
				}
				continue
			}
			err = c.service.PublishConversationSignal(memberId, toId, req.Signal)
			if err != nil {
				handleWebsocketError(ctx, conn, errors.New("fail to publish"))
				return
			}
		}
	}
}

func (c *Controller) RelaySignal(ctx context.Context, fromId, toId uuid.UUID, signal []byte) error {
	resp := dto.ConversationSignalResponse{
		FromIds: []uuid.UUID{fromId},
		Signal:  signal,
	}
	payload, err := json.Marshal(resp)
	if err != nil {
		slog.Error("fail to marshal ConversationSignalResponse",
			"err", err)
		return err
	}

	err = c.connections[toId].Write(ctx, websocket.MessageText, payload)
	if err != nil {
		slog.Error("fail to write payload",
			"err", err,
		)
		return err
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
