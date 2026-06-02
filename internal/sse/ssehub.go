package sse

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type Client struct {
	AccountID uint
	Channel   chan []byte
	IsClosed  bool
	mu        sync.RWMutex
}

type Hub struct {
	clients   map[uint][]*Client
	broadcast chan []byte
	mu        sync.RWMutex
	isRunning bool
	stopChan  chan struct{}
}

func NewHub() *Hub {
	return &Hub{
		clients:   make(map[uint][]*Client),
		broadcast: make(chan []byte, 100),
		isRunning: true,
		stopChan:  make(chan struct{}),
	}
}

func (h *Hub) Run() {
	defer func() {
		if r := recover(); r != nil {
		}
	}()

	for {
		select {
		case <-h.stopChan:
			return
		case msg := <-h.broadcast:
			h.mu.RLock()
			for _, clients := range h.clients {
				for _, client := range clients {
					client.mu.RLock()
					isClosed := client.IsClosed
					client.mu.RUnlock()

					if !isClosed {
						select {
						case client.Channel <- msg:
						default:
						}
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) Stop() {
	h.mu.Lock()
	if !h.isRunning {
		h.mu.Unlock()
		return
	}
	h.isRunning = false
	h.mu.Unlock()

	select {
	case <-h.stopChan:
	default:
		close(h.stopChan)
	}

	h.mu.RLock()
	for _, clients := range h.clients {
		for _, client := range clients {
			client.Close()
		}
	}
	h.mu.RUnlock()
}

func (c *Client) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.IsClosed {
		c.IsClosed = true
		close(c.Channel)
	}
}

func (c *Client) IsClientClosed() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.IsClosed
}

func (h *Hub) Subscribe(accountID uint, c *gin.Context) error {
	client := &Client{
		AccountID: accountID,
		Channel:   make(chan []byte, 256),
		IsClosed:  false,
	}

	h.mu.Lock()
	h.clients[accountID] = append(h.clients[accountID], client)
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		clients := h.clients[accountID]
		for i, c := range clients {
			if c == client {
				h.clients[accountID] = append(clients[:i], clients[i+1:]...)
				break
			}
		}
		if len(h.clients[accountID]) == 0 {
			delete(h.clients, accountID)
		}
		h.mu.Unlock()
		client.Close()
	}()

	c.Stream(func(w io.Writer) bool {
		if client.IsClientClosed() {
			return false
		}

		select {
		case msg, ok := <-client.Channel:
			if !ok {
				return false
			}
			c.SSEvent("message", string(msg))
			return true
		case <-c.Request.Context().Done():
			return false
		case <-time.After(30 * time.Second):
			c.SSEvent("ping", "")
			return true
		}
	})

	return nil
}

func (h *Hub) SendTo(accountID uint, msg interface{}) bool {
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("SSE SendTo error: failed to marshal message: %v", err)
		return false
	}

	h.mu.RLock()
	clients, ok := h.clients[accountID]
	h.mu.RUnlock()

	if !ok {
		log.Printf("SSE SendTo error: no client connected for account %d", accountID)
		return false
	}

	log.Printf("SSE SendTo: sending message to account %d, %d clients connected", accountID, len(clients))

	sent := false
	for _, client := range clients {
		client.mu.RLock()
		isClosed := client.IsClosed
		client.mu.RUnlock()

		if !isClosed {
			select {
			case client.Channel <- data:
				sent = true
				log.Printf("SSE SendTo: message sent successfully to account %d", accountID)
			default:
				log.Printf("SSE SendTo: client channel is full for account %d", accountID)
			}
		}
	}
	return sent
}

func (h *Hub) Broadcast(msg interface{}) bool {
	data, err := json.Marshal(msg)
	if err != nil {
		return false
	}

	select {
	case h.broadcast <- data:
		return true
	default:
		return false
	}
}

func (h *Hub) ServeSSE(c *gin.Context) {
	accountIDVal, exists := c.Get("accountID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	accountID, ok := accountIDVal.(uint)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account ID"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Expose-Headers", "Content-Type")
	c.Header("X-Accel-Buffering", "no")

	h.Subscribe(accountID, c)
}

type Notification struct {
	Type    string      `json:"type"` // "like", "comment", "follow", "message"
	FromID  uint        `json:"from_id"`
	From    string      `json:"from"`
	Content interface{} `json:"content"`
}
