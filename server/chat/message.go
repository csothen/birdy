package chat

import (
	"encoding/json"
	"log"
)

const (
	JoinRoom Type = iota
	LeaveRoom
	CreateRoom
	DeleteRoom
	SendMessage
	Authenticate
)

type Type int

type Message struct {
	IsDM    bool   `json:"isDM"`
	Target  string `json:"target"`
	Sender  string `json:"sender"`
	Type    Type   `json:"type"`
	Content string `json:"content"`
}

func (m *Message) encode() []byte {
	msg, err := json.Marshal(m)
	if err != nil {
		log.Println("Error encoding message: ", err)
		return []byte{}
	}
	return msg
}
