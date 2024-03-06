package tokenutil

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Option func(m *Manager)

func applyOptions(m *Manager, options []Option) {
	for _, option := range options {
		option(m)
	}
}

func WithStorage(storage Storage) Option {
	return func(m *Manager) {
		m.storage = storage
	}
}

// WithMaxTokenPerUser set the max token count for per user, default is 10, and 0 means no limit.
// It means that user can generate different tokens for different clients.
func WithMaxTokenPerUser(max int) Option {
	return func(m *Manager) {
		m.maxTokenPerUser = max
	}
}

type Manager struct {
	storage         Storage
	maxTokenPerUser int
}

func NewManager(options ...Option) *Manager {
	m := &Manager{maxTokenPerUser: 10}
	applyOptions(m, options)
	if m.storage == nil {
		m.storage = NewMemoryStorage()
	}
	return m
}

func (m *Manager) GenToken(ctx context.Context, userId string, clientId string, claims Claims) (string, error) {
	tokens, err := m.getTokens(userId)
	if err != nil {
		return "", err
	}
	var token *Token
	if token = tokens.Get(clientId); token == nil {
		token = &Token{ClientId: clientId, Secret: randString(16)}
	}
	token.Timestamp = time.Now().UnixMilli()
	tokens.Set(token, m.maxTokenPerUser)
	if err = m.storage.Set(ctx, userId, tokens.JSON(), 7*24*time.Hour); err != nil {
		return "", err
	}

	head := Claims{}
	head["uid"] = userId
	head["cid"] = clientId
	head["nonce"] = randString(8)
	head["timestamp"] = strconv.FormatInt(token.Timestamp, 10)

	text := head.Encode() + claims.Encode() + token.Secret
	signature := toMd5([]byte(text))
	tokenString := fmt.Sprintf("%s.%s.%s", head.String(), claims.String(), signature)

	return tokenString, nil
}

func (m *Manager) ValidateToken(ctx context.Context, tokenString string) (Claims, error) {
	var head, claims = Claims{}, Claims{}
	var texts = strings.Split(tokenString, ".")
	if len(texts) != 3 {
		return nil, ErrTokenInvalid
	}
	if err := head.Decode(texts[0]); err != nil {
		return nil, err
	}
	if err := claims.Decode(texts[1]); err != nil {
		return nil, err
	}
	uid := head["uid"]
	cid := head["cid"]
	if uid == "" || cid == "" {
		return nil, ErrTokenInvalid
	}
	tokens, err := m.getTokens(uid)
	if err != nil {
		return nil, err
	}
	token := tokens.Get(cid)
	if token == nil {
		return nil, ErrTokenInvalid
	}
	if time.Now().UnixMilli() > token.Timestamp+7*24*86400*1000 {
		return nil, ErrTokenExpired
	}
	var signature = toMd5([]byte(head.Encode() + claims.Encode() + token.Secret))
	if signature != texts[2] {
		return nil, ErrTokenInvalid
	}
	return claims, nil
}

func (m *Manager) ValidateHTTPRequest(request *http.Request) (Claims, error) {
	var tokenString = request.Header.Get("Authorization")
	if len(tokenString) == 0 {
		return nil, ErrTokenInvalid
	}
	if strings.ToLower(tokenString[:7]) == "bearer " {
		tokenString = tokenString[7:]
	}
	return m.ValidateToken(request.Context(), tokenString)
}

func (m *Manager) getTokens(uid string) (Tokens, error) {
	data, err := m.storage.Get(context.Background(), uid)
	if err != nil {
		return Tokens{}, err
	}
	if len(data) == 0 {
		return Tokens{}, nil
	}
	var tokens = Tokens{}
	err = json.Unmarshal(data, &tokens)
	return tokens, err
}
