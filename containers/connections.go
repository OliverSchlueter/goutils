package containers

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/OliverSchlueter/goutils/sloki"
	_ "github.com/mattn/go-sqlite3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/readpref"
)

// ConnectToMinIOE2E connects to a local MinIO instance without authentication.
func ConnectToMinIOE2E() *minio.Client {
	opts := &minio.Options{
		Creds:  credentials.NewStaticV4("adminadmin", "admin", ""),
		Secure: false,
	}
	minioClient, err := minio.New("localhost:9090", opts)
	if err != nil {
		slog.Error("Could not connect to MinIO", sloki.WrapError(err))
		os.Exit(1)
	}

	slog.Info("Connected to MinIO")
	return minioClient
}

// ConnectToMinIO connects to a MinIO instance using the provided parameters.
func ConnectToMinIO(endpoint, user, password string) *minio.Client {
	opts := &minio.Options{
		Creds: credentials.NewStaticV4(user, password, ""),
	}
	minioClient, err := minio.New(endpoint, opts)
	if err != nil {
		slog.Error("Could not connect to MinIO", sloki.WrapError(err))
		os.Exit(1)
	}

	slog.Info("Connected to MinIO")
	return minioClient
}

// DisconnectMinIO closes the connection to the MinIO server.
func DisconnectMinIO(minioClient *minio.Client) {
	// MinIO client does not have a close method, so we just log the disconnection.
	slog.Info("Disconnected from MinIO")
}

// ConnectToNatsE2E connects to a local NATS instance without authentication.
func ConnectToNatsE2E() *nats.Conn {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		slog.Error("Could not connect to NATS", sloki.WrapError(err))
		os.Exit(1)
	}

	slog.Info("Connected to NATS")
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

	slog.Info("Connected to NATS")
	return nc
}

// DisconnectNats closes the connection to the NATS server.
func DisconnectNats(nc *nats.Conn) {
	nc.Close()
	slog.Info("Disconnected from NATS")
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

	slog.Info("Connected to Clickhouse")
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

	slog.Info("Connected to Clickhouse")
	return ch
}

// DisconnectClickhouse closes the connection to the Clickhouse server.
func DisconnectClickhouse(ch driver.Conn) {
	if err := ch.Close(); err != nil {
		slog.Error("Could not close Clickhouse connection", sloki.WrapError(err))
		return
	}

	slog.Info("Disconnected from Clickhouse")
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

	slog.Info("Connected to MongoDB")
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

	slog.Info("Connected to MongoDB")
	return mc.Database(db)
}

// DisconnectMongo disconnects from the MongoDB server.
func DisconnectMongo(mc *mongo.Database) {
	if err := mc.Client().Disconnect(context.Background()); err != nil {
		slog.Error("Could not disconnect from MongoDB", sloki.WrapError(err))
		return
	}

	slog.Info("Disconnected from MongoDB")
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

	slog.Info("Connected to Redis")
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

	slog.Info("Connected to Redis")
	return rc
}

// DisconnectRedis closes the connection to the Redis server.
func DisconnectRedis(rc *redis.Client) {
	if err := rc.Close(); err != nil {
		slog.Error("Could not close Redis connection", sloki.WrapError(err))
		return
	}

	slog.Info("Disconnected from Redis")
}

// ConnectSqlite connects to a SQLite database at the given path.
func ConnectSqlite(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		slog.Error("Failed to open sqlite database", sloki.WrapError(err))
		os.Exit(1)
	}

	if err := db.Ping(); err != nil {
		slog.Error("Failed to ping sqlite database", sloki.WrapError(err))
		os.Exit(1)
	}

	slog.Info("Connected to SQLite")
	return db
}

// DisconnectSqlite closes the connection to the SQLite database.
func DisconnectSqlite(db *sql.DB) {
	if err := db.Close(); err != nil {
		slog.Error("Failed to close sqlite database", sloki.WrapError(err))
	}

	slog.Info("Disconnected from SQLite")
}
