package common

import (
	"bufio"
	"fmt"
	"net"
)

type AtomicConn struct {
	conn net.Conn
	Name string
	send chan string
	rec  chan string
	done chan bool
}

func NewAtomicConn(conn net.Conn, name string) *AtomicConn {
	return &AtomicConn{conn: conn, Name: name, send: make(chan string, 16), rec: make(chan string, 16), done: make(chan bool)}
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
			fmt.Fprint(conn.conn, msg)
		case <-conn.done:
			return
		}
	}
}
