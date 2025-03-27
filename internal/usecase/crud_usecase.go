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
	crud.log.Debug("Create request needs to unmarshal data",
		"data", rq.Value,
	)

	// // var value map[string]any
	// err := json.Unmarshal(rq.Value, &value)
	// if err != nil {
	// 	crud.log.Debug("Failed to unmarshal JSON into `map[string]any`",
	// 		"error", err,
	// 	)
	// 	return domain.AppPostResponse{}, err
	// }

	ttrsp, err := crud.repo.Insert(rq.Key, rq.Value)

	return domain.AppPostResponse{
		Message: "created",
		Key:     ttrsp.Key,
		Size:    uint(len(ttrsp.Value)),
	}, err
}

func (crud *CRUD) Update(c *gin.Context, rq *domain.AppPutRequest) (domain.AppPutResponse, error) {
	crud.log.Debug("Update request needs to unmarshal data",
		"data", rq.Value,
	)

	// var value map[string]any
	// err := json.Unmarshal(rq.Value, &value)
	// if err != nil {
	// 	crud.log.Debug("Failed to unmarshal JSON into `map[string]any`",
	// 		"error", err,
	// 	)
	// 	return domain.AppPutResponse{}, err
	// }
	ttrsp, err := crud.repo.Update(rq.Key, rq.Value)

	return domain.AppPutResponse{
		Message: "updated",
		Key:     rq.Key,
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

	if err != nil {
		return domain.AppGetResponse{}, err
	}

	return domain.AppGetResponse{
		Value: ttrsp.Value,
	}, err
}
