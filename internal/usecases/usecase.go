package usecases

import (
	"tarantool-app/internal/domain"
	"tarantool-app/internal/interfaces"
)

type UserUseCase struct {
	repo interfaces.Repository
	log  interfaces.Logger
}

func NewUserUseCase(repo interfaces.Repository, log interfaces.Logger) UserUseCase {
	return UserUseCase{repo: repo, log: log}
}

func (uc UserUseCase) Create(ap domain.Payload) error {
	return uc.repo.Insert(ap)
}

func (uc UserUseCase) Update(ap domain.Payload) error {
	return uc.repo.Update(ap)
}

func (uc UserUseCase) Delete(ap domain.Payload) (domain.Payload, error) {
	return uc.repo.Delete(ap)
}

func (uc UserUseCase) Read(ap domain.Payload) (domain.Payload, error) {
	return uc.repo.Select(ap)
}
