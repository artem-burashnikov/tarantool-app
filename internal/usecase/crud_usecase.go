package usecase

import (
	"tarantool-app/internal/domain"
	"tarantool-app/internal/repository"

	"github.com/gin-gonic/gin"
)

type CRUD struct {
	repo *repository.Tarantool
}

func NewCRUD(repo *repository.Tarantool) *CRUD {
	return &CRUD{repo: repo}
}

func (crud *CRUD) Create(c *gin.Context, key string, value domain.Value) error {
	return nil
}

func (crud *CRUD) Read(c *gin.Context, key string) (domain.Value, error) {
	var value domain.Value
	return value, nil
}

func (crud *CRUD) Update(c *gin.Context, key string, value domain.Value) error {
	return nil
}

func (crud *CRUD) Delete(c *gin.Context, key string) (domain.Value, error) {
	var value domain.Value
	return value, nil
}
