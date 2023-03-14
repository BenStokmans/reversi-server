package game

import (
	"github.com/BenStokmans/reversi-server/game/handlers"
	"github.com/BenStokmans/reversi-server/snowflake"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
	"net"
	"time"
)

type Client struct {
	id   snowflake.Snowflake
	game *Game

	owner  *Server
	conn   net.Conn
	sendCh chan *anypb.Any
	recvCh chan *anypb.Any

	closed  bool
	closeCh chan struct{}

	lastHeartbeat time.Time
}

func NewClient(conn net.Conn, owner *Server) Client {
	c := Client{
		id:            snowflake.Next(),
		owner:         owner,
		conn:          conn,
		lastHeartbeat: time.Now(),
	}
	c.sendCh = make(chan *anypb.Any)
	c.recvCh = make(chan *anypb.Any)
	return c
}

func (c *Client) Start() {
	go c.recvLoop()
	go c.sendLoop()
	go c.handleLoop()
	go c.heartbeatLoop()

	logrus.Infof("new client with id %v", c.id)
	c.Send(&Connected{PlayerId: int64(c.id)})
}

func (c *Client) Send(msg proto.Message) {
	if msg == nil {
		return
	}
	a := &anypb.Any{}
	_ = anypb.MarshalFrom(a, msg, proto.MarshalOptions{})
	c.sendCh <- a
}

func (c *Client) heartbeatLoop() {
	for {
		select {
		case <-time.After(1 * time.Second):
			if c.closed {
				return
			}

			if time.Now().Sub(c.lastHeartbeat) > 5*time.Second {
				logrus.Error("heartbeat timeout")
				c.Close()
				return
			}
		case <-c.closeCh:
			return
		}
	}
}

func (c *Client) handleLoop() {
	for {
		select {
		case data := <-c.recvCh:
			err := handlers.HandleMessage(data)
			if err != nil {
				logrus.Debug(err)
				continue
			}
			if data.MessageIs(&Heartbeat{}) {
				c.lastHeartbeat = time.Now()
			}
		case <-c.closeCh:
			return
		}
	}
}

func (c *Client) recvLoop() {
	for {
		data := make([]byte, 1024)
		n, err := c.conn.Read(data)
		if err != nil {
			if !c.closed && err.Error() != "EOF" {
				logrus.Error(err)
			}
			c.Close()
			return
		}
		msg := &anypb.Any{}

		err = proto.Unmarshal(data[:n], msg)
		if err != nil {
			return
		}
		c.recvCh <- msg
	}
}

func (c *Client) sendLoop() {
	for {
		select {
		case msg := <-c.sendCh:
			data, err := proto.Marshal(msg)
			if err != nil {
				logrus.Error(err)
				return
			}
			_, err = c.conn.Write(data)
			if err != nil {
				logrus.Error(err)
				return
			}
		case <-c.closeCh:
			return
		}
	}
}

func (c *Client) Close() {
	if c.closed {
		return
	}
	logrus.Infof("closing connection: %v", c.conn.RemoteAddr())
	_ = c.conn.Close()
	// remove client from owner connection
	for i, client := range c.owner.clients {
		if client == c {
			c.owner.clients = append(c.owner.clients[:i], c.owner.clients[i+1:]...)
		}
	}
	c.closed = true
	// tell other goroutines to close
	c.closeCh <- struct{}{}
}
