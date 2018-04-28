package logging

import (
	"net/http"
	"reflect"
	"time"

	"github.com/sirupsen/logrus"
)

type logging struct {
	logger  logrus.FieldLogger
	handler http.Handler
}

func NewLogging(logger logrus.FieldLogger, handler http.Handler) http.Handler {
	return logging{
		logger:  logger.WithField("build.package", reflect.TypeOf(logging{}).PkgPath()),
		handler: handler,
	}
}

func (h logging) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()

	logWriter := &responseLogger{w: w, status: http.StatusOK}

	h.handler.ServeHTTP(logWriter, r)

	logger := h.logger

	if context := r.Context().Value(LoggerContext); context != nil {
		logger = logger.WithFields(context.(logrus.Fields))
	}

	t1 := time.Now()

	logger.WithFields(logrus.Fields{
		"http.response.time":     t1.Format(time.RFC3339),
		"http.response.duration": t1.Sub(t0) / time.Millisecond,
		"http.response.status":   logWriter.Status(),
		"http.response.size":     logWriter.Size(),
	}).Infof("completed")
}

type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

func (l *responseLogger) Status() int {
	return l.status
}

func (l *responseLogger) Size() int {
	return l.size
}

func (l *responseLogger) Flush() {
	f, ok := l.w.(http.Flusher)
	if ok {
		f.Flush()
	}
}
