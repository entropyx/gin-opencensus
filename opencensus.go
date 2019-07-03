package opencensus

import (
	"fmt"

	propagation "github.com/entropyx/opencensus-propagation"
	"github.com/entropyx/soul/env"
	"github.com/gin-gonic/gin"
	"go.opencensus.io/trace"
)

type Config struct{}

func Middleware(config *Config) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		request := c.Request
		spanCtx, _ := propagation.Extract(propagation.FormatTextMap, propagation.HTTPHeader(request.Header))
		_, span := trace.StartSpanWithRemoteParent(c, fmt.Sprintf("[%s] %s", request.Method, request.URL.Path), spanCtx, trace.WithSampler(setSampler()))
		defer span.End()
		c.Set("span", span)
		c.Next()
	})
}

func setSampler() trace.Sampler {
	switch env.Mode {
	case env.ModeTest, env.ModeDebug:
		return trace.NeverSample()
	default:
		return trace.AlwaysSample()
	}
}
