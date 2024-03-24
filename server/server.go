package server

import (
	"fmt"
	"log"
	"net"
	"strings"
)

type server struct {
	rooms    map[string]*room
	commands chan command
}

func NewServer() *server {
	log.Println("Creating new server")

	return &server{
		rooms:    make(map[string]*room),
		commands: make(chan command),
	}
}

func (s *server) Run() {
	for cmd := range s.commands {
		switch cmd.id {
		case CMD_NAME:
			s.name(cmd.client, cmd.args)
		case CMD_JOIN:
			s.join(cmd.client, cmd.args)
		case CMD_ROOMS:
			s.listRooms(cmd.client)
		case CMD_MSG:
			s.msg(cmd.client, cmd.args)
		case CMD_QUIT:
			s.quit(cmd.client)
		case CMD_HELP:
			s.displayHelp(cmd.client)
		case CMD_VERSION:
			s.showVersion(cmd.client)
		}

	}
}

func (s *server) displayHelp(c *client) {
	c.msg("Available commands: /name, /join, /rooms, /msg, /quit, /help, /version")
}

func (s *server) showVersion(c *client) {
	c.msg("Version 1.0")
}

func (s *server) NewClient(conn net.Conn) *client {
	log.Printf("new client has connected: %s", conn.RemoteAddr().String())

	conn.Write([]byte("Welcome to ez-chat\n Enter /help for available commands\n"))



	return &client{
		conn:     conn,
		name:     "anonymous",
		commands: s.commands,
	}
}

func (s *server) name(c *client, args []string) {
	if len(args) < 2 {
		c.msg("name is required")
		return
	}

	c.name = args[1]

	c.msg(fmt.Sprintf("Hello %s", c.name))
}

func (s *server) join(c *client, args []string) {
	if len(args) < 2 {
		c.msg("room name is required")
		return
	}

	roomName := args[1]

	r, ok := s.rooms[roomName]
	if !ok {
		r = &room{
			name:    roomName,
			members: make(map[net.Addr]*client),
		}

		s.rooms[roomName] = r
	}

	r.members[c.conn.RemoteAddr()] = c
	c.room = r

	r.broadcast(c, fmt.Sprintf("%s has joined the room", c.name))

	c.msg(fmt.Sprintf("> Welcome to %s", r.name))
}

func (s *server) listRooms(c *client) {
	var rooms string
	for name := range s.rooms {
		rooms += name + ", "
	}
	c.msg(rooms)
}

func (s *server) msg(c *client, args []string) {
	if len(args) < 2 {
		c.msg("message is required")
		return
	}

	msg := strings.Join(args[1:], " ")
	c.room.broadcast(c, c.name+": "+msg)
}

func (s *server) quit(c *client) {
	s.quitCurrentRoom(c)

	c.msg("Bye Bye!")
	c.conn.Close()
}

func (s *server) quitCurrentRoom(c *client) {
	if c.room != nil {
		delete(c.room.members, c.conn.RemoteAddr())
		c.room.broadcast(c, fmt.Sprintf("%s has left the room", c.name))
	}
}
