package domain

import (
	"fmt"
	"reflect"

	"github.com/vmihailenco/msgpack/v5"
)

type Payload struct {
	Key   string         `json:"key"`
	Value map[string]any `json:"value"`
}

func (p *Payload) EncodeMsgpack(e *msgpack.Encoder) error {
	payloadLength := reflect.TypeOf(Payload{}).NumField()

	if err := e.EncodeArrayLen(payloadLength); err != nil {
		return err
	}
	if err := e.EncodeString(p.Key); err != nil {
		return err
	}
	if err := e.EncodeMap(p.Value); err != nil {
		return err
	}
	return nil
}

func (p *Payload) DecodeMsgpack(d *msgpack.Decoder) error {
	var err error
	var structLength int
	if structLength, err = d.DecodeArrayLen(); err != nil {
		return err
	}

	payloadLength := reflect.TypeOf(Payload{}).NumField()

	if structLength != payloadLength {
		return fmt.Errorf("array len doesn't match: %d", structLength)
	}
	if p.Key, err = d.DecodeString(); err != nil {
		return err
	}
	if p.Value, err = d.DecodeMap(); err != nil {
		return err
	}
	return nil
}
