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
	tt  *tarantool.Connection
	log *logging.Logger
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
		tt:  conn,
		log: zlog,
	}, nil
}

// Gracefully close tarantool connection.
func (t *Tarantool) Close() {
	if err := t.tt.CloseGraceful(); err != nil {
		t.log.Error("Error closing Tarantool connection")
		t.log.Debug("Failed to close the connection to database",
			"error", err,
		)
	} else {
		t.log.Info("Database connection closed.")
	}
}

func (tt *Tarantool) SaveKeyValue(key string, value domain.AppPostRequest) error {
	return nil
}

func (tt *Tarantool) GetByKey(key string) (domain.AppPostRequest, error) {
	var value domain.AppPostRequest
	return value, nil
}

func (tt *Tarantool) SetValueByKey(key string, value domain.AppPostRequest) error {
	return nil
}

func (tt *Tarantool) DeleteByKey(key string) error {
	return nil
}

// 	// Close connection
// 	conn.CloseGraceful()
// 	log.Debug("Connection is closed", log.Str("user", cfg.Storage.Credentials.User))

// Interact with the database
// ...
// Insert data
// tuples := [][]interface{}{
// 	{1, "Roxette", 1986},
// 	{2, "Scorpions", 1965},
// 	{3, "Ace of Base", 1987},
// 	{4, "The Beatles", 1960},
// }
// var futures []*tarantool.Future
// for _, tuple := range tuples {
// 	request := tarantool.NewInsertRequest("bands").Tuple(tuple)
// 	futures = append(futures, conn.Do(request))
// }
// fmt.Println("Inserted tuples:")
// for _, future := range futures {
// 	result, err := future.Get()
// 	if err != nil {
// 		fmt.Println("Got an error:", err)
// 	} else {
// 		fmt.Println(result)
// 	}
// }

// // Select by primary key
// data, err := conn.Do(
// 	tarantool.NewSelectRequest("bands").
// 		Limit(10).
// 		Iterator(tarantool.IterEq).
// 		Key([]interface{}{uint(1)}),
// ).Get()
// if err != nil {
// 	fmt.Println("Got an error:", err)
// }
// fmt.Println("Tuple selected by the primary key value:", data)

// // Select by secondary key
// data, err = conn.Do(
// 	tarantool.NewSelectRequest("bands").
// 		Index("band").
// 		Limit(10).
// 		Iterator(tarantool.IterEq).
// 		Key([]interface{}{"The Beatles"}),
// ).Get()
// if err != nil {
// 	fmt.Println("Got an error:", err)
// }
// fmt.Println("Tuple selected by the secondary key value:", data)

// // Update
// data, err = conn.Do(
// 	tarantool.NewUpdateRequest("bands").
// 		Key(tarantool.IntKey{2}).
// 		Operations(tarantool.NewOperations().Assign(1, "Pink Floyd")),
// ).Get()
// if err != nil {
// 	fmt.Println("Got an error:", err)
// }
// fmt.Println("Updated tuple:", data)

// // Delete
// data, err = conn.Do(
// 	tarantool.NewDeleteRequest("bands").
// 		Key([]interface{}{uint(5)}),
// ).Get()
// if err != nil {
// 	fmt.Println("Got an error:", err)
// }
// fmt.Println("Deleted tuple:", data)
