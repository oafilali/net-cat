package main

import (
	"log"
	"net"
)

// pendingConversation is mapping the group name to the conversation string
type client struct {
	name            string
	currActiveGroup string
	conn            net.Conn
	pendingConv     map[string]string
}

type clients []client

var clientsArr clients

func isClientIn(conn net.Conn) bool {
	for _, c := range clientsArr {
		if c.conn == conn {
			return true
		}
	}
	return false
}

func registerClient(name, currGroup string, conn net.Conn) int {
	if isClientIn(conn) {
		log.Println("Can't register client. It's already in.")
		return getClientId(conn)
	}
	clientsArr = append(clientsArr, client{
		name:            name,
		currActiveGroup: currGroup,
		conn:            conn,
		pendingConv:     make(map[string]string),
	})
	return len(clientsArr) - 1
}

func getClientById(clientId int) *client {
	if clientId < len(clientsArr) {
		return &clientsArr[clientId]
	}
	return nil
}

func getClientByConn(conn net.Conn) *client {
	for i, c := range clientsArr {
		if c.conn == conn {
			return &clientsArr[i]
		}
	}
	return nil
}

func getClientId(conn net.Conn) int {
	for i, c := range clientsArr {
		if c.conn == conn {
			return i
		}
	}
	return -1
}
