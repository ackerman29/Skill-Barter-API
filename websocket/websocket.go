package websocket

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true 
	},
}

var connections = make(map[string]*websocket.Conn)
var mu sync.Mutex
func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		http.Error(w, "email required", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Failed to upgrade to WebSocket", http.StatusInternalServerError)
		return
	}

	mu.Lock()
	connections[email] = conn
	mu.Unlock()

	for {
		if _, _, err := conn.NextReader(); err != nil {
			break
		}
	}

	mu.Lock()
	delete(connections, email)
	mu.Unlock()
	conn.Close()
}

func SendToUser(email string, message string) {
	mu.Lock()
	conn, ok := connections[email]
	mu.Unlock()

	if ok {
		conn.WriteMessage(websocket.TextMessage, []byte(message))
	}
}
// SendChatMessage sends a chat message to all users in the chat
func SendChatMessage(users []string, msg string) {
	mu.Lock()
	defer mu.Unlock()

	for _, u := range users {
		if conn, ok := connections[u]; ok {
			conn.WriteMessage(websocket.TextMessage, []byte(msg))
		}
	}
}
