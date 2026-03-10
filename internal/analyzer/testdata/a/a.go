package a

import (
	"log/slog"

	"go.uber.org/zap"
)

func slogPackage(password string) {
	slog.Info("Starting server")       // want "log message should start with a lowercase letter"
	slog.Info("привет")                // want "log message should contain only English letters"
	slog.Info("started!!!")            // want "log message should not contain special symbols or emoji"
	slog.Info("password: " + password) // want "log message contains potentially sensitive data"
	slog.Info("started server")
}

func slogMethod(l *slog.Logger) {
	l.Info("Bad") // want "log message should start with a lowercase letter"
}

func zapMethod(z *zap.Logger, token string) {
	z.Info("Bad")             // want "log message should start with a lowercase letter"
	z.Info("token: " + token) // want "log message contains potentially sensitive data"
	z.Info("started")
}
