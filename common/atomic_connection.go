package common

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
)

var colorBase int = 30
var bgBase int = 40

type AtomicConn struct {
	conn  net.Conn
	Name  string
	send  chan string
	rec   chan string
	done  chan bool
	Color string
}

func NewAtomicConn(conn net.Conn, name string) *AtomicConn {
	color := getColor()
	return &AtomicConn{conn: conn, Name: name, send: make(chan string), rec: make(chan string), done: make(chan bool), Color: color}
}

func (conn *AtomicConn) Close() {
	conn.done <- true
	conn.conn.Close()
	close(conn.rec)
	close(conn.send)
}

func (conn *AtomicConn) Write(msg string) {
	conn.send <- msg
}

func (conn *AtomicConn) Read() string {
	return <-conn.rec
}

func (conn *AtomicConn) ReadLoop() {
	scanner := bufio.NewScanner(conn.conn)
	for scanner.Scan() {
		conn.rec <- scanner.Text()
	}
}

func (conn *AtomicConn) Start() {
	for {
		select {
		case msg := <-conn.send:
			fmt.Fprintln(conn.conn, msg)
		case <-conn.done:
			return
		}
	}
}

func getColor() string {
	colorBase++
	bgBase += 2
	if colorBase > 37 {
		colorBase = 30
	}
	if bgBase > 47 {
		bgBase = 40 + (colorBase % 2)
	}

	fgColor := strconv.Itoa(colorBase)
	bgColor := strconv.Itoa(bgBase)
	bold := strconv.Itoa(colorBase % 2)
	fmt.Println("COLOR: x1b[0" + bold + ";" + fgColor + ";" + bgColor + "m")
	return "\x1b[0" + bold + ";" + fgColor + ";" + bgColor + "m"
}
