package http

import (
	"net/http"

	"go.uber.org/zap"
)

type Opt func(l *loadShedderMiddleware)

func WithZapLogger(log *zap.Logger) Opt {
	return func(l *loadShedderMiddleware) {
		l.log = log
	}
}

func WithOverloadHandler(h http.Handler) Opt {
	return func(l *loadShedderMiddleware) {
		l.overloadHandler = h
	}
}
