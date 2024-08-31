package telemetry

import (
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/uber/jaeger-client-go"
	traceconfig "github.com/uber/jaeger-client-go/config"
	"github.com/uber/jaeger-lib/metrics/prometheus"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

func Listen(ctx context.Context, zap *zap.SugaredLogger, telemetryAddr string) {
	MustSetup(ctx, "kod")

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	zap.Infof("Listening on %s", telemetryAddr)
	if err := http.ListenAndServe(telemetryAddr, mux); err != nil {
		zap.Errorln(err)
	}
}

func MustSetup(ctx context.Context, serviceName string) {
	cfg := traceconfig.Configuration{
		ServiceName: serviceName,
		Sampler: &traceconfig.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &traceconfig.ReporterConfig{
			CollectorEndpoint: "http://localhost:14268/api/traces",
			LogSpans:          true,
		},
	}

	tracer, closer, err := cfg.NewTracer(traceconfig.Logger(jaeger.StdLogger), traceconfig.Metrics(prometheus.New()))
	if err != nil {
		log.Fatalf("ERROR: cannot init Jaeger %s", err)
	}

	go func() {
		onceCloser := sync.OnceFunc(func() {
			log.Warn(ctx, "closing tracer")
			if err = closer.Close(); err != nil {
				log.Errorf("error closing tracer: %s", err)
			}
		})

		for {
			<-ctx.Done()
			onceCloser()
		}
	}()

	opentracing.SetGlobalTracer(tracer)
}