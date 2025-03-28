// Usecase layer. Decouples repository from route handlers.

package usecase

import (
	"tarantool-app/internal/domain"
	"tarantool-app/internal/infrastructure/logging"
	"tarantool-app/internal/repository"

	"github.com/gin-gonic/gin"
)

type CRUD struct {
	repo *repository.Tarantool
	log  *logging.Logger
}

func NewCRUD(repo *repository.Tarantool, zloger *logging.Logger) *CRUD {
	return &CRUD{
		repo: repo,
		log:  zloger,
	}
}

func (crud *CRUD) Create(c *gin.Context, rq *domain.AppPack) error {
	crud.log.Debug("Create request was made",
		"key", rq.Key,
		"value", rq.Value,
	)
	return crud.repo.Insert(rq)
}

func (crud *CRUD) Update(c *gin.Context, rq *domain.AppPack) error {
	crud.log.Debug("Update request was made",
		"key", rq.Key,
		"value", rq.Value,
	)
	return crud.repo.Update(rq)
}

func (crud *CRUD) Delete(c *gin.Context, rq *domain.AppPack) (*domain.AppPack, error) {
	crud.log.Debug("Delete request was made",
		"key", rq.Key,
		"value", rq.Value,
	)
	return crud.repo.Delete(rq)
}

func (crud *CRUD) Read(c *gin.Context, rq *domain.AppPack) (*domain.AppPack, error) {
	crud.log.Debug("Read request was made",
		"key", rq.Key,
		"value", rq.Value,
	)
	resp, err := crud.repo.Select(rq)
	if err != nil {
		return &domain.AppPack{}, err
	}
	return resp, nil
}
