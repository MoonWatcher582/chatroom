package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/MoonWatcher582/chatroom/common"
)

type ChatHistory struct {
	name   string
	write  chan string
	done   chan bool
	writer *os.File
}

func NewChatHistory(fileName string) (*ChatHistory, error) {
	writer, err := os.Create(fileName)
	if err != nil {
		return nil, err
	}

	return &ChatHistory{name: fileName, write: make(chan string), done: make(chan bool), writer: writer}, nil
}

func (c *ChatHistory) Close() error {
	c.done <- true
	close(c.write)
	return c.writer.Close()
}

func (c *ChatHistory) Write(conn *common.AtomicConn, msg string) {
	c.write <- msg
}

func (c *ChatHistory) ReadAll(conn *common.AtomicConn) error {
	reader, err := os.Open(c.name)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		conn.Write(scanner.Text())
	}

	return scanner.Err()
}

func (c *ChatHistory) Start() {
	for {
		select {
		case msg := <-c.write:
			fmt.Fprintln(c.writer, msg)
		case <-c.done:
			return
		}
	}
}
