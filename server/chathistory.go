package main

import (
	"fmt"
	"os"
)

type ChatHistory struct {
	write chan string
	done  chan bool
	f     *os.File
}

func NewChatHistory(fileName string) (*ChatHistory, err) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}

	return &ChatHistory{write: make(chan string), done: make(chan bool), f: f}, nil
}

func (c *ChatHistory) Close() {
	c.done <- true
	close(c.write)
	close(c.f)
}

func (c *ChatHistory) Write(msg string) {
	c.write <- msg
}

func (c *ChatHistory) Start() {
	for {
		select {
		case msg := <-c.write:
			fmt.Fprint(c.f, msg)
		case <-c.done:
			return
		}
	}
}
