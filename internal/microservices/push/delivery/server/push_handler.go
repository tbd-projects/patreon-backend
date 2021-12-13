package push_server

import (
	"bytes"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	bh "patreon/internal/app/delivery/http/handlers/base_handler"
	session_client "patreon/internal/microservices/auth/delivery/grpc/client"
	"patreon/internal/microservices/auth/sessions/middleware"
	"time"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type PushHandler struct {
	sessionClient session_client.AuthCheckerClient
	bh.BaseHandler
}

func NewPushHandler(log *logrus.Logger, sManager session_client.AuthCheckerClient) *PushHandler {
	h := &PushHandler{
		BaseHandler:   *bh.NewBaseHandler(log),
		sessionClient: sManager,
	}
	h.AddMethod(http.MethodGet, h.GET, middleware.NewSessionMiddleware(h.sessionClient, log).CheckFunc)
	return h
}

// GET Push Creator
// @Summary create creator
// @Description create creator with id from path, and respond created creator
// @Param creator body http_models.RequestCreator true "Request body for creators"
// @Produce json
// @tags creators
// @Success 201 {object} http_models.IdResponse
// @Failure 409 {object} http_models.ErrResponse "creator already exist"
// @Failure 404 {object} http_models.ErrResponse "user with this id not found"
// @Failure 500 {object} http_models.ErrResponse "can not do bd operation", "server error"
// @Failure 422 {object} http_models.ErrResponse "invalid creator nickname", "invalid creator category-description", "invalid creator category", "invalid body in request"
// @Failure 403 {object} http_models.ErrResponse "csrf token is invalid, get new token"
// @Failure 401 "user are not authorized"
// @Router /user/push [GET]
func (h *PushHandler) GET(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	userId, ok := r.Context().Value("user_id").(int64)
	if !ok {
		h.Log(r).Errorf("not found user_id in context")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	go func(conn *websocket.Conn) {
		ticker := time.NewTicker(pingPeriod)
		defer func() {
			ticker.Stop()
			_ = conn.Close()
		}()
		conn.SetReadLimit(maxMessageSize)
		_ = conn.SetReadDeadline(time.Now().Add(pongWait))
		conn.SetPongHandler(func(string) error { _ = conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

		massages := make(chan []byte)
		errorCh := make(chan error)
		done := make(chan bool)

		go func() {
			for {
				select {
				case <-done:
					return
				default:
					break
				}
				_, tmpmessage, tmperr := conn.ReadMessage()
				massages <- tmpmessage
				errorCh <- tmperr
			}
		}()

		defer func() { done <- true }()

		for {
			var message []byte
			select {
			case <-ticker.C:
				_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
				if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
					return
				}
				continue
			case message = <-massages:
				err = <-errorCh
				break
			}
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					h.Log(r).Warnf("error: %v", err)
				}
				break
			}
			message = bytes.TrimSpace(bytes.Replace(message, []byte{'\n'}, []byte{' '}, -1))

			_ = conn.SetWriteDeadline(time.Now().Add(writeWait))
			w, err := conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			message = append(message, []byte{' '}...)
			message = append(message, []byte(fmt.Sprintf("%d", userId))...)
			h.Log(r).Infof("Send massage %v", massages)
			_, _ = w.Write(message)
			if err := w.Close(); err != nil {
				return
			}

		}
	}(conn)
}
