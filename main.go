package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	clients     = make(map[net.Conn]string)
	clientMutex sync.Mutex
)

func main() {
	port := getPort()
	go startServer(port)
	startClient(port)
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
		addClient(conn)
		go handleConnection(conn)
	}
}

func addClient(conn net.Conn) {
	conn.Write([]byte("Welcome to TCP-Chat!\n"))
	linuxlogo, err := os.ReadFile("linuxlogo.txt")
	errorCheck("",err)
	conn.Write([]byte(linuxlogo))
	name := getName(conn)
	clientMutex.Lock()
	clients[conn] = name
	clientMutex.Unlock()
}

func removeClient(conn net.Conn) {
	clientMutex.Lock()
	delete(clients, conn)
	clientMutex.Unlock()
	conn.Close()
	fmt.Println("Client disconnected")
}

func getName(conn net.Conn) string {
	fmt.Print("[ENTER YOUR NAME]: ")
	reader := bufio.NewReader(conn)
	name, err := reader.ReadString('\n')
	errorCheck("Error reading the name:", err)
	return strings.TrimSpace(name)
}

func handleConnection(conn net.Conn) {
	defer removeClient(conn)

	fmt.Println("Client connected:", clients[conn])

	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Connection closed:", err)
			return
		}
		message = strings.TrimSpace(message)
		fmt.Printf("Message from %s: %s\n", clients[conn], message)
		broadcastMessage(conn, message)
	}
}

func broadcastMessage(sender net.Conn, message string) {
	clientMutex.Lock()
	defer clientMutex.Unlock()

	for client := range clients {
		if client == sender {
			continue
		}
		_, err := client.Write([]byte(message + "\n"))
		if err != nil {
			log.Printf("Error sending message to %s: %v\n", clients[client], err)
		}
	}
}

func startClient(port string) {
	conn, err := net.Dial("tcp", ":"+port)
	errorCheck("Error connecting to server: ", err)
	defer conn.Close()

	go listenForMessages(conn)

	fmt.Println("Connected to the server. Type your messages and press Enter:")
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		message, _ := reader.ReadString('\n')
		message = strings.TrimSpace(message)
		if message == "exit" {
			fmt.Println("Exiting client...")
			return
		}
		_, err := conn.Write([]byte(message + "\n"))
		if err != nil {
			log.Println("Error sending message:", err)
			return
		}
	}
}

func listenForMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Error receiving message:", err)
			return
		}
		fmt.Printf("[Server]: %s\n", strings.TrimSpace(message))
		fmt.Print("> ")
	}
}

func errorCheck(msg string, err error) {
	if err != nil {
		log.Println(msg, err)
		os.Exit(1)
	}
}
