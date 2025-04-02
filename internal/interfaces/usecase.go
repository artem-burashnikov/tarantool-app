package interfaces

import (
	"tarantool-app/internal/domain"
)

type UserUseCase interface {
	Create(domain.Payload) error
	Update(domain.Payload) error
	Delete(domain.Payload) (domain.Payload, error)
	Read(domain.Payload) (domain.Payload, error)
}
