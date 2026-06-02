package realtime

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/google/uuid"
)

type HubConnectionWS struct {
	clients map[uuid.UUID]map[*Client]bool
	mu      *sync.RWMutex
}

func NewHubConnectionWS() *HubConnectionWS {
	return &HubConnectionWS{
		clients: make(map[uuid.UUID]map[*Client]bool),
		mu:      &sync.RWMutex{},
	}
}

func (h *HubConnectionWS) AddClient(client *Client, userId uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if h.clients[userId] == nil {
		h.clients[userId] = make(map[*Client]bool)
	}

	h.clients[userId][client] = true
	log.Println("Client connect to server with id: ", userId)
}

func (h *HubConnectionWS) RemoveClient(client *Client, userId uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if userClients, ok := h.clients[userId]; ok {
		delete(userClients, client)

		if len(userClients) == 0 {
			delete(h.clients, userId)
		}
	}

	log.Println("Remove Client connect to server with id: ", userId)
}

func (h *HubConnectionWS) sendToUser(userId uuid.UUID, message []byte) {
	h.mu.RLock()
	clients, exists := h.clients[userId]
	h.mu.RUnlock()

	if !exists {
		return
	}

	for client := range clients {
		client.Send(message)
	}
}

func (h *HubConnectionWS) SendEventToUser(userId uuid.UUID, payload any) error {
	data, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	h.sendToUser(userId, data)

	return nil
}
