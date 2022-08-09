package main

import (
	"fmt"
)

type message struct {
	data     []byte
	room     string
	username float64
}

type subscription struct {
	conn *connection
	room string
	user float64
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	rooms map[string]map[userconnection]bool

	// Inbound messages from the connections.
	broadcast chan message

	// Register requests from the connections.
	register chan subscription

	// Unregister requests from connections.
	unregister chan subscription
}

type userconnection struct {
	connection *connection
	username   float64
}

var h = hub{
	broadcast:  make(chan message),
	register:   make(chan subscription),
	unregister: make(chan subscription),
	rooms:      make(map[string]map[userconnection]bool),
}

var usercounter int

func (h *hub) run() {
	for {
		select {
		case s := <-h.register:
			connections := h.rooms[s.room]
			if connections == nil {

				connections := make(map[userconnection]bool)
				h.rooms[s.room] = connections
			}
			usrconn := userconnection{
				connection: s.conn,
				username:   s.user,
			}
			h.rooms[s.room][usrconn] = true
			usercounter++
			s.conn.send <- []byte(fmt.Sprintf("u r registered with id : %f and user %d", s.user, usercounter))

		case s := <-h.unregister:
			connections := h.rooms[s.room]
			if connections != nil {
				usrconn := userconnection{
					connection: s.conn,
					username:   s.user,
				}
				if _, ok := connections[usrconn]; ok {
					delete(connections, usrconn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.rooms, s.room)
					}
				}
			}
		case m := <-h.broadcast:
			fmt.Println("msg : ", m)
			connections := h.rooms[m.room]
			for c := range connections {
				select {
				case c.connection.send <- m.data:
				default:
					close(c.connection.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.rooms, m.room)
					}
				}
			}
		}
	}
}
