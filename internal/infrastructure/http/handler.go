package http

import (
	"tarantool-app/internal/domain"
	"tarantool-app/internal/usecases"
)

type RequestsHandler struct {
	GetUseCase    *usecases.GetUseCase
	CreateUseCase *usecases.CreateUseCase
	UpdateUseCase *usecases.UpdateUseCase
	DeleteUseCase *usecases.DeleteUseCase
}

func NewRequestHandler(
	getUseCase *usecases.GetUseCase,
	crtUseCase *usecases.CreateUseCase,
	updUseCase *usecases.UpdateUseCase,
	delUseCase *usecases.DeleteUseCase,
) *RequestsHandler {
	return &RequestsHandler{
		GetUseCase:    getUseCase,
		CreateUseCase: crtUseCase,
		UpdateUseCase: updUseCase,
		DeleteUseCase: delUseCase,
	}
}

// GET /kv/{key}
func (r *RequestsHandler) Get(key string) error {
	return nil
}

// POST /kv body: {key: "key", "value": JSON}
func (r *RequestsHandler) Create(key string, value domain.Value) error {
	return nil
}

// PUT kv/{key} body: {"value": JSON}
func (r *RequestsHandler) Set(key string, value domain.Value) error {
	return nil
}

// DELETE kv/{key}
func (r *RequestsHandler) Delete(key string) error {
	return nil
}
