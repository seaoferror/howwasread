package controller

import (
	"backend/chat/internal/dto"
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

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
	c.csMutex.Lock()
	c.conns[memberId] = conn
	slog.Info("success to make connection",
		"number of current connection", len(c.conns))
	c.csMutex.Unlock()

	defer func() {
		destroy := context.Background()
		err = conn.Close(websocket.StatusNormalClosure, "")
		if err != nil {
			slog.Error("fail to close conn", "err", err)
		}
		c.service.RemoveServerIP(destroy, memberId[:])
		c.csMutex.Lock()
		delete(c.conns, memberId)
		slog.Info("success to close connection",
			"number of current connection", len(c.conns))
		c.csMutex.Unlock()
	}()

	ip, _ := getPodIP()

	init := context.Background()
	err = c.service.SetServerIP(init, memberId, ip)
	if err != nil {
		handleWebsocketError(init, conn, err)
		return
	}
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
}

func (c *Controller) RelayMessaging(ctx context.Context, id uuid.UUID, toIds [][]byte, roomId, fromId uuid.UUID, contentType, content string) error {
	var wg sync.WaitGroup
	var pushToIds [][]byte
	var mu sync.Mutex
	res := dto.MessagingResponse{
		Id:          id,
		FromId:      fromId,
		ContentType: contentType,
		Content:     content,
	}
	if roomId != uuid.Nil {
		res.RoomId = roomId
	}
	resRaw, _ := json.Marshal(res)
	for _, toId := range toIds {
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.csMutex.RLock()
			ct, ok := c.conns[uuid.UUID(toId)]
			c.csMutex.RUnlock()
			var err1 error
			if ok {
				err1 = ct.Write(ctx, websocket.MessageText, resRaw)
				if err1 != nil {
					slog.Error("fail to write payload",
						"err", err1,
					)
					err2 := ct.CloseNow()
					if err2 != nil {
						slog.Error("fail to close zombie connection", "err", err2)
					}
				}
			}
			if !ok || err1 != nil {
				mu.Lock()
				pushToIds = append(pushToIds, toId[:])
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if pushToIds != nil {
		for _, toId := range pushToIds {
			wg.Add(1)
			go func() {
				defer wg.Done()
				err1 := c.service.RemoveServerIP(ctx, toId)
				if err1 != nil {
					return
				}
			}()
		}
		wg.Wait()
		err1 := c.service.NotifyMessaging(ctx, pushToIds, roomId, fromId, contentType, content)
		if err1 != nil {
			slog.Error("fail to notify messaging", "err", err1)
			return err1
		}
	}
	return nil
}
