package repository

import "tarantool-app/internal/domain"

type Tarantool struct{}

func NewTarantoolRepository() *Tarantool {
	return &Tarantool{}
}

func (tt *Tarantool) SaveKeyValue(key string, value domain.Value) error {
	return nil
}

func (tt *Tarantool) GetByKey(key string) (domain.Value, error) {
	var value domain.Value
	return value, nil
}

func (tt *Tarantool) SetValueByKey(key string, value domain.Value) error {
	return nil
}

func (tt *Tarantool) DeleteByKey(key string) error {
	return nil
}
