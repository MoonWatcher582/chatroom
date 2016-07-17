package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

const (
	default_host     = "127.0.0.1"
	default_port     = "8000"
	default_username = ""
	start_private    = "@private"
	end_private      = "@end"
	ls               = "@who"
	quit             = "@exit"
)

var host *string = flag.String("host", default_host, "The name of the host machine running the server.")
var port *string = flag.String("p", default_port, "The port number of the server process.")
var username *string = flag.String("u", default_username, "Your displayed username. A username must be provided.")

func usage() {
	fmt.Fprintln(os.Stderr, "useage: go run client.go -host=<host> -p=<port> -u=<username>")
	flag.PrintDefaults()
	os.Exit(1)
}

func writeInput(conn *net.Conn, done chan bool) {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		command_arg := ""
		fmt.Fprint(os.Stdout, ">")
		scanner.Scan()
		in := scanner.Text()

		if len(in) < 1 {
			fmt.Fprintln(os.Stderr, "Something went wrong, you may not have entered a valid command. Try again.")
			continue
		}

		if m := strings.SplitN(in, " ", 2); in[0] == '@' && len(m) > 1 {
			command_arg = m[1]
			in = m[0]
		}

		switch in {
		case start_private:
			fmt.Println("Beginning private session with", command_arg)
			// @TODO Write to conn requesting private session
		case end_private:
			fmt.Println("Ending private session with", command_arg)
			// @TODO Write to conn ending private session
		case ls:
			// Request list of all usernames
			fmt.Fprintln(*conn, "ls")
		case quit:
			// Disconnect from server
			done <- true
		default:
			fmt.Fprintln(*conn, in)
		}
	}
}

func readMsg(conn *net.Conn, done chan bool) {
	scanner := bufio.NewScanner(*conn)
	for scanner.Scan() {
		fmt.Fprintln(os.Stdout, scanner.Text())
	}

	fmt.Fprintln(os.Stderr, "Disconnected from server")
	done <- true
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *username == "" {
		flag.Usage()
	}

	if p, err := strconv.Atoi(*port); err != nil || p < 1024 || p > 65535 {
		fmt.Fprintln(os.Stderr, "Please provide a usable port number (1024 - 65535)")
		os.Exit(1)
	}

	done := make(chan bool, 1)

	conn, err := net.Dial("tcp", *host+":"+*port)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to "+*host+":"+*port+":", err)
		os.Exit(1)
	}

	fmt.Fprintln(conn, *username)

	resp, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to chatroom:", err)
		os.Exit(1)
	}

	switch resp {
	case "success\n":
		fmt.Fprintln(os.Stdout, "Connected to chatroom server.")
	case "full\n":
		fmt.Fprintln(os.Stderr, "Sorry, chatroom is currently full.")
		os.Exit(1)
	case "failed\n":
		fmt.Fprintln(os.Stderr, "Failed to connect to chatroom.")
		os.Exit(1)
	}

	go writeInput(&conn, done)
	go readMsg(&conn, done)

	<-done
	fmt.Fprintln(os.Stdout, "Exiting chatroom. Thank you!")
}
