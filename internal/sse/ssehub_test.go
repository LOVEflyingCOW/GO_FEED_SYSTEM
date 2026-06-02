package sse

import (
	"testing"
	"time"
)

func TestHub_Broadcast(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Stop()

	client := &Client{
		AccountID: 1,
		Channel:   make(chan []byte, 10),
		IsClosed:  false,
	}

	hub.mu.Lock()
	hub.clients[1] = append(hub.clients[1], client)
	hub.mu.Unlock()

	msg := Notification{
		Type:    "test",
		FromID:  2,
		From:    "testuser",
		Content: "test content",
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		hub.Broadcast(msg)
	}()

	select {
	case received := <-client.Channel:
		if len(received) == 0 {
			t.Error("Broadcast() message should not be empty")
		}
	case <-time.After(2 * time.Second):
		t.Error("Broadcast() timeout waiting for message")
	}
}

func TestHub_SendTo(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Stop()

	client := &Client{
		AccountID: 1,
		Channel:   make(chan []byte, 10),
		IsClosed:  false,
	}

	hub.mu.Lock()
	hub.clients[1] = append(hub.clients[1], client)
	hub.mu.Unlock()

	msg := Notification{
		Type:    "test",
		FromID:  2,
		From:    "testuser",
		Content: "test content",
	}

	go func() {
		time.Sleep(100 * time.Millisecond)
		hub.SendTo(1, msg)
	}()

	select {
	case received := <-client.Channel:
		if len(received) == 0 {
			t.Error("SendTo() message should not be empty")
		}
	case <-time.After(2 * time.Second):
		t.Error("SendTo() timeout waiting for message")
	}
}

func TestHub_SendTo_NoClient(t *testing.T) {
	hub := NewHub()
	go hub.Run()
	defer hub.Stop()

	msg := Notification{
		Type:    "test",
		FromID:  2,
		From:    "testuser",
		Content: "test content",
	}

	result := hub.SendTo(999, msg)
	if result {
		t.Error("SendTo() should return false when no client exists")
	}
}

func TestClient_Close(t *testing.T) {
	client := &Client{
		AccountID: 1,
		Channel:   make(chan []byte, 10),
		IsClosed:  false,
	}

	client.Close()

	if !client.IsClientClosed() {
		t.Error("Client should be closed")
	}
}

func TestHub_Stop(t *testing.T) {
	hub := NewHub()
	go hub.Run()

	time.Sleep(100 * time.Millisecond)

	hub.Stop()

	hub.mu.RLock()
	running := hub.isRunning
	hub.mu.RUnlock()

	if running {
		t.Error("Hub should be stopped")
	}
}
