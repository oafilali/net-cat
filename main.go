package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"
	"github.com/jroimartin/gocui"
	"net-cat/autocorrector"
)

var (
	clients     = make(map[net.Conn]string)
	clientMutex sync.Mutex
)

func main() {
	port := getPort()
	done := make(chan bool)
	go startServer(port)
	go startGui()
	<-done
}

func getPort() string {
	var port string
	switch len(os.Args) {
	case 1:
		port = "8989"
	case 2:
		port = os.Args[1]
	default:
		fmt.Println("[USAGE]: ./TCPChat $port")
		os.Exit(1)
	}
	return port
}

func startServer(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	errorCheck(fmt.Sprintf("Error starting server on port %s: ", port), err)
	defer listener.Close()

	fmt.Printf("Server listening on port %s...\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go handleNewClient(conn)
	}
}

func handleNewClient(conn net.Conn) {
	if len(clients) < 10 {
		conn.Write([]byte("Welcome to TCP-Chat!\n"))
		linuxlogo, err := os.ReadFile("linuxlogo.txt")
		errorCheck("Error reading linux logo:", err)
		conn.Write(linuxlogo)
		name := getName(conn)
		clientMutex.Lock()
		clients[conn] = name
		clientMutex.Unlock()
		joinMsg := fmt.Sprintf("%s has joined our chat...\n", name)
		broadcastMessage(conn, joinMsg)
		go handleConnection(conn)
	}
}

func removeClient(conn net.Conn) {
	clientMutex.Lock()
	name := clients[conn]
	delete(clients, conn)
	clientMutex.Unlock()
	leaveMsg := fmt.Sprintf("%s has left our chat...\n", name)
	broadcastMessage(conn, leaveMsg)
	conn.Close()
	log.Printf("Client %s disconnected", name)
}

func getName(conn net.Conn) string {
	conn.Write([]byte("[ENTER YOUR NAME]: "))
	reader := bufio.NewReader(conn)
	for {
		name, err := reader.ReadString('\n')
		if err != nil {
			errorCheck("Error reading the name:", err)
			conn.Write([]byte("[ENTER YOUR NAME]: "))
			continue
		}
		name = strings.TrimSpace(name)
		if name != "" {
			return name
		}
		conn.Write([]byte("[ENTER YOUR NAME]: "))
	}
}

func handleConnection(conn net.Conn) {
	defer removeClient(conn)

	fmt.Println("Client connected:", clients[conn])
	
	reader := bufio.NewReader(conn)
	for {
		currentTime := time.Now().Format("2006-01-02 15:04:05")
		conn.Write([]byte(fmt.Sprintf("[%s][%s]:", currentTime, clients[conn])))
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Connection closed:", err)
			return
		}
		message = strings.TrimSpace(message)
		if message == "" {
			continue
		}
		message = autocorrector.Input(message)
		formattedMessage := formatMessage(clients[conn], message)
		fmt.Printf("Message from %s: %s\n", clients[conn], message)
		broadcastMessage(conn, formattedMessage)
	}
}

func formatMessage(name, message string) string {
	if message == "" {
		return ""
	}
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s][%s]:%s\n", currentTime, name, message)
}

func broadcastMessage(sender net.Conn, message string) {
	clientMutex.Lock()
	defer clientMutex.Unlock()
	for client := range clients {
		if client == sender {
			continue
		}
		_, err := client.Write([]byte(message))
		if err != nil {
			log.Printf("Error sending message to %s: %v\n", clients[client], err)
		}
	}
}

func errorCheck(msg string, err error) {
	if err != nil {
		log.Println(msg, err)
		os.Exit(1)
	}
}

//trying to implement the terminal ui

func startGui() {
    g, err := gocui.NewGui(gocui.OutputNormal)
    if err != nil {
        log.Panicln(err)
    }
    defer g.Close()

    g.SetManagerFunc(layout)

    if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
        log.Panicln(err)
    }

    if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
        log.Panicln(err)
    }
}

func layout(g *gocui.Gui) error {
    maxX, maxY := g.Size()
    if v, err := g.SetView("log", 0, 0, maxX-1, maxY-3); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = "Log"
    }
    if v, err := g.SetView("input", 0, maxY-3, maxX-1, maxY-1); err != nil {
        if err != gocui.ErrUnknownView {
            return err
        }
        v.Title = "Input"
        v.Editable = true
    }
    return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
    return gocui.ErrQuit
}