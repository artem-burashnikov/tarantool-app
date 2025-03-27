package repository

import (
	"context"
	"encoding/json"
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

func (tt *Tarantool) Insert(key string, value json.RawMessage) (domain.TTPack, error) {
	tt.log.Debug("Insert request",
		"key", key,
		"value", value,
	)

	ttrq := domain.TTPack{
		Key:   key,
		Value: value,
	}

	request := tarantool.NewInsertRequest("kv_storage").Tuple(&ttrq)

	futureResp, err := tt.conn.Do(request).GetResponse()
	if err != nil {
		tt.log.Debug("Failed to Insert data",
			"key", ttrq.Key,
			"value", ttrq.Value,
			"error", err,
		)
		return domain.TTPack{}, ErrInsertOperationFail
	}

	if futureResp.Header().Error != tarantool.ErrorNo && futureResp.Header().Error.String() == "ER_TUPLE_FOUND" {
		tt.log.Debug("response error code", "error", futureResp.Header().Error)
		return domain.TTPack{}, ErrAlreadyExists
	}

	tt.log.Debug("Insert result",
		"key", key,
		"value", value,
	)

	return ttrq, nil
}

func (tt *Tarantool) Select(key string) (domain.TTPack, error) {
	tt.log.Debug("Select request",
		"key", key,
	)
	request := tarantool.NewSelectRequest("kv_storage").Key(tarantool.StringKey{S: key})

	var resp []domain.TTPack
	err := tt.conn.Do(request).GetTyped(&resp)
	if err != nil {
		tt.log.Debug("Failed to Select data",
			"key", key,
			"error", err,
		)
		return domain.TTPack{}, ErrSelectOperationFail
	}

	if len(resp) == 0 {
		return domain.TTPack{}, ErrNotFound
	}

	tt.log.Debug("Select response",
		"key", key,
	)

	return resp[0], nil
}

func (tt *Tarantool) Update(key string, value json.RawMessage) (domain.TTPack, error) {
	tt.log.Debug("Update request",
		"key", key,
		"value", value,
	)

	request := tarantool.NewUpdateRequest("kv_storage").
		Key(tarantool.StringKey{S: key}).
		Operations(tarantool.NewOperations().Assign(2, value))

	var result []json.RawMessage
	err := tt.conn.Do(request).GetTyped(&result)
	if err != nil {
		tt.log.Debug("Update request failed",
			"key", key,
			"value", value,
			"error", err,
		)
		return domain.TTPack{}, ErrUpdateOperationFail
	}

	if len(result) == 0 {
		tt.log.Debug("Update returned zero length array",
			"key", key,
			"value", value,
		)
		return domain.TTPack{}, ErrNotFound
	}

	return domain.TTPack{Key: key, Value: result[0]}, nil
}

func (tt *Tarantool) Delete(key string) (domain.TTPack, error) {
	tt.log.Debug("Delete request",
		"key", key,
	)

	request := tarantool.NewDeleteRequest("kv_storage").Key(tarantool.StringKey{S: key})

	var result []domain.TTPack
	err := tt.conn.Do(request).GetTyped(&result)
	if err != nil {
		tt.log.Debug("Delete request failed",
			"key", key,
			"error", err,
		)
		return domain.TTPack{}, ErrDeleteOperationFail
	}

	if len(result) == 0 {
		return domain.TTPack{}, ErrNotFound
	}

	return result[0], nil
}
