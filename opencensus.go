package opencensus

import (
	"fmt"
	"strconv"

	propagation "github.com/entropyx/opencensus-propagation"
	"github.com/entropyx/soul/env"
	"github.com/gin-gonic/gin"
	"go.opencensus.io/trace"
)

type Config struct{}

var statusList = map[int]int32{
	200: 0,
	499: 1,
	// 500: 2,
	400: 3,
	504: 4,
	404: 5,
	409: 6,
	403: 7,
	429: 8,
	// 400: 9,
	// 409: 10,
	// 400: 11,
	501: 12,
	500: 13,
	503: 14,
	// 500: 15,
	401: 16,
}

func Middleware(config *Config) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		request := c.Request
		spanCtx, _ := propagation.Extract(propagation.FormatTextMap, propagation.HTTPHeader(request.Header))
		_, span := trace.StartSpanWithRemoteParent(c, fmt.Sprintf("[%s] %s", request.Method, c.FullPath()), spanCtx, trace.WithSampler(setSampler()))
		defer span.End()
		c.Set("span", span)
		span.SpanContext()
		c.Next()
		addHTTP(c, span)
		addParams(c, span)
	})
}

func addHTTP(c *gin.Context, span *trace.Span) {
	request := c.Request
	s := c.Writer.Status()
	status := strconv.Itoa(s)
	if s >= 400 {
		span.AddAttributes(trace.StringAttribute("error.msg", status))
		span.SetStatus(trace.Status{Code: statusList[s]})
	}
	span.AddAttributes(trace.StringAttribute("http.url", request.URL.Path))
	span.AddAttributes(trace.StringAttribute("http.method", request.Method))
	span.AddAttributes(trace.StringAttribute("http.status_code", status))
}

func addParams(c *gin.Context, span *trace.Span) {
	for _, param := range c.Params {
		span.AddAttributes(trace.StringAttribute("http.params."+param.Key, param.Value))
	}
}

func setSampler() trace.Sampler {
	switch env.Mode {
	case env.ModeTest, env.ModeDebug:
		return trace.NeverSample()
	default:
		return trace.AlwaysSample()
	}
}
