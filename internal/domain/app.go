// Structures used internally to handle http requests and responses.

package domain

import (
	"fmt"

	"github.com/vmihailenco/msgpack/v5"
)

type AppPack struct {
	Key   string         `json:"key"`
	Value map[string]any `json:"value"`
}

// Packs TTPack structure in msgPack array.
func (ttpr *AppPack) EncodeMsgpack(e *msgpack.Encoder) error {
	if err := e.EncodeArrayLen(2); err != nil {
		return err
	}
	if err := e.EncodeString(ttpr.Key); err != nil {
		return err
	}
	if err := e.EncodeMap(ttpr.Value); err != nil {
		return err
	}
	return nil
}

// Unpacks msgPack array into TTPack structure.
func (ttpr *AppPack) DecodeMsgpack(d *msgpack.Decoder) error {
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
	if ttpr.Value, err = d.DecodeMap(); err != nil {
		return err
	}
	return nil
}
