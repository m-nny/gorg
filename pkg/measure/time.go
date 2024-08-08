package measure

import (
	"log/slog"
	"time"
)

func Timer(msg string) (string, time.Time) {
	return msg, time.Now()
}

func TimerStop(msg string, start time.Time) {
	slog.Info(msg, "time.Since", time.Since(start))
}
