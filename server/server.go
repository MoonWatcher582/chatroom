package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/MoonWatcher582/chatroom/common"
)

var users *common.AtomicMap
var chatHistory *ChatHistory

const (
	maxUsers = 20
	host     = "127.0.0.1"
	port     = "8000"
)

func usage() {
	fmt.Fprintln(os.Stderr, "go run server.go")
	flag.PrintDefaults()
	os.Exit(1)
}

func listAllUsers(conn *common.AtomicConn) {
	allUsers := ""
	keys := users.Keys()
	for _, user := range keys {
		if user == "" {
			continue
		}
		allUsers += user + "\n"
	}
	if allUsers == "" {
		conn.Write("No users in chat\n")
		return
	}
	conn.Write(allUsers)
}

func sendToAll(conn *common.AtomicConn, msg string) {
	for _, k := range users.Keys() {
		u := users.Get(k)
		if u.Name != conn.Name {
			u.Write(msg)
		}
	}
}

func handleClient(conn *common.AtomicConn) {
	go conn.Start()
	go conn.ReadLoop()

	chatHistory.ReadAll(conn)
	// @TODO Put chathistory through conn.Write
	for {
		msg := conn.Read()
		switch msg {
		case "ls":
			conn.Write("Users in chat:\n")
			listAllUsers(conn)
		case "private":
			break
		case "end":
			break
		case "quit":
			users.Remove(conn.Name)
			conn.Close()
			break
		default:
			msg = fmt.Sprintf("%s%s says: %s\x1b[0m", conn.Color, conn.Name, msg)
			chatHistory.Write(conn, msg)
			sendToAll(conn, msg)
		}
	}
}

func main() {
	fmt.Fprintln(os.Stdout, "Chatroom server starting up...")

	var err error
	chatHistory, err = NewChatHistory("chathistory.txt")
	go chatHistory.Start()

	users = common.NewAtomicMap()
	userCount := 0

	ln, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to start server:", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unsuccessfully tried to accept connection:", err)
			continue
		}
		if userCount >= maxUsers {
			fmt.Fprintln(os.Stderr, "Sorry, chatroom is currently full.")
			fmt.Fprintln(conn, "full")
			continue
		}

		username, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not receive username:", err)
			fmt.Fprintln(conn, "failed")
			continue
		}
		username = strings.Trim(username, "\n")

		//if users.Get(username) != nil {
		//	fmt.Fprintln(os.Stdout, "User "+username+" already exists!")
		//	fmt.Fprintln(conn, "A user already exists with that name!")
		//	continue
		//}

		atconn := common.NewAtomicConn(conn, username)
		users.Set(username, atconn)
		fmt.Fprintln(conn, "success")
		userCount++

		fmt.Fprintln(os.Stdout, "Successfully accepted connection. Serving new user ", username)
		go handleClient(atconn)
	}
}
