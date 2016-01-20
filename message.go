package main

// Структура сообщений
type Message struct {
	Body string `json:"body"`
}

func (self *Message) String() string {
	return "> " + self.Body
}
