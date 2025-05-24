package middleware

import (
	"github.com/nats-io/nats.go"
	"log/slog"
	"time"
)

func NatsLogging(next nats.MsgHandler) nats.MsgHandler {
	return func(msg *nats.Msg) {
		startTime := time.Now()

		next(msg)

		elapsedTime := time.Since(startTime)

		slog.Info(
			"NATS message received",
			slog.String("subject", msg.Subject),
			slog.String("reply", msg.Reply),
			slog.Int64("elapsed_time", elapsedTime.Milliseconds()),
		)
	}
}
