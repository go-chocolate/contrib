package tokenutil

import (
	"context"
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

type claimsContextKey struct{}

var _claimsContextKey = &claimsContextKey{}

func (c Claims) WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, _claimsContextKey, c)
}

func FromContext(ctx context.Context) Claims {
	v := ctx.Value(_claimsContextKey)
	if v == nil {
		return make(Claims)
	}
	c, ok := v.(Claims)
	if !ok {
		return Claims{}
	}
	return c
}
