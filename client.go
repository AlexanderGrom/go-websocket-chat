package main

import (
	"golang.org/x/net/websocket"
	"io"
)

const (
	chanBufferSize = 100
)

var (
	clientId int = 0
)

// Клиент
type Client struct {
	id      int
	ws      *websocket.Conn
	server  *Server
	wCh     chan *Message
	rCh     chan *Message
	doneWCh chan bool
	doneRCh chan bool
	doneCh  chan bool
}

// Создаем нового клиента
func NewClient(ws *websocket.Conn, server *Server) *Client {
	clientId++
	return &Client{
		id:      clientId,
		ws:      ws,
		server:  server,
		wCh:     make(chan *Message, chanBufferSize),
		rCh:     make(chan *Message),
		doneWCh: make(chan bool),
		doneRCh: make(chan bool),
		doneCh:  make(chan bool),
	}
}

// Пишем сообщение клиенту
func (c *Client) Write(msg *Message) {
	c.wCh <- msg
}

// Читаем из вебсокета
func (c *Client) read() {
	var msg Message
	err := websocket.JSON.Receive(c.ws, &msg)
	if err != nil {
		if err == io.EOF {
			c.doneCh <- true
			return
		}
		c.server.Err(err)
	} else {
		c.rCh <- &msg
	}
}

// Слушаем чтение и запись, висим пока не поступит сигнал на выход
func (c *Client) Listen() {
	go c.listenWrite()
	go c.listenRead()
	select {
	case <-c.doneCh:
		c.server.Del(c)
		c.doneWCh <- true
		c.doneRCh <- true
		return
	}
}

// Слушаем запросы на запись
func (c *Client) listenWrite() {
	for {
		select {
		// Ждем послупления сообщения для записи в вебсокет
		case msg := <-c.wCh:
			websocket.JSON.Send(c.ws, msg)

		// Поступил сигнал на выход
		case <-c.doneWCh:
			return
		}
	}
}

// Слушаем запросы на чтение
func (c *Client) listenRead() {
	for {
		go c.read()
		select {
		// Ждем поступения сообщения из вебсокета, рассылаем
		case msg := <-c.rCh:
			c.server.SendAll(msg)

		// Поступил сигнал на выход
		case <-c.doneRCh:
			return
		}
	}
}
