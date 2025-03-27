package repository

import (
	"context"
	"tarantool-app/internal/config"
	"tarantool-app/internal/domain"
	"tarantool-app/internal/infrastructure/logging"
	"time"

	"github.com/tarantool/go-tarantool/v2"

	_ "github.com/tarantool/go-tarantool/v2/datetime"
	_ "github.com/tarantool/go-tarantool/v2/decimal"
	_ "github.com/tarantool/go-tarantool/v2/uuid"
)

type Tarantool struct {
	conn *tarantool.Connection
	log  *logging.Logger
}

func NewTarantoolRepository(cfg *config.Config, zlog *logging.Logger) (*Tarantool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	zlog.Info("Initializing database.")

	dialer := tarantool.NetDialer{
		Address:  cfg.Storage.Host + ":" + cfg.Storage.Port,
		User:     cfg.Storage.Credentials.User,
		Password: cfg.Storage.Credentials.Password,
	}
	opts := tarantool.Opts{
		Timeout: time.Second,
	}

	conn, err := tarantool.Connect(ctx, dialer, opts)
	if err != nil {
		zlog.Error("Connection refused.")
		zlog.Debug("Failed to connect to database",
			"address", dialer.Address,
			"user", dialer.User,
			"error", err,
		)
		return nil, err
	}

	zlog.Debug("Connected to database",
		"user", cfg.Storage.Credentials.User,
		"database", cfg.Storage.Host,
	)

	zlog.Info("Database initialization succsessful.")

	return &Tarantool{
		conn: conn,
		log:  zlog,
	}, nil
}

// Gracefully close tarantool connection.
func (tt *Tarantool) Close() {
	if err := tt.conn.CloseGraceful(); err != nil {
		tt.log.Error("Error closing Tarantool connection")
		tt.log.Debug("Failed to close the connection to database",
			"error", err,
		)
	} else {
		tt.log.Info("Database connection closed.")
	}
}

// CRUD operations

func (tt *Tarantool) Insert(rq *domain.AppPack) error {
	tt.log.Debug("Insert request",
		"key", rq.Key,
		"value", rq.Value,
	)

	// We can only check the header to determine the outcome of Insert operation.
	request := tarantool.NewInsertRequest("kv_storage").Tuple(&rq)

	futureResp, err := tt.conn.Do(request).GetResponse()
	if err != nil {
		tt.log.Debug("Failed to Insert data",
			"key", rq.Key,
			"value", rq.Value,
			"error", err,
		)
		return ErrInsertOperationFail
	}

	// If key already exists the response header will contain erorr string.
	// Is there a better way?
	if futureResp.Header().Error != tarantool.ErrorNo && futureResp.Header().Error.String() == "ER_TUPLE_FOUND" {
		tt.log.Debug("response error code", "error", futureResp.Header().Error)
		return ErrAlreadyExists
	}

	return nil
}

func (tt *Tarantool) Select(rq *domain.AppPack) (*domain.AppPack, error) {
	tt.log.Debug("Select request was made",
		"key", rq.Key,
		"value", rq.Value,
	)

	request := tarantool.NewSelectRequest("kv_storage").Key(tarantool.StringKey{S: rq.Key})

	future := tt.conn.Do(request)

	// Check the response without decoding.
	// If there is an error, assume Tarantool failed to do the operation.
	futureResp, err := future.GetResponse()
	if err != nil {
		tt.log.Debug("Failed to Select data",
			"key", rq.Key,
			"value", rq.Value,
			"error", err,
		)
		return &domain.AppPack{}, err
	}

	// If there were no errors, decode and return the result.
	var result []domain.AppPack
	errDecode := futureResp.DecodeTyped(&result)
	if errDecode != nil {
		tt.log.Debug("Failed to decode typed data after Select response returned no errors",
			"key", rq.Key,
			"value", rq.Value,
			"error", errDecode,
		)
		return &domain.AppPack{}, ErrSelectOperationFail
	}

	// If the result is zero length array, then the key does not exist.
	if len(result) == 0 {
		tt.log.Debug("Select returned zero length array",
			"key", rq.Key,
			"value", rq.Value,
		)
		return &domain.AppPack{}, ErrNotFound
	}

	tt.log.Debug("Select response",
		"key", result[0].Key,
		"value", result[0].Value,
	)

	return &result[0], nil
}

func (tt *Tarantool) Update(rq *domain.AppPack) error {
	tt.log.Debug("Update request",
		"key", rq.Key,
		"value", rq.Value,
	)

	request := tarantool.NewUpdateRequest("kv_storage").
		Key(tarantool.StringKey{S: rq.Key}).
		Operations(tarantool.NewOperations().Assign(1, &rq.Value))

	future := tt.conn.Do(request)

	var result []domain.AppPack
	err := future.GetTyped(&result)
	if err != nil {
		tt.log.Debug("Update request failed",
			"key", rq.Key,
			"value", rq.Value,
			err,
		)
		return ErrUpdateOperationFail
	}

	// Key was not found.
	if len(result) == 0 {
		tt.log.Debug("Update returned zero length array",
			"key", rq.Key,
			"value", rq.Value,
		)
		return ErrNotFound
	}

	return nil
}

func (tt *Tarantool) Delete(rq *domain.AppPack) (*domain.AppPack, error) {
	tt.log.Debug("Delete request",
		"key", rq.Key,
		"value", rq.Value,
	)

	request := tarantool.NewDeleteRequest("kv_storage").Key(tarantool.StringKey{S: rq.Key})

	future := tt.conn.Do(request)

	var result []domain.AppPack
	err := future.GetTyped(&result)
	if err != nil {
		tt.log.Debug("Delete request failed",
			"key", rq.Key,
			err,
		)
		return &domain.AppPack{}, ErrDeleteOperationFail
	}

	if len(result) == 0 {
		tt.log.Debug("Delete request returned zero length array.",
			"key", rq.Key,
			"value", rq.Value,
		)
		return &domain.AppPack{}, ErrNotFound
	}

	return &result[0], nil
}
