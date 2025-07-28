package logging

import (
	"log/slog"
	"os"
)

func getLogLevel() slog.Level {
	switch os.Getenv("LOG_LEVEL") {
	case "DEBUG", "-4":
		return slog.LevelDebug
	case "INFO", "0":
		return slog.LevelInfo
	case "WARN", "4":
		return slog.LevelWarn
	case "ERROR", "8":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

var Logger *slog.Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true, Level: getLogLevel()}))
