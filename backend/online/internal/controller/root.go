package controller

import (
	"backend/online/internal/service"
	"log/slog"
	"net/http"

	"github.com/coder/websocket"
)

type HTTPMethod int

const (
	GET HTTPMethod = iota
	POST
	DELETE
	PUT
)

type Controller struct {
	service     *service.Service
	mux         *http.ServeMux
	connections map[string]*websocket.Conn
}

func NewController(s *service.Service, m *http.ServeMux) *Controller {

	c := &Controller{
		service: s,
		mux:     m,
	}

	conversationRouter(c)
	profileRouter(c)

	return c
}

func (c *Controller) Router(httpMethod HTTPMethod, path string, handler http.HandlerFunc) {
	m := c.mux

	switch httpMethod {
	case GET:
		m.HandleFunc("GET "+path, handler)
	case POST:
		m.HandleFunc("POST "+path, handler)
	case PUT:
		m.HandleFunc("PUT "+path, handler)
	case DELETE:
		m.HandleFunc("DELETE "+path, handler)

	default:
		panic("This HTTP method is not supported")
	}
}

func getStatusCode(err error) int {
	return http.StatusBadRequest
}

func handleParseError(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnprocessableEntity)
	_, err := w.Write([]byte("Not valid request"))
	if err != nil {
		slog.Error("fail to write response body",
			"err", err,
		)
	}
}

func handleServiceError(w http.ResponseWriter, err error) {
	w.WriteHeader(getStatusCode(err))
	_, err = w.Write([]byte(err.Error()))
	if err != nil {
		slog.Error("fail to write response body", "err", err)
	}
}
