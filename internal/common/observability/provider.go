package observability

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

// Config holds observability configuration
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	OTLPEndpoint   string // e.g., "localhost:4317" for OTLP gRPC
	EnableTracing  bool
	EnableMetrics  bool
	SampleRate     float64 // 0.0 to 1.0
}

// Provider holds observability providers
type Provider struct {
	TracerProvider     *trace.TracerProvider
	MeterProvider      *metric.MeterProvider
	PrometheusExporter *prometheus.Exporter
	Shutdown           func(context.Context) error
}

// New initializes OpenTelemetry with tracing and metrics
func New(ctx context.Context, cfg Config) (*Provider, error) {
	if cfg.SampleRate == 0 {
		cfg.SampleRate = 1.0 // Default to 100% sampling
	}

	// Create resource with service information
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.DeploymentEnvironment(cfg.Environment),
			attribute.String("library.language", "go"),
		),
		resource.WithHost(),
		resource.WithProcess(),
		resource.WithTelemetrySDK(),
	)
	if err != nil {
		return nil, fmt.Errorf("create resource: %w", err)
	}

	provider := &Provider{}
	var shutdownFuncs []func(context.Context) error

	// Initialize tracing
	if cfg.EnableTracing {
		tp, shutdown, err := initTracer(ctx, cfg, res)
		if err != nil {
			return nil, fmt.Errorf("init tracer: %w", err)
		}
		provider.TracerProvider = tp
		shutdownFuncs = append(shutdownFuncs, shutdown)

		// Set global tracer provider
		otel.SetTracerProvider(tp)
	}

	// Initialize metrics
	if cfg.EnableMetrics {
		mp, promExporter, shutdown, err := initMeter(ctx, cfg, res)
		if err != nil {
			return nil, fmt.Errorf("init meter: %w", err)
		}
		provider.MeterProvider = mp
		provider.PrometheusExporter = promExporter
		shutdownFuncs = append(shutdownFuncs, shutdown)

		// Set global meter provider
		otel.SetMeterProvider(mp)
	}

	// Set global propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Create combined shutdown function
	provider.Shutdown = func(ctx context.Context) error {
		var errs []error
		for _, fn := range shutdownFuncs {
			if err := fn(ctx); err != nil {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return fmt.Errorf("shutdown errors: %v", errs)
		}
		return nil
	}

	return provider, nil
}

func initTracer(ctx context.Context, cfg Config, res *resource.Resource) (*trace.TracerProvider, func(context.Context) error, error) {
	var exporter trace.SpanExporter
	var err error

	if cfg.OTLPEndpoint != "" {
		// OTLP exporter for Jaeger/Tempo/etc.
		exporter, err = otlptracegrpc.New(ctx,
			otlptracegrpc.WithEndpoint(cfg.OTLPEndpoint),
			otlptracegrpc.WithInsecure(),
		)
		if err != nil {
			return nil, nil, fmt.Errorf("create OTLP trace exporter: %w", err)
		}
	} else {
		// No-op exporter when no endpoint configured
		exporter = noopSpanExporter{}
	}

	// Create sampler based on sample rate
	var sampler trace.Sampler
	if cfg.SampleRate >= 1.0 {
		sampler = trace.AlwaysSample()
	} else if cfg.SampleRate <= 0 {
		sampler = trace.NeverSample()
	} else {
		sampler = trace.TraceIDRatioBased(cfg.SampleRate)
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter,
			trace.WithBatchTimeout(5*time.Second),
			trace.WithMaxExportBatchSize(512),
		),
		trace.WithResource(res),
		trace.WithSampler(sampler),
	)

	shutdown := func(ctx context.Context) error {
		return tp.Shutdown(ctx)
	}

	return tp, shutdown, nil
}

func initMeter(ctx context.Context, cfg Config, res *resource.Resource) (*metric.MeterProvider, *prometheus.Exporter, func(context.Context) error, error) {
	// Prometheus exporter for /metrics endpoint
	promExporter, err := prometheus.New()
	if err != nil {
		return nil, nil, nil, fmt.Errorf("create prometheus exporter: %w", err)
	}

	readers := []metric.Reader{promExporter}

	// OTLP exporter for centralized metrics (optional)
	if cfg.OTLPEndpoint != "" {
		otlpExporter, err := otlpmetricgrpc.New(ctx,
			otlpmetricgrpc.WithEndpoint(cfg.OTLPEndpoint),
			otlpmetricgrpc.WithInsecure(),
		)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("create OTLP metric exporter: %w", err)
		}
		readers = append(readers, metric.NewPeriodicReader(otlpExporter,
			metric.WithInterval(15*time.Second),
		))
	}

	mp := metric.NewMeterProvider(
		metric.WithResource(res),
		metric.WithReader(promExporter),
	)

	shutdown := func(ctx context.Context) error {
		return mp.Shutdown(ctx)
	}

	return mp, promExporter, shutdown, nil
}

// noopSpanExporter is a no-op exporter when tracing is not configured
type noopSpanExporter struct{}

func (noopSpanExporter) ExportSpans(ctx context.Context, spans []trace.ReadOnlySpan) error {
	return nil
}

func (noopSpanExporter) Shutdown(ctx context.Context) error {
	return nil
}
