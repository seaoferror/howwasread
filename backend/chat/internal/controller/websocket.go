package controller

import (
	"backend/chat/internal/dto"
	"backend/common/payload"
	"bytes"
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"sync"

	"github.com/coder/websocket"
	"github.com/google/uuid"
)

func (c *Controller) connectMessaging(w http.ResponseWriter, r *http.Request) {
	memberId, err := uuid.Parse(r.Header.Get("X-User-Id"))
	if err != nil {
		slog.Error("fail to parse member id from header",
			"err", err)
		handleError(w, errors.New("fail to parse"))
		return
	}
	deviceId, err := uuid.Parse(r.Header.Get("Device-Id"))
	if err != nil {
		slog.Error("fail to parse device id from header",
			"err", err)
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
	c.csMutex.Lock()
	if c.conns[memberId] == nil {
		c.conns[memberId] = make(map[uuid.UUID]*websocket.Conn)
	}
	c.conns[memberId][deviceId] = conn
	c.numbers = c.numbers + 1
	slog.Info("success to make connection",
		"number of current connection", c.numbers)
	c.csMutex.Unlock()

	ip, _ := getPodIP()
	defer func() {
		destroy := context.Background()
		err = conn.Close(websocket.StatusNormalClosure, "")
		if err != nil {
			slog.Error("fail to close conn", "err", err)
		}
		c.csMutex.Lock()
		delete(c.conns[memberId], deviceId)
		c.numbers = c.numbers - 1
		c.csMutex.Unlock()
		c.csMutex.RLock()
		slog.Info("success to close connection",
			"number of current connection", c.numbers)
		if len(c.conns[memberId]) == 0 {
			c.service.RemoveServerIP(destroy, memberId[:], ip)
		}
		c.csMutex.RUnlock()
	}()

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
		slog.Error("incorrect payload types",
			"msgType", msgType,
			"data", data)
		handleWebsocketError(ctx, conn, errors.New("incorrect payload types"))
		return
	}
	if err != nil {
		slog.Error("read error",
			"err", err)
		handleWebsocketError(ctx, conn, errors.New("read error"))
		return
	}
}

func (c *Controller) RelayMessaging(ctx context.Context, id uuid.UUID, toIds [][]byte, roomId, fromId uuid.UUID, contentType string, contents []string) ([][]byte, error) {
	var wg sync.WaitGroup
	var pushToIds [][]byte
	var mu sync.Mutex
	res := dto.MessagingResponse{
		Id:          id,
		RoomId:      roomId,
		FromId:      fromId,
		ContentType: contentType,
		Contents:    contents,
	}
	resRaw := payload.Marshal(res)
	for _, toId := range toIds {
		r := resRaw
		if bytes.Equal(toId, roomId[:]) {
			res1 := res
			res1.RoomId = fromId
			resRaw1 := payload.Marshal(res1)
			r = resRaw1
		}
		wg.Add(1)
		go func() {
			defer wg.Done()
			c.csMutex.RLock()
			ct, ok := c.conns[uuid.UUID(toId)]
			var sm sync.Mutex
			var success bool
			if ok {
				var wg1 sync.WaitGroup
				for _, cd := range ct {
					wg1.Add(1)
					go func() {
						defer wg1.Done()
						err1 := cd.Write(ctx, websocket.MessageText, r)
						if err1 != nil {
							slog.Error("fail to write payload",
								"err", err1,
							)
							err2 := cd.CloseNow()
							if err2 != nil {
								slog.Error("fail to close zombie connection", "err", err2)
							}
							return
						}
						sm.Lock()
						success = true
						sm.Unlock()
					}()
				}
				wg1.Wait()
			}
			c.csMutex.RUnlock()
			if !ok || !success {
				mu.Lock()
				pushToIds = append(pushToIds, toId[:])
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	return pushToIds, nil
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
