package push_server

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"time"
)

type Client struct {
	logger   *logrus.Entry
	hub      *SendHub
	clientId int64
	conn     *websocket.Conn
	send     chan interface{}
	close    chan bool
}

func NewClient(hub *SendHub, clientId int64, conn *websocket.Conn, logger *logrus.Entry) *Client {
	return &Client{
		hub:      hub,
		clientId: clientId,
		conn:     conn,
		send:     make(chan interface{}),
		close:    make(chan bool),
		logger:   logger,
	}
}

func (c *Client) CloseClient() {
	c.hub.UnregisterClient(c)
	c.close <- true
}

func (c *Client) SenderProcesses() {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		_ = c.conn.Close()
		c.CloseClient()
	}()
	c.hub.SendMessage([]int64{c.clientId}, fmt.Sprintf("Hello user with Id: %d", c.clientId))

	for {
		select {
		case <-c.close:
			_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			c.logger.Infof("client with id %d was close", c.clientId)
			return
		case msg, ok := <-c.send:
			if !ok {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				c.logger.Infof("client with id %d was close by sendHub", c.clientId)
				return
			}
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))

			err := c.conn.WriteJSON(msg)
			if err != nil {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte("server error"))
				c.logger.Errorf("client with id %d write msg with error %s", c.clientId, err)
				return
			}
		case <-ticker.C:
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

type SendHub struct {
	clients    map[int64]*Client
	broadcast  chan *message
	register   chan *Client
	unregister chan *Client
	stopHub    chan bool
}

type message struct {
	users   []int64
	message interface{}
}

func NewHub() *SendHub {
	return &SendHub{
		broadcast:  make(chan *message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[int64]*Client),
		stopHub:    make(chan bool),
	}
}

func (h *SendHub) RegisterClient(client *Client) {
	h.register <- client
}

func (h *SendHub) UnregisterClient(client *Client) {
	h.unregister <- client
}

func (h *SendHub) SendMessage(users []int64, hsg interface{}) {
	h.broadcast <- &(message{users: users, message: hsg})
}

func (h* SendHub) StopHub() {
	h.stopHub <- true
}

func (h* SendHub) unregisterAll() {
	for key, client := range h.clients {
		close(client.send)
		client.CloseClient()
		delete(h.clients, key)
	}
}

func (h *SendHub) sendMessage(msg *message) {
	for _, id := range msg.users {
		if client, ok := h.clients[id]; ok {
			select {
			case client.send <- msg.message:
				break
			default:
				h.UnregisterClient(client)
			}
		}
	}
}

func (h *SendHub) Run() {
	for {
		select {
		case client, ok := <-h.register:
			if ok {
				h.clients[client.clientId] = client
			}
			break
		case client, ok := <-h.unregister:
			if ok {
				if _, ok = h.clients[client.clientId]; ok {
					delete(h.clients, client.clientId)
					close(client.send)
				}
			}
			break
		case msg, ok := <-h.broadcast:
			if ok {
				h.sendMessage(msg)
			}
			break
		case <-h.stopHub:
			h.unregisterAll()
			return
		}
	}
}
