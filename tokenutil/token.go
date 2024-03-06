package tokenutil

import (
	"encoding/base64"
	"encoding/json"
	"sort"
)

type Claims map[string]string

func (c Claims) Get(key string) string {
	v := c[key]
	return v
}

func (c Claims) Set(key, value string) {
	c[key] = value
}

func (c Claims) Sort() []string {
	keys := make([]string, len(c))
	i := 0
	for k := range c {
		keys[i] = k
		i++
	}
	sort.Strings(keys)
	return keys
}

func (c Claims) Encode() string {
	if len(c) == 0 {
		return ""
	}
	var text string
	for _, key := range c.Sort() {
		text += key + "=" + c[key] + "&"
	}
	return text[:len(text)-1]
}

func (c *Claims) Decode(text string) error {
	b, err := base64.StdEncoding.WithPadding(base64.NoPadding).DecodeString(text)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, c)
}

func (c Claims) String() string {
	var b []byte
	if len(c) == 0 {
		b = []byte("{}")
	} else {
		b, _ = json.Marshal(c)
	}
	return base64.StdEncoding.WithPadding(base64.NoPadding).EncodeToString(b)
}

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
