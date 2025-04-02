package repository

import (
	"context"
	"os"
	"tarantool-app/config"
	"tarantool-app/internal/domain"
	"tarantool-app/internal/interfaces"
	"time"

	"github.com/tarantool/go-iproto"
	"github.com/tarantool/go-tarantool/v2"

	_ "github.com/tarantool/go-tarantool/v2/datetime"
	_ "github.com/tarantool/go-tarantool/v2/decimal"
	_ "github.com/tarantool/go-tarantool/v2/uuid"
)

type Tarantool struct {
	conn *tarantool.Connection
	log  interfaces.Logger
}

var _ interfaces.Repository = Tarantool{} // Tarantool must satisfy Repository

func NewTarantoolRepository(cfg *config.Config, log interfaces.Logger) (Tarantool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	// TODO: fixme
	dialer := tarantool.NetDialer{
		Address:  cfg.Storage.URI,
		User:     os.Getenv("TT_USER"),
		Password: os.Getenv("TT_PASSWORD"),
	}
	opts := tarantool.Opts{
		Timeout: time.Second,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		return Tarantool{}, err
	}

	return Tarantool{
		conn: conn,
		log:  log,
	}, nil
}

func (tt Tarantool) Close() {
	if err := tt.conn.CloseGraceful(); err != nil {
		tt.log.Error("Error closing Tarantool connection")
	} else {
		tt.log.Info("Database connection closed.")
	}
}

// POST ---> Insert
func (tt Tarantool) Insert(rq domain.Payload) error {
	request := tarantool.NewInsertRequest("kv_storage").Tuple(&rq)

	futureResp, err := tt.conn.Do(request).GetResponse()

	if futureResp.Header().Error != tarantool.ErrorNo {
		if tntErr, ok := err.(tarantool.Error); !ok || tntErr.Code == iproto.ER_TUPLE_FOUND {
			return ErrAlreadyExists
		} else {
			return ErrInsertOperationFail
		}
	}

	return nil
}

// GET ---> Select
func (tt Tarantool) Select(rq domain.Payload) (domain.Payload, error) {
	request := tarantool.NewSelectRequest("kv_storage").Key(tarantool.StringKey{S: rq.Key})

	future := tt.conn.Do(request)

	futureResp, err := future.GetResponse()
	if err != nil {
		return domain.Payload{}, err
	}

	var result []domain.Payload
	errDecode := futureResp.DecodeTyped(&result)
	if errDecode != nil {
		return domain.Payload{}, ErrSelectOperationFail
	}

	if len(result) == 0 {
		return domain.Payload{}, ErrNotFound
	}

	return result[0], nil
}

// PUT ---> Update
func (tt Tarantool) Update(rq domain.Payload) error {
	request := tarantool.NewUpdateRequest("kv_storage").
		Key(tarantool.StringKey{S: rq.Key}).
		Operations(tarantool.NewOperations().Assign(1, &rq.Value))

	future := tt.conn.Do(request)

	var result []domain.Payload
	err := future.GetTyped(&result)
	if err != nil {
		return ErrUpdateOperationFail
	}

	if len(result) == 0 {
		return ErrNotFound
	}

	return nil
}

// DELETE ---> Delete
func (tt Tarantool) Delete(rq domain.Payload) (domain.Payload, error) {
	request := tarantool.NewDeleteRequest("kv_storage").Key(tarantool.StringKey{S: rq.Key})

	future := tt.conn.Do(request)

	var result []domain.Payload
	err := future.GetTyped(&result)
	if err != nil {
		return domain.Payload{}, ErrDeleteOperationFail
	}

	if len(result) == 0 {
		return domain.Payload{}, ErrNotFound
	}

	return result[0], nil
}
