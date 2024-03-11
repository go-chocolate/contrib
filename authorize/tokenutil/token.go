package tokenutil

import (
	"encoding/json"
)

type Token struct {
	ClientId  string `json:"clientId"`
	Secret    string `json:"secret"`
	Timestamp int64  `json:"timestamp"`
}

type Tokens map[string]*Token

func (t Tokens) JSON() []byte {
	b, _ := json.Marshal(t)
	return b
}

func (t Tokens) Remove(clientId string) {
	delete(t, clientId)
}

func (t Tokens) Get(clientId string) *Token {
	return t[clientId]
}

func (t Tokens) Set(token *Token, max int) {
	t[token.ClientId] = token
	if max > 0 && len(t) > max {
		t.removeOldest()
	}
}

func (t Tokens) removeOldest() {
	if len(t) == 0 {
		return
	}
	var oldest *Token
	for _, v := range t {
		if oldest == nil || v.Timestamp < oldest.Timestamp {
			oldest = v
		}
	}
	delete(t, oldest.ClientId)
}
