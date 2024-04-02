package tracer

import (
	"context"
	"log"

	"gopkg.in/DataDog/dd-trace-go.v1/ddtrace/tracer"
)

func Init(host string, service string) {
	log.Println("Starting tracer ddog...")
	tracer.Start(
		tracer.WithAgentAddr(host),
		tracer.WithService(service),
	)
}

func Stop() {
	log.Println("Stopping tracer ddog...")
	tracer.Stop()
}

func StartSpanWithContext(ctx context.Context, name string, tags map[string]string) (tracer.Span, context.Context) {
	span, newCtx := tracer.StartSpanFromContext(ctx, name)
	for k, v := range tags {
		span.SetTag(k, v)
	}

	return span, newCtx
}

func SpanWithContext(ctx context.Context, tags map[string]string) tracer.Span {
	span, _ := tracer.SpanFromContext(ctx)
	for k, v := range tags {
		span.SetTag(k, v)
	}

	return span
}
