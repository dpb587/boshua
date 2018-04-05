package middleware

import (
	"context"
	"net/http"
	"time"

	uuid "github.com/nu7hatch/gouuid"
	"github.com/sirupsen/logrus"
)

const LoggerContext = iota

type loggerContext struct {
	handler http.Handler
}

func NewLoggerContext(handler http.Handler) http.Handler {
	return loggerContext{
		handler: handler,
	}
}

func (h loggerContext) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var uuidString string

	uuidResult, err := uuid.NewV4()
	if err != nil {
		// TODO bail?
		uuidString = "unknown"
	} else {
		uuidString = uuidResult.String()
	}

	ctx := context.WithValue(r.Context(), LoggerContext, logrus.Fields{
		"http.request.time":        time.Now().Format(time.RFC3339),
		"http.request.id":          uuidString,
		"http.request.remote_addr": r.RemoteAddr,
		"http.request.method":      r.Method,
		"http.request.uri":         r.RequestURI,
		"http.request.user_agent":  r.Header.Get("user-agent"),
	})

	h.handler.ServeHTTP(w, r.WithContext(ctx))
}
