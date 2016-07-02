package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"

	"github.com/MoonWatcher582/chatroom/common"
)

var users *common.AtomicMap

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

func handleClient(conn *common.AtomicConn) {
	go conn.Start()
	go conn.ReadLoop()
	// @TODO Put chathistory through conn.Write
	for {
		msg := conn.Read()
		switch msg {
		case "ls":
			conn.Write("Users in chat:\n")
			listAllUsers(conn)
		case "quit":
			users.Remove(conn.Name)
			conn.Close()
			break
		default:
			// @TODO If user sends new message, write to chat history
			fmt.Fprintln(os.Stdout, msg)
		}
	}
	// @TODO If new message is written to chat history that isn't by the user, write to user
}

func main() {
	fmt.Fprintln(os.Stdout, "Chatroom server starting up...")

	users = common.NewAtomicMap()
	userCount := 0

	ln, err := net.Listen("tcp", host+":"+port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to start server:", err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if userCount >= maxUsers {
			fmt.Fprintln(os.Stderr, "Sorry, chatroom is currently full.")
			fmt.Fprintln(conn, "full")
			continue
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unsuccessfully tried to accept connection:", err)
			fmt.Fprintln(conn, "failed")
			continue
		}

		username, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not receive username:", err)
			fmt.Fprintln(conn, "failed")
			continue
		}

		fmt.Fprintln(conn, "success")

		atconn := common.NewAtomicConn(conn, username)
		users.Set(username, atconn)
		userCount++

		fmt.Fprint(os.Stdout, "Successfully accepted connection. Serving new user ", username)
		go handleClient(atconn)
	}
}
