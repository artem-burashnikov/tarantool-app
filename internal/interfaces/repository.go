package interfaces

import "tarantool-app/internal/domain"

type Repository interface {
	Insert(domain.Payload) error
	Select(domain.Payload) (domain.Payload, error)
	Update(domain.Payload) error
	Delete(domain.Payload) (domain.Payload, error)
	Close()
}
