package http

import (
	"math/rand"
	"net/http"

	"github.com/fhke/loadsheding-go/usage"
	"go.uber.org/zap"
)

type loadShedderMiddleware struct {
	log *zap.Logger

	tracker         usage.Tracker
	limitSoft       float64
	limitHard       float64
	next            http.Handler
	overloadHandler http.Handler
}

func NewLoadShedderMiddleware(tracker usage.Tracker, limitSoft, limitHard float64, next http.Handler, opts ...Opt) http.Handler {
	lsm := &loadShedderMiddleware{
		log:             zap.NewNop(),
		tracker:         tracker,
		limitSoft:       limitSoft,
		limitHard:       limitHard,
		next:            next,
		overloadHandler: http.HandlerFunc(defaultOverloadHandler),
	}

	for _, opt := range opts {
		opt(lsm)
	}

	return lsm
}

func (l *loadShedderMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	nextHandler := l.next
	if l.shouldReject() {
		nextHandler = l.overloadHandler
	}
	nextHandler.ServeHTTP(w, r)
}

func (l *loadShedderMiddleware) shouldReject() bool {
	overloadFactor := l.overloadFactor()
	randNum := rand.Float64()
	reject := randNum < overloadFactor
	log := l.log.With(
		zap.Float64("overload_factor", overloadFactor),
		zap.Float64("overload_rand", randNum),
	)

	if reject {
		log.Warn("Rejecting request")
	} else {
		log.Debug("Permitting request")
	}

	return reject
}

func (l *loadShedderMiddleware) overloadFactor() float64 {
	currentUtilization := l.tracker.Utilization()
	limitRange := l.limitHard - l.limitSoft
	overSoftLimit := currentUtilization - l.limitSoft
	overloadFactor := max(overSoftLimit/limitRange, 0.0)

	return overloadFactor
}

func defaultOverloadHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusServiceUnavailable)
	w.Write([]byte("server overloaded\n"))
}
