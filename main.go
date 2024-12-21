package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"net-cat/basic"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	BoldYellow  = "\033[33;1m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	BoldMagenta = "\033[35;1m"
	Cyan        = "\033[36m"
	Gray        = "\033[37m"
	White       = "\033[97m"
)

var (
	groupChats  = make(map[string][]int)
	clientMutex sync.Mutex
)

func main() {
	deleteChatFiles()
	clearChat()
	port := getPort()
	done := make(chan bool)
	go startServer(port)
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
	conn.Write([]byte(BoldYellow + "\nTo add/join a group chat:\n" + Reset + ":chat: <name of group chat>\n"))
	conn.Write([]byte(BoldYellow + "To change your name:\n" + Reset + ":name: <new name>\n"))
	conn.Write([]byte(BoldYellow + "To exit the current group chat:\n" + Reset + ":exit:\n"))
	conn.Write([]byte(BoldMagenta + "By default, you'll be added to the global chat unless it's full.\n\n" + Reset))
	go handleConnection(conn)
}

func removeClient(conn net.Conn, currentGroup string) {
	if currentGroup != "" {
		clientMutex.Lock()
		name := getClientByConn(conn).name
		for i, clientId := range groupChats[currentGroup] {
			c := getClientById(clientId)
			if c.conn == conn {
				groupChats[currentGroup][i], groupChats[currentGroup][len(groupChats[currentGroup])-1] = groupChats[currentGroup][len(groupChats[currentGroup])-1], groupChats[currentGroup][i]
				groupChats[currentGroup] = groupChats[currentGroup][:len(groupChats[currentGroup])-1]
				break
			}
		}
		clientMutex.Unlock()
		leaveMsg := fmt.Sprintf("%s has left our chat...\n", name)
		broadcastMessage(currentGroup, conn, Yellow+leaveMsg+Reset)
		if currentGroupName(conn) == "" {
			c := getClientByConn(conn)
			c.conn.Close()
			*c = client{}
			log.Printf("Client %s disconnected", name)
		}
	}
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
		if isDuplicateName(name) {
			conn.Write([]byte(Red + "NAME IS TAKEN\n" + Reset))
			conn.Write([]byte("[ENTER ANOTHER NAME]: "))
			continue
		}
		if name != "" {
			return name
		}
		conn.Write([]byte("[ENTER YOUR NAME]: "))
	}
}

// checks if chat is full, and creates it if it doesn't exist
func checkGroupChat(groupName string, conn net.Conn) error {
	_, ok := groupChats[groupName]
	if !ok {
		groupChats[groupName] = []int{}
	} else if c := getClientByConn(conn); c != nil && groupName == c.currActiveGroup {
		conn.Write([]byte("YOU'RE ALREADY IN " + groupName + "\n"))
		return errors.New("client already in group")
	} else if len(groupChats[groupName]) == 10 {
		conn.Write([]byte(BoldMagenta + "Oops, " + groupName + " chat is packed right now! Try again in a bit " + Reset + ":)\n"))
		return errors.New("group is full")
	}
	return nil
}

func isClientInGroup(groupName string, clientId int) bool {
	for _, id := range groupChats[groupName] {
		if id == clientId {
			return true
		}
	}
	return false
}

// returns bool value whether it has added client to group or not
func addClientToGroup(groupName string, clientId int) bool {
	// check if client id is already in the group

	// if the client is already in group
	if isClientInGroup(groupName, clientId) {
		return false
	}

	groupChats[groupName] = append(groupChats[groupName], clientId)
	return true
}

// addChat adds conn to the new group, and if it's the first time joining a group
// it registers the conn in the clientsArr
func joinChat(groupName string, conn net.Conn) {
	var clientName string
	if err := checkGroupChat(groupName, conn); err != nil {
		return
	}

	conn.Write([]byte("Welcome to " + groupName + " Chat!\n"))
	writeLogo(groupName, conn)

	c := getClientByConn(conn)
	if c == nil {
		clientName = getName(conn)
	} else {
		clientName = c.name
	}

	clientMutex.Lock()
	id := registerClient(clientName, groupName, conn)
	isAdded := addClientToGroup(groupName, id)
	clientsArr[id].currActiveGroup = groupName
	clientMutex.Unlock()

	joinMsg := fmt.Sprintf("%s has joined %s...\n", clientName, groupName)
	broadcastMessage(groupName, conn, Magenta+joinMsg+Reset)

	c = getClientByConn(conn)
	if isAdded {
		loadChat(c.conn, groupName)
	}
	if conv, ok := c.pendingConv[groupName]; ok {
		conn.Write([]byte(conv))
		delete(c.pendingConv, groupName)
	}
}

func writeLogo(groupName string, conn net.Conn) {
	if groupName != "global" {
		conn.Write([]byte(basic.Basic(cap(groupName), "standard")))
	} else {
		linuxlogo, err := os.ReadFile("linuxlogo.txt")
		errorCheck("Error reading linux logo:", err)
		conn.Write(linuxlogo)
	}
}

func welcomeBackTo(groupName string, conn net.Conn) {
	conn.Write([]byte("Welcome back to " + groupName + "\n"))
	writeLogo(groupName, conn)
	c := getClientByConn(conn)
	if conv, ok := c.pendingConv[groupName]; ok {
		conn.Write([]byte(conv))
		delete(c.pendingConv, groupName)
	}
}

func currentGroupName(conn net.Conn) string {
	for groupName, clients := range groupChats {
		for _, clientId := range clients {
			if getClientById(clientId).conn == conn {
				return groupName
			}
		}
	}
	return ""
}

func exitClient(conn net.Conn) string {
	cl := getClientByConn(conn)
	groupName := cl.currActiveGroup
	removeClient(conn, groupName)
	cl.currActiveGroup = currentGroupName(conn)
	if cl.currActiveGroup == "" {
		return "EXIT"
	} else {
		welcomeBackTo(cl.currActiveGroup, conn)
		return "CONTINUE"
	}
}

func processMessage(msg string, conn net.Conn) string {
	cl := getClientByConn(conn)

	if cl != nil && len(msg) > 7 && msg[0:7] == ":name: " {
		currAcGroup := cl.currActiveGroup
		newName := msg[7:]
		if newName == cl.name {
			conn.Write([]byte(Green + "You're already using that name, aren't you" + Reset + " :)\n"))
			return "CONTINUE"
		}
		if isDuplicateName(newName) {
			conn.Write([]byte(Red + "NAME IS TAKEN\n" + Reset))
			return "CONTINUE"
		}
		newNameMsg := "Heads up! [" + cl.name + "] is now going by [" + newName + "].\n"
		saveChat(newNameMsg, currAcGroup)
		broadcastMessage(currAcGroup, conn, Blue+newNameMsg+Reset)
		conn.Write([]byte(Green + "You've successfully changed your name\n" + Reset))
		cl.name = newName
		return "CONTINUE"
	}
	if len(msg) > 7 && msg[0:7] == ":chat: " {
		if len(msg) == 8 {
			conn.Write([]byte(Red+"Invalid chat name, 1 character isn't descriptive enough.\n"+Reset))
		} else {
			joinChat(strings.TrimSpace(msg[7:]), conn)
		}
		return "CONTINUE"
	}
	if msg == ":exit:" {
		return exitClient(conn)
	}
	if cl == nil {
		return "CONTINUE"
	}
	return "Broadcast Message"
}

func handleConnection(conn net.Conn) {
	joinChat("global", conn)
	reader := bufio.NewReader(conn)
	for {
		cl := getClientByConn(conn)
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Println("Connection closed:", err)
			return
		}
		message = strings.TrimSpace(message)
		if p := processMessage(message, conn); p == "CONTINUE" {
			continue
		} else if p == "EXIT" {
			break
		}
		formattedMessage := formatMessage(cl.name, message)
		msg := fmt.Sprintf("Message in %s from %s: %s\n", cl.currActiveGroup, getClientByConn(conn).name, sanitize(message))
		fmt.Print(msg)
		saveChat(formattedMessage, cl.currActiveGroup)
		broadcastMessage(cl.currActiveGroup, conn, formattedMessage)
	}
}

func formatMessage(name, message string) string {
	if message == "" {
		return ""
	}
	message = sanitize(message)
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s][%s]:%s\n", currentTime, name, message)
}

func broadcastMessage(brGroupName string, sender net.Conn, message string) {
	clientMutex.Lock()
	defer clientMutex.Unlock()
	for _, clientId := range groupChats[brGroupName] {
		c := getClientById(clientId)
		if c == nil {
			log.Fatal("SOMETHING IS WRONG\n", clientId)
		}
		if c.currActiveGroup == brGroupName {
			_, err := c.conn.Write([]byte(message))
			if err != nil {
				log.Printf("Error sending message to %s: %v\n", c.name, err)
			}
		} else {
			c.pendingConv[brGroupName] += message
		}

	}
}

func saveChat(message, chatName string) {
	file, err := os.OpenFile(chatName+".chat", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening chat log file:", err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(message)
	if err != nil {
		log.Println("Error writing message to chat log file:", err)
	}
}

func loadChat(client net.Conn, chatName string) {
	chat, err := os.ReadFile(chatName + ".chat")
	if err != nil {
		log.Println("Error loading the chat", err)
	}
	client.Write(chat)
}

func clearChat() {
	file, err := os.OpenFile("log.txt", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error clearing chat log file:", err)
		return
	}
	file.Close()
}

func deleteChatFiles() {
	entries, err := os.ReadDir("./")

	if err != nil {
		log.Println(err)
	}
	for _, e := range entries {
		if !e.IsDir() && len(e.Name()) > 5 && e.Name()[len(e.Name())-5:] == ".chat" {
			err := os.Remove(e.Name())
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func errorCheck(msg string, err error) {
	if err != nil {
		log.Println(msg, err)
		os.Exit(1)
	}
}

func cap(s string) string {
	return strings.ToUpper(s[:1]) + strings.ToLower(s[1:])
}

func sanitize(msg string) string {
	var s string

	for _, char := range msg {
		if (char >= 32 && char <= 126) || char == 'ö' || char == 'ä' || char == 'å' {
			s += string(char)
		}
	}
	return s
}
