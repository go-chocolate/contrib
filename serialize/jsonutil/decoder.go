package jsonutil

import (
	"bytes"
	"encoding/json"
	"io"
)

type decoder struct {
	d *json.Decoder
}

type DecodeOption func(d *decoder)

func WithDisallowUnknownFields() DecodeOption {
	return func(d *decoder) {
		d.d.DisallowUnknownFields()
	}
}

func applyDecodeOptions(d *json.Decoder, options ...DecodeOption) {
	de := &decoder{d: d}
	for _, option := range options {
		option(de)
	}
}

func NewDecoder(reader io.Reader, options ...DecodeOption) *json.Decoder {
	d := json.NewDecoder(reader)
	d.UseNumber()
	applyDecodeOptions(d, options...)
	return d
}

func Unmarshal(b []byte, v any, options ...DecodeOption) error {
	return NewDecoder(bytes.NewReader(b), options...).Decode(v)
}
