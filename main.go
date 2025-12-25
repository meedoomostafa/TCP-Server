package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

type Message struct {
	from    string
	payload []byte
}

type Server struct {
	listenAddr string
	ln         net.Listener
	quitch     chan struct{}
	msgch      chan Message

	peerLock sync.Mutex
	peers    map[net.Conn]bool
}

func NewServer(listenAddr *string) *Server {
	return &Server{
		listenAddr: *listenAddr,
		quitch:     make(chan struct{}),
		msgch:      make(chan Message, 10),
		peers:      make(map[net.Conn]bool),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()
	go s.inputLoop()

	<-s.quitch

	close(s.msgch)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			fmt.Println("log: ", err)
			continue
		}

		s.peerLock.Lock()
		s.peers[conn] = true
		s.peerLock.Unlock()

		fmt.Println("new connection to the server : ", conn.RemoteAddr())
		go s.ReadLoop(conn)
	}
}

func (s *Server) ReadLoop(conn net.Conn) {
	defer func() {
		conn.Close()
		s.peerLock.Lock()
		delete(s.peers, conn)
		s.peerLock.Unlock()
	}()
	buf := make([]byte, 2048)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read error : ", err)
			return
		}

		s.msgch <- Message{
			from:    conn.RemoteAddr().String(),
			payload: buf[:n],
		}
	}
}

func (s *Server) inputLoop() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		fullMsg := fmt.Sprintf("Server: %s\n", msg)
		s.peerLock.Lock()

		for conn := range s.peers {
			conn.Write([]byte(fullMsg))
		}
		s.peerLock.Unlock()
	}
}

func main() {
	addr := ":3000"
	server := NewServer(&addr)
	go func() {
		for msg := range server.msgch {
			fmt.Printf("Received message: %s: %s", msg.from, string(msg.payload))
		}
	}()
	log.Fatal(server.Start())
}
