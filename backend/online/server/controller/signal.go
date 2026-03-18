package controller

import (
	"backend/online/server/dto"
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
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err = w.Write([]byte("incorrect header"))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		return
	}
	conversationIdRaw := r.URL.Query().Get("id")
	conversationId, err := bson.ObjectIDFromHex(conversationIdRaw)
	if err != nil {
		slog.Error("fail to parse conversation object id from raw string",
			"conversationIdRaw", conversationIdRaw)
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err = w.Write([]byte("incorrect header"))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		return
	}

	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		//OriginPatterns:     []string{"example.com"},
		InsecureSkipVerify: true,
	})
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, err = w.Write([]byte("incorrect header"))
		if err != nil {
			slog.Error("fail to write response body",
				"err", err,
			)
		}
		slog.Info("fail to accept connection", err)
		return
	}
	if c.connections == nil {
		c.connections = make(map[string]*websocket.Conn)
	}
	c.connections[memberIdRaw] = conn
	slog.Info("success to make connection",
		"number of current connection", len(c.connections))

	ip, _ := getPodIP()

	defer func() {
		destroy := context.Background()
		conn.Close(websocket.StatusNormalClosure, "")
		delete(c.connections, memberIdRaw)
		c.service.RemoveServerIP(destroy, memberId)
		c.service.RemoveParticipant(destroy, conversationId, memberId)
		slog.Info("success to close connection",
			"number of current connection", len(c.connections))
	}()

	init := context.Background()

	err = c.service.SetServerIP(init, memberId, ip)
	if err != nil {
		err = conn.Write(init, websocket.MessageText, []byte(err.Error()))
		if err != nil {
			slog.Error("fail to write payload",
				"err", err,
			)
		}
		return
	}
	pids, err := c.service.GetParticipants(init, conversationId, memberId)
	if err != nil {
		err = conn.Write(init, websocket.MessageText, []byte(err.Error()))
		if err != nil {
			slog.Error("fail to write payload",
				"err", err,
			)
		}
		return
	}
	if pids != nil {
		resp := dto.ConversationSignalResponse{FromIds: pids}
		payload, err := json.Marshal(resp)
		if err != nil {
			slog.Error("fail to marshal")
			err = conn.Write(init, websocket.MessageText, []byte("fail to get participants"))
			if err != nil {
				slog.Error("fail to write payload",
					"err", err,
				)
			}
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
		err = conn.Write(init, websocket.MessageText, []byte(err.Error()))
		if err != nil {
			slog.Error("fail to write payload",
				"err", err,
			)
		}
		return
	}
	for _, pid := range pids {
		resp := dto.ConversationSignalResponse{
			FromIds: []string{memberIdRaw},
		}
		payload, err := json.Marshal(resp)
		if err != nil {
			slog.Error("fail to marshal")
			err = conn.Write(init, websocket.MessageText, []byte("fail to get participants"))
			if err != nil {
				slog.Error("fail to write payload",
					"err", err,
				)
			}
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
		err = c.service.PublishConversationSignal(memberIdRaw, pid, []byte{})
		if err != nil {
			err = conn.Write(init, websocket.MessageText, []byte("fail to publish"))
			if err != nil {
				slog.Error("fail to write payload",
					"err", err,
				)
			}
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
			slog.Error("incorrect message type",
				"msgType", msgType,
				"data", data)
			err = conn.Write(ctx, websocket.MessageText, []byte("incorrect message type"))
			if err != nil {
				slog.Error("fail to write payload",
					"err", err)
			}
			return
		}
		if err != nil {
			slog.Error("read error",
				"err", err)
			err = conn.Write(ctx, websocket.MessageText, []byte("read error"))
			if err != nil {
				slog.Error("fail to write payload",
					"err", err,
				)
			}
			return
		}
		var req dto.ConversationSignalRequest
		err = json.Unmarshal(data, &req)
		if err != nil {
			slog.Error("fail to unmarshalling data",
				"err", err)
			err = conn.Write(ctx, websocket.MessageText, []byte("fail to unmarshal"))
			if err != nil {
				slog.Error("fail to write payload",
					"err", err,
				)
			}
			return
		}

		for _, toId := range req.ToIds {
			to, ok := c.connections[toId]
			if ok {
				resp := dto.ConversationSignalResponse{
					FromIds: []string{memberIdRaw},
					Signal:  req.Signal,
				}
				payload, err := json.Marshal(resp)
				if err != nil {
					slog.Error("fail to marshalling conversationSignal",
						"resp", resp)
					err = conn.Write(ctx, websocket.MessageText, []byte("fail to marshal"))
					if err != nil {
						slog.Error("fail to write payload",
							"err", err,
						)
					}
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
			err = c.service.PublishConversationSignal(memberIdRaw, toId, req.Signal)
			if err != nil {
				err = conn.Write(ctx, websocket.MessageText, []byte("fail to publish"))
				if err != nil {
					slog.Error("fail to write payload",
						"err", err,
					)
				}
				return
			}
		}
	}
}

func (c *Controller) RelaySignal(ctx context.Context, fromId, toId string, signal []byte) error {
	resp := dto.ConversationSignalResponse{
		FromIds: []string{fromId},
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
