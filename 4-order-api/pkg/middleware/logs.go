package middleware

import (
	"4-order-api/configs"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

type Logger struct {
	*logrus.Logger
}

func NewLogger(config *configs.Config) *Logger {
	logger := logrus.New()
	logger.Formatter = new(logrus.JSONFormatter)
	logger.Formatter = new(logrus.TextFormatter)
	logger.Level = logrus.TraceLevel
	logger.Out = os.Stdout

	file, err := os.OpenFile(
		config.Logger.LogFile,
		os.O_CREATE|os.O_WRONLY|os.O_APPEND,
		0666,
	)
	if err != nil {
		logger.Info("Failed to log to file, using only stderr")
		file = nil
	}
	if file != nil {
		logger.SetOutput(io.MultiWriter(os.Stdout, file))
	} else {
		logger.SetOutput(os.Stdout)
	}

	return &Logger{logger}
}

func (logger *Logger) Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		wrapper := &WrapperWriter{
			ResponseWriter: w,
			StatusCode:     http.StatusOK,
		}
		next.ServeHTTP(wrapper, r)

		logEntry := logger.WithFields(logrus.Fields{
			"method":   r.Method,
			"path":     r.URL.Path,
			"status":   wrapper.StatusCode,
			"duration": time.Since(start).String(),
			"tag":      "middleware",
		})

		switch {
		case wrapper.StatusCode >= 500:
			logEntry.Error("Server error")
		case wrapper.StatusCode >= 400:
			logEntry.Warning("Client error")
		case wrapper.StatusCode >= 300:
			logEntry.Info("Redirection")
		default:
			logEntry.Info("Success")
		}
	})
}
