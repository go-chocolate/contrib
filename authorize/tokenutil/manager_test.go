package tokenutil

import (
	"context"
	"testing"
)

func TestNewManager(t *testing.T) {
	manager := NewManager()
	token, err := manager.GenToken(context.Background(), "1", "1", Claims{"foo": "bar"})
	if err != nil {
		t.Error(err)
	}
	if token == "" {
		t.Errorf("token is empty")
	}
	t.Log(token)
	claims, err := manager.ValidateToken(context.Background(), token)
	if err != nil {
		t.Error(err)
		return
	}
	if claims == nil {
		t.Errorf("claims is empty")
	}
	t.Log(map[string]string(claims))
	if claims["foo"] != "bar" {
		t.Errorf("claims is not equal")
	}
}
