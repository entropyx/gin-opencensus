package opencensus

import (
	"github.com/gin-gonic/gin"
	. "github.com/smartystreets/goconvey/convey"
	"go.opencensus.io/trace"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMiddleware(t *testing.T) {
	Convey("Given a new server", t, func() {
		var span *trace.Span
		r := gin.New()
		r.Use(Middleware(nil))
		r.GET("users/:userID/posts/:postID", func(c *gin.Context) {
			span = c.MustGet("span").(*trace.Span)
		})

		Convey("When it is requested", func() {
			req, _ := http.NewRequest("GET", "users/1/posts/2", nil)
			resp := httptest.NewRecorder()
			r.ServeHTTP(resp, req)

			Convey("The span should be valid", func() {
				So(span.SpanContext().TraceID.String(), ShouldNotBeEmpty)
				So(span.SpanContext().SpanID.String(), ShouldNotBeEmpty)
			})
		})
	})
}
