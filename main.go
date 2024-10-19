package main

import (
	"math"
	"math/rand/v2"
	nethttp "net/http"
	"time"

	"github.com/fhke/loadsheding-go/http"
	"github.com/fhke/loadsheding-go/testserver"
	"github.com/fhke/loadsheding-go/usage"
	"github.com/fhke/loadsheding-go/usage/cpu"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

func main() {
	log := lo.Must(zap.NewDevelopment())
	cpuTracker := usage.NewBackgroundTracker(lo.Must(cpu.New()), time.Millisecond*200)

	go func() {
		err := nethttp.ListenAndServe(
			"127.0.0.1:8000",
			http.
				NewLoadShedderMiddleware(cpuTracker, 2.5, 7, testserver.NewHandler(), http.WithZapLogger(log.Named("shedder"))),
		)
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		for range 20 {
			log.Info("Adding a worker")
			go func() {
				for {
					math.Sqrt(rand.Float64())
				}
			}()
			time.Sleep(time.Second * 20)
		}
	}()

	select {}
}
