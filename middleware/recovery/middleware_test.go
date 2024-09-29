package recovery

import (
	"log"
	"testing"
	"toyweb"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	s := toyweb.NewHTTPServer()
	s.Get("/user", func(ctx *toyweb.Context) {
		ctx.RespData = []byte("hello, world")
	})

	s.Get("/", func(ctx *toyweb.Context) {
		panic(" panic了")
	})

	s.Use((&MiddlewareBuilder{
		StatusCode: 500,
		ErrMsg:     "你 Panic 了",
		LogFunc: func(ctx *toyweb.Context) {
			log.Println(ctx.Req.URL.Path)
		},
	}).Build())

	s.Start(":8081")
}
