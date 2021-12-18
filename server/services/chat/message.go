package chat

import (
	"encoding/json"
	"log"
)

const (
	Join Action = iota
	Leave
	Send
	Notify
)

type Action int

type Message struct {
	Action  Action
	Sender  *Client
	Content string
	Target  string
}

func (m *Message) encode() []byte {
	msg, err := json.Marshal(m)
	if err != nil {
		log.Println("Error encoding message: ", err)
		return []byte{}
	}
	return msg
}
