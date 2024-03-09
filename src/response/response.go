package response

import (
	"encoding/json"
	"log"

	"github.com/valyala/fasthttp"
)

// RJSON - Function that returns a JSON response
func RJSON(ctx *fasthttp.RequestCtx, statusCode int, data interface{}) {
	ctx.SetContentType("application/json; charset=utf-8")
	ctx.SetStatusCode(statusCode)

	if data != nil {
		if err := json.NewEncoder(ctx).Encode(data); err != nil {
			log.Fatal(err)
		}
	}
}

// RError - Function that returns an error response
func RError(ctx *fasthttp.RequestCtx, statusCode int, err string) {
	RJSON(ctx, statusCode, struct {
		Error string `json:"error"`
	}{
		Error: err,
	})
}
