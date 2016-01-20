package main

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

// Сервер
type Server struct {
	messages  []*Message
	clients   map[int]*Client
	addCh     chan *Client
	delCh     chan *Client
	sendAllCh chan *Message
	doneCh    chan bool
	errCh     chan error
}

// Создаем новый сервер
func NewServer() *Server {
	return &Server{
		messages:  make([]*Message, 0),
		clients:   make(map[int]*Client),
		addCh:     make(chan *Client),
		delCh:     make(chan *Client),
		sendAllCh: make(chan *Message),
		doneCh:    make(chan bool),
		errCh:     make(chan error),
	}
}

// Добавляем клиента
func (s *Server) Add(c *Client) {
	s.addCh <- c
}

// Удаляем клиента
func (s *Server) Del(c *Client) {
	s.delCh <- c
}

// Рассылаем сообщение всем клиентам
func (s *Server) SendAll(msg *Message) {
	s.sendAllCh <- msg
}

// Пишем ошибку в лог
func (s *Server) Err(err error) {
	s.errCh <- err
}

// Отправляем все последние сообщения клиенту
func (s *Server) sendPastMessages(c *Client) {
	go func() {
		for _, msg := range s.messages {
			c.Write(msg)
		}
	}()
}

// Рассылаем сообщение всем клиентам
func (s *Server) sendAll(msg *Message) {
	for _, c := range s.clients {
		go c.Write(msg)
	}
}

// Обработчик для http.Handle
func (s *Server) Handler() http.Handler {
	return websocket.Handler(func(ws *websocket.Conn) {
		defer func() {
			if err := ws.Close(); err != nil {
				s.errCh <- err
			}
		}()

		client := NewClient(ws, s)
		s.Add(client)
		client.Listen()
	})
}

// Прослушка
func (s *Server) Listen() {
	for {
		select {
		// Потупил новый клиент
		case c := <-s.addCh:
			s.clients[c.id] = c
			log.Println("Added new client. Clients connected:", len(s.clients))
			s.sendPastMessages(c)

		// Поступил сигнал на удаление клиента
		case c := <-s.delCh:
			log.Println("Delete client number:", c.id)
			delete(s.clients, c.id)

		// Трансляция сообщения всем клиентам
		case msg := <-s.sendAllCh:
			s.messages = append(s.messages, msg)
			s.sendAll(msg)

		// Поступила ошибка
		case err := <-s.errCh:
			log.Println("Error:", err.Error())

		// Поступил сигнал на выход
		case <-s.doneCh:
			return
		}
	}
}
