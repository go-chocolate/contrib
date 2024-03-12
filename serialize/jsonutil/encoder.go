package jsonutil

import (
	"bytes"
	"encoding/json"
	"io"
)

type encoder struct {
	e *json.Encoder
}

type EncodeOption func(e *encoder)

func WithIndent(prefix, indent string) EncodeOption {
	return func(e *encoder) {
		e.e.SetIndent(prefix, indent)
	}
}

func WithEscapeHTML(on bool) EncodeOption {
	return func(e *encoder) {
		e.e.SetEscapeHTML(on)
	}
}

func applyEncodeOption(e *json.Encoder, options ...EncodeOption) {
	en := &encoder{e: e}
	for _, option := range options {
		option(en)
	}
}

func NewEncoder(w io.Writer, options ...EncodeOption) *json.Encoder {
	e := json.NewEncoder(w)
	applyEncodeOption(e, options...)
	return e
}

func Marshal(v any, options ...EncodeOption) ([]byte, error) {
	buf := bytes.NewBuffer([]byte{})
	if err := NewEncoder(buf, options...).Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func String(v any, options ...EncodeOption) (string, error) {
	b, err := Marshal(v, options...)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func MaybeString(v any, options ...EncodeOption) string {
	s, _ := String(v, options...)
	return s
}

func MaybeBytes(v any, options ...EncodeOption) []byte {
	b, _ := Marshal(v, options...)
	return b
}

type errorReader struct {
	err error
	re  io.Reader
}

func (r *errorReader) Read(b []byte) (int, error) {
	if r.err != nil {
		return 0, r.err
	}
	return r.re.Read(b)
}

func Reader(v any, options ...EncodeOption) io.Reader {
	b, err := Marshal(v, options...)
	return &errorReader{
		err: err,
		re:  bytes.NewReader(b),
	}
}
