package server

import (
	"net"
)

type room struct {
	name    string
	members map[net.Addr]*client
}

func (r *room) broadcast(sender *client, msg string) {
	for addr, client := range r.members {
		if addr != sender.conn.RemoteAddr() {
			client.msg(msg)
		}
	}
}
