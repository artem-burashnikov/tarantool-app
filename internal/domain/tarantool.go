package domain

import (
	"encoding/json"
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
)

// Structure used to communicate with Tarantool.
type TTPack struct {
	Key   string          `json:"key"`
	Value json.RawMessage `json:"value"`
}

// Packs TTPack structure in msgPack array.
func (ttpr *TTPack) EncodeMsgpack(e *msgpack.Encoder) error {
	if err := e.EncodeArrayLen(2); err != nil {
		return err
	}
	if err := e.EncodeString(ttpr.Key); err != nil {
		return err
	}
	if err := e.EncodeBytes(ttpr.Value); err != nil {
		return err
	}
	return nil
}

// Unpacks msgPack array into TTPack structure.
func (ttpr *TTPack) DecodeMsgpack(d *msgpack.Decoder) error {
	var err error
	var l int
	if l, err = d.DecodeArrayLen(); err != nil {
		return err
	}
	if l != 2 {
		return fmt.Errorf("array len doesn't match: %d", l)
	}
	if ttpr.Key, err = d.DecodeString(); err != nil {
		return err
	}
	if ttpr.Value, err = d.DecodeBytes(); err != nil {
		return err
	}
	return nil
}
