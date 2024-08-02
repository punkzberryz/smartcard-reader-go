package server

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/punkzberryz/smartcard-reader-go/pkg/model"
)

type SocketServer struct {
	Server
	//Socket related
	Broadcast chan model.Message
}

//go:embed index.html
var indexPage []byte

func (s *SocketServer) RunServerWithWebSocket() error {
	socketServer := NewSocketIO()
	go func() {
		if err := socketServer.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer socketServer.Close()

	webSocket := NewWS()
	go func() {
		for {
			msg, ok := <-s.Broadcast
			if ok {
				socketServer.Broadcast(msg)
				webSocket.Broadcast(msg)
			}
		}
	}()

	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(200)
			w.Write(indexPage)
		}))
	mux.Handle("/health", http.HandlerFunc(s.handleHome))
	mux.Handle("/read", s.auth(http.HandlerFunc(s.handleRead)))
	mux.Handle("/ws", s.auth(http.HandlerFunc(webSocket.handleWS)))
	http.Handle("/socket.io/", socketServer)
	return http.ListenAndServe(fmt.Sprintf(":%s", s.Server.ServerConfig.Port), mux)
}

type connection struct {
	// The websocket connection.
	ws *websocket.Conn
}

func (c *connection) write(mt int, payload []byte) error {
	c.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
	return c.ws.WriteMessage(mt, payload)
}

type subscriber struct {
	// put registered clients.
	clients map[*connection]bool
}

func (s *subscriber) register(c *connection) {
	if s.clients == nil {
		s.clients = make(map[*connection]bool)
	}
	s.clients[c] = true
}
func (s *subscriber) unregister(c *connection) {
	if s.clients != nil {
		delete(s.clients, c)
	}
}

type ws struct {
	subscriber
}

func NewWS() *ws {
	s := subscriber{}
	return &ws{s}
}

func (s *ws) handleWS(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true },
	}
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		//failed to upgrade
		fmt.Println(err)
		http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		return
	}

	//register the client
	c := &connection{ws: ws}
	s.subscriber.register(c)

	defer ws.Close()

	// helpful log statement to show connections
	log.Println("Client Connected")
	//greet the client
	if err := ws.WriteMessage(websocket.TextMessage, []byte("Hello Client!")); err != nil {
		fmt.Println(err)
		return
	}

	//infinite loop to read message from client
	for {
		//read in a message
		messageType, message, err := ws.ReadMessage()
		if err != nil {
			fmt.Println(err)
			s.subscriber.unregister(c)
			break
		}
		//print out the message for clarity
		if messageType == websocket.TextMessage {
			fmt.Printf("Received: %s\n", message)
			//handle the message or broadcast it
			// return
		}
		if messageType == websocket.CloseMessage {
			s.subscriber.unregister(c)
			break
		}

	}

}

func (s *ws) Broadcast(msg model.Message) {
	m, _ := json.Marshal(msg)
	for c := range s.subscriber.clients {
		c.write(websocket.TextMessage, m)
	}
}
