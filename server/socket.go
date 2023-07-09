package server

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	gosocketio "github.com/ambelovsky/gosf-socketio"
	"github.com/ambelovsky/gosf-socketio/transport"
)

type Message struct {
	ChannelID string    `json:"channelId"`
	MessageID string    `json:"messageId"`
	Email     string    `json:"email"`
	Content   string    `json:"content"`
	CreatedAt string    `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Channel struct {
	ChannelID     string `json:"channelId"`
	PrevChannelID string `json:"prevChannelId"`
}

type SocketServer struct {
	server *gosocketio.Server
}

func newSocketServer() *SocketServer {
	server := gosocketio.NewServer(transport.GetDefaultWebsocketTransport())

	server.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Println("New client connected", c.Id())
	})

	server.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Println("Disconnected", c.Id())
	})

	server.On("join.channel", func(c *gosocketio.Channel, str string) string {
		var channel Channel
		json.Unmarshal([]byte(str), &channel)
		c.Leave("chat:" + channel.PrevChannelID)
		c.Join("chat:" + channel.ChannelID)
		return "Joined channel [" + "chat:" + channel.ChannelID + "]"
	})

	server.On("send.message", func(c *gosocketio.Channel, msg string) {
		var message Message
		json.Unmarshal([]byte(msg), &message)
		loc, _ := time.LoadLocation("Asia/Seoul")
		parse, _ := time.Parse(time.RFC3339, message.CreatedAt)
		parse = parse.In(loc)
		message.UpdatedAt = parse
		message.CreatedAt = message.UpdatedAt.String()
		c.BroadcastTo("chat:"+message.ChannelID, "send.message", message)
	})

	server.On("edit.message", func(c *gosocketio.Channel, msg string) string {
		var message Message
		json.Unmarshal([]byte(msg), &message)
		c.BroadcastTo("chat:"+message.ChannelID, "edit.message", message)
		return message.MessageID + " Edited"
	})

	server.On("delete.message", func(c *gosocketio.Channel, msg string) string {
		var message Message
		json.Unmarshal([]byte(msg), &message)
		c.BroadcastTo("chat:"+message.ChannelID, "delete.message", message)
		return message.MessageID + " Deleted"
	})

	server.On("edit.channel", func(s *gosocketio.Channel, msg string) string {
		var message Message
		json.Unmarshal([]byte(msg), &message)
		server.BroadcastToAll("edit.channel", message)
		return ""
	})

	server.On("delete.channel", func(s *gosocketio.Channel, msg string) string {
		var message Message
		message.ChannelID = msg
		server.BroadcastToAll("delete.channel", message)
		return ""
	})

	return &SocketServer{
		server: server,
	}
}

func StartSocketServer(serveMux *http.ServeMux) {
	socketServer := newSocketServer()
	socket := os.Getenv("SOCKET_SERVER_PORT")
	lis, _ := net.Listen("tcp", socket)
	serveMux.Handle("/socket.io/", socketServer.server)
	http.Serve(lis, serveMux)
}
