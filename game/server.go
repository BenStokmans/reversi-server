package game

import (
	"github.com/BenStokmans/reversi-server/snowflake"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/anypb"
	"math/rand"
	"net"
	"sync"
	"time"
)

type MessageHandler func(msg *anypb.Any, client *Client) error

type Server struct {
	clients    map[snowflake.Snowflake]*Client
	clientsMut sync.Mutex
	games      map[snowflake.Snowflake]*Game
	gamesMut   sync.Mutex

	messageHandleFunc MessageHandler
}

func NewServer(handler MessageHandler) Server {
	return Server{
		clients:           make(map[snowflake.Snowflake]*Client),
		games:             make(map[snowflake.Snowflake]*Game),
		messageHandleFunc: handler,
	}
}

func (g *Server) Start() {
	rand.Seed(time.Now().UnixMilli())
	l, err := net.Listen("tcp", ":8080")
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.Infof("listening on %v", l.Addr())
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err != nil {
			logrus.Fatal(err)
		}
		go g.handleNewConnection(conn)
	}
}

func (g *Server) handleNewConnection(conn net.Conn) {
	count := 0
	for _, client := range g.clients {
		if client.conn.RemoteAddr() == conn.RemoteAddr() {
			count++
		}
	}
	if count > 5 {
		logrus.Error("5 or more connection from same ip")
		_ = conn.Close()
		return
	}
	logrus.Infof("new connection from %v", conn.RemoteAddr())

	client := NewClient(conn, g, g.messageHandleFunc)
	client.Start()

	g.clientsMut.Lock()
	g.clients[client.id] = &client
	g.clientsMut.Unlock()
}
