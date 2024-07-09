package opentelemetry

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"web"
)

type MiddlewareBuilder struct {
	Tracer trace.Tracer
}

const instrumentationName = "gitee.com/ObambooO/go_geektime_bootcamp/web/middlewares/opentelemetry"

func (m MiddlewareBuilder) Build() web.Middleware {
	if m.Tracer == nil {
		m.Tracer = otel.GetTracerProvider().Tracer(instrumentationName)
	}
	return func(next web.HandleFunc) web.HandleFunc {
		return func(ctx *web.Context) {

			reqCtx := ctx.Req.Context()

			// 尝试和客户端的trace结合
			reqCtx = otel.GetTextMapPropagator().Extract(reqCtx, propagation.HeaderCarrier(ctx.Req.Header))

			_, span := m.Tracer.Start(reqCtx, "unknown")
			defer span.End()

			span.SetAttributes(attribute.String("http.method", ctx.Req.Method))
			span.SetAttributes(attribute.String("http.url", ctx.Req.URL.String()))
			span.SetAttributes(attribute.String("http.schema", ctx.Req.URL.Scheme))
			span.SetAttributes(attribute.String("http.host", ctx.Req.Host))

			// 你这里还可以继续加

			// 直接调用下一步
			next(ctx)
			// 这个是只有执行完next才可能有值
			span.SetName(ctx.MatchedRoute)
		}
	}
}

// 支持自己传参
//func NewMiddlewareBuilder(tracer trace.Tracer) *MiddlewareBuilder {
//	return &MiddlewareBuilder{}
//}
