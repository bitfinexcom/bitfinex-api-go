package tests

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type client struct {
	parent *TestWsService
	*websocket.Conn
	send     chan []byte
	received []string
	lock     sync.Mutex
}

func (c *client) writePump() {
	for msg := range c.send {
		err := c.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Printf("could not send message (%s) to client: %s", string(msg), err.Error())
			continue
		}
	}
}

func (c *client) readPump() {
	defer func() {
		c.parent.unregister <- c
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			log.Printf("test ws service drop client: %s", err.Error())
			return
		}
		message = bytes.TrimSpace(bytes.Replace(message, []byte("\n"), []byte(" "), -1))
		c.lock.Lock()
		log.Printf("[DEBUG] WsClient -> WsService: %s", string(message))
		c.received = append(c.received, string(message))
		c.lock.Unlock()
	}
}

type TestWsService struct {
	clients  map[*client]bool
	listener net.Listener
	port     int

	register     chan *client
	unregister   chan *client
	broadcast    chan []byte
	totalClients int
	lock         *sync.RWMutex

	publishOnConnect string
}

func (s *TestWsService) WaitForClientCount(count int) error {
	loops := 80
	delay := time.Millisecond * 50
	for i := 0; i < loops; i++ {
		s.lock.RLock()
		if s.totalClients == count {
			return nil
		}
		s.lock.RUnlock()
		time.Sleep(delay)
	}
	return fmt.Errorf("client peer #%d did not connect", count)
}

func (s *TestWsService) TotalClientCount() int {
	return s.totalClients
}

func (s *TestWsService) PublishOnConnect(msg string) {
	s.publishOnConnect = msg
}

func NewTestWsService(port int) *TestWsService {
	return &TestWsService{
		port:       port,
		clients:    make(map[*client]bool),
		register:   make(chan *client),
		unregister: make(chan *client),
		broadcast:  make(chan []byte),
		lock:       &sync.RWMutex{},
	}
}

// Broadcast sends a message to all connected clients.
func (s *TestWsService) Broadcast(msg string) {
	s.broadcast <- []byte(msg)

}

// ReceivedCount starts indexing clients at position 0.
func (s *TestWsService) ReceivedCount(clientNum int) int {
	i := 0
	for client := range s.clients {
		if i == clientNum {
			client.lock.Lock()
			defer client.lock.Unlock()
			return len(client.received)
		}
		i++
	}
	return 0
}

// Received starts indexing clients and message positions at position 0.
func (s *TestWsService) Received(clientNum int, msgNum int) (string, error) {
	var client *client
	i := 0
	for client = range s.clients {
		if i == clientNum {
			break
		}
		i++
	}
	if client != nil {
		client.lock.Lock()
		defer client.lock.Unlock()
		if len(client.received) > msgNum {
			return string(client.received[msgNum]), nil
		}
		return "", fmt.Errorf("could not find message index %d, %d messages exist", msgNum, len(client.received))
	}
	return "", fmt.Errorf("could not find client %d", clientNum)
}

func (s *TestWsService) WaitForMessage(clientNum int, msgNum int) (string, error) {
	loops := 80
	delay := time.Millisecond * 50
	var msg string
	var err error
	for i := 0; i < loops; i++ {
		msg, err = s.Received(clientNum, msgNum)
		if err != nil {
			time.Sleep(delay)
		} else {
			return msg, nil
		}
	}
	return "", err
}

func (s *TestWsService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.serveWs(w, r)
}

func (s *TestWsService) Stop() {
	//s.lock.RLock()
	//defer s.lock.RUnlock()
	s.listener.Close() // stop listening to http
	for c := range s.clients {
		c.Close()
	}
}

//nolint
func (s *TestWsService) Start() error {
	go s.loop()
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}
	s.listener = l
	go http.Serve(s.listener, s)
	return nil
}

//nolint
func (s *TestWsService) serveWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print(err)
		return
	}
	s.totalClients++
	client := &client{parent: s, Conn: conn, send: make(chan []byte, 256), received: make([]string, 0)}
	go client.writePump()
	go client.readPump()
	s.clients[client] = true
	if s.publishOnConnect != "" {
		s.Broadcast(s.publishOnConnect)
	}
}

func (s *TestWsService) loop() {
	for {
		select {
		case client := <-s.register:
			//s.lock.Lock()
			s.clients[client] = true
			//s.lock.Unlock()
		case client := <-s.unregister:
			if _, ok := s.clients[client]; ok {
				//s.lock.Lock()
				delete(s.clients, client)
				close(client.send)
				//s.lock.Unlock()
			}
		case msg := <-s.broadcast:
			for client := range s.clients {
				select {
				case client.send <- msg:
				default: // send failure
					//s.lock.Lock()
					close(client.send)
					delete(s.clients, client)
					//s.lock.Unlock()
				}
			}
		}
	}
}
