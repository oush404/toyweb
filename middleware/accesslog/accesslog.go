package accesslog

import (
	"encoding/json"
	"log"
	"toyweb"
)

type LogHandlerFunc func(accessLog string)

type MiddlewareBuilder struct {
	logFunc LogHandlerFunc
}

func (b *MiddlewareBuilder) LogFunc(logFunc LogHandlerFunc) *MiddlewareBuilder {
	b.logFunc = logFunc
	return b
}

func NewBuilder() *MiddlewareBuilder {
	return &MiddlewareBuilder{
		logFunc: func(accessLog string) {
			log.Println(accessLog)
		},
	}
}

type accessLog struct {
	Host       string
	Route      string
	HTTPMethod string `json:"http_method"`
	Path       string
}

func (b MiddlewareBuilder) Build() toyweb.Middleware {
	return func(next toyweb.HandleFunc) toyweb.HandleFunc {
		return func(ctx *toyweb.Context) {
			defer func() {
				l := accessLog{
					Host:       ctx.Req.Host,
					Route:      ctx.MatchedRoute,
					Path:       ctx.Req.URL.Path,
					HTTPMethod: ctx.Req.Method,
				}
				val, _ := json.Marshal(l)
				b.logFunc(string(val))
			}()
			next(ctx)
		}
	}
}
