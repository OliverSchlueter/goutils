package containers

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/OliverSchlueter/goutils/sloki"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// ConnectToNatsE2E connects to a local NATS instance without authentication.
func ConnectToNatsE2E() *nats.Conn {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		slog.Error("Could not connect to NATS", sloki.WrapError(err))
		os.Exit(1)
	}

	return nc
}

// ConnectToNats connects to a NATS instance using the provided URL and authentication token.
// Example url: "nats://127.0.0.1:4222"
func ConnectToNats(url, authToken string) *nats.Conn {
	nc, err := nats.Connect(url, nats.Token(authToken))
	if err != nil {
		slog.Error("Could not connect to NATS", sloki.WrapError(err))
		os.Exit(1)
	}

	return nc
}

// ConnectToClickhouseE2E connects to a local Clickhouse instance.
// Username and password are both "admin" by default.
func ConnectToClickhouseE2E(database string) driver.Conn {
	ch, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{"127.0.0.1:9000"},
		Auth: clickhouse.Auth{
			Database: database,
			Username: "admin",
			Password: "admin",
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "myapp", Version: "0.0.1"},
			},
		},
		//Debug: true,
		Debugf: func(format string, v ...interface{}) {
			slog.Debug(fmt.Sprintf(format, v...))
		},
	})
	if err != nil {
		slog.Error("Could not connect to Clickhouse", sloki.WrapError(err))
		os.Exit(1)
	}
	if err := ch.Ping(context.Background()); err != nil {
		slog.Error("Could not ping Clickhouse", sloki.WrapError(err))
		os.Exit(1)
	}

	return ch
}

// ConnectToClickhouse connects to a Clickhouse instance using the provided parameters.
// The addr parameter should be in the format "host:port" (example: "127.0.0.1:9000").
func ConnectToClickhouse(addr, db, user, pwd, appName, appVersion string) driver.Conn {
	ch, err := clickhouse.Open(&clickhouse.Options{
		Addr: []string{addr},
		Auth: clickhouse.Auth{
			Database: db,
			Username: user,
			Password: pwd,
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: appName, Version: appVersion},
			},
		},
		//Debug: true,
		Debugf: func(format string, v ...interface{}) {
			slog.Debug(fmt.Sprintf(format, v...))
		},
	})
	if err != nil {
		slog.Error("Could not connect to Clickhouse", sloki.WrapError(err))
		os.Exit(1)
	}

	if err := ch.Ping(context.Background()); err != nil {
		slog.Error("Could not ping Clickhouse", sloki.WrapError(err))
		os.Exit(1)
	}

	return ch
}

// ConnectToMongoE2E connects to a local MongoDB instance without authentication.
func ConnectToMongoE2E(database string) *mongo.Database {
	mc, err := mongo.Connect(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		slog.Error("Could not connect to MongoDB", sloki.WrapError(err))
		os.Exit(1)
	}
	err = mc.Ping(context.Background(), readpref.Primary())
	if err != nil {
		slog.Error("Could not ping MongoDB", sloki.WrapError(err))
		os.Exit(1)
	}

	return mc.Database(database)
}

// ConnectToMongo connects to a MongoDB instance using the provided connection string and database name.
// Example conn: "mongodb://user:password@localhost:27017"
func ConnectToMongo(conn, db string) *mongo.Database {
	mc, err := mongo.Connect(options.Client().ApplyURI(conn))
	if err != nil {
		slog.Error("Could not connect to MongoDB", sloki.WrapError(err))
		os.Exit(1)
	}

	err = mc.Ping(context.Background(), readpref.Primary())
	if err != nil {
		slog.Error("Could not ping MongoDB", sloki.WrapError(err))
		os.Exit(1)
	}

	return mc.Database(db)
}

// ConnectToRedisE2E connects to a local Redis instance without authentication.
func ConnectToRedisE2E() *redis.Client {
	rc := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0, // use default DB
	})

	ping := rc.Ping(context.Background())
	if ping.Err() != nil {
		slog.Error("Could not connect to Redis", sloki.WrapError(ping.Err()))
		os.Exit(1)
	}

	return rc
}

// ConnectToRedis connects to a Redis instance using the provided address and password.
// Example addr: "127.0.0.1:6379"
func ConnectToRedis(addr, pwd string) *redis.Client {
	rc := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pwd,
		DB:       0, // use default DB
	})

	ping := rc.Ping(context.Background())
	if ping.Err() != nil {
		slog.Error("Could not connect to Redis", sloki.WrapError(ping.Err()))
		os.Exit(1)
	}

	return rc
}
