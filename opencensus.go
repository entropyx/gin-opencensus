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

func Middleware(config *Config) gin.HandlerFunc {
	return gin.HandlerFunc(func(c *gin.Context) {
		request := c.Request
		spanCtx, _ := propagation.Extract(propagation.FormatTextMap, propagation.HTTPHeader(request.Header))
		_, span := trace.StartSpanWithRemoteParent(c, fmt.Sprintf("[%s] %s", request.Method, request.URL.Path), spanCtx, trace.WithSampler(setSampler()))
		defer span.End()
		c.Set("span", span)
		c.Next()
		status := c.Writer.Status()
		if status >= 400 {
			span.AddAttributes(trace.StringAttribute("error.msg", strconv.Itoa(status)))
		}
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
