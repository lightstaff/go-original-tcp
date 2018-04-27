package protocol

import (
	"encoding/json"
	"io"
)

type Message struct {
	Headers map[string]string
	Body    interface{}
}

func (m *Message) Write(w io.Writer) error {
	return json.NewEncoder(w).Encode(m)
}

func ReadMessage(r io.Reader) (*Message, error) {
	var m Message
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

type Status int

const (
	OK Status = iota
	ERR
)

type Reply struct {
	Status Status
}

func (r *Reply) Write(w io.Writer) error {
	return json.NewEncoder(w).Encode(r)
}

func ReadReply(r io.Reader) (*Reply, error) {
	var rep Reply
	if err := json.NewDecoder(r).Decode(&rep); err != nil {
		return nil, err
	}

	return &rep, nil
}
