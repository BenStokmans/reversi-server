package game

import (
	"github.com/sirupsen/logrus"
	"net"
)

type Server struct {
	clients []*Client
}

func (g *Server) Start() {
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

	client := NewClient(conn, g)
	client.Start()
	g.clients = append(g.clients, &client)
}
