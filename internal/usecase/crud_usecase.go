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

func (crud *CRUD) Create(c *gin.Context, rq *domain.AppPostRequest) (domain.AppPostResponse, error) {
	crud.log.Debug("Inserting key and value into database",
		"key", rq.Key,
		"value", rq.Value,
	)
	ttrsp, err := crud.repo.Insert(rq.Key, rq.Value)

	crud.log.Debug("Inserted succesfully. Returning to handler.")
	return domain.AppPostResponse{
		Message: "created",
		Key:     ttrsp.Key,
		Size:    uint(len(ttrsp.Value)),
	}, err
}

func (crud *CRUD) Update(c *gin.Context, rq *domain.AppPutRequest) (domain.AppPutResponse, error) {
	ttrsp, err := crud.repo.Update(rq.Key, rq.Value)

	return domain.AppPutResponse{
		Message: "updated",
		Key:     ttrsp.Key,
		Size:    uint(len(ttrsp.Value)),
	}, err
}

func (crud *CRUD) Delete(c *gin.Context, rq *domain.AppDeleteRequest) (domain.AppDeleteResponse, error) {
	ttrsp, err := crud.repo.Delete(rq.Key)

	return domain.AppDeleteResponse{
		Message: "deleted",
		Key:     ttrsp.Key,
	}, err
}

func (crud *CRUD) Read(c *gin.Context, rq *domain.AppGetRequest) (domain.AppGetResponse, error) {
	ttrsp, err := crud.repo.Select(rq.Key)

	return domain.AppGetResponse{
		Value: ttrsp.Value,
	}, err
}
