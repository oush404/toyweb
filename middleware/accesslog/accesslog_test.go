package accesslog

import (
	"log"
	"testing"
	"toyweb"
)

func TestMiddlewareBuilder_Build(t *testing.T) {
	s := toyweb.NewHTTPServer()
	s.Get("/", func(ctx *toyweb.Context) {
		ctx.Resp.Write([]byte("hello,world"))
	})
	//s.Use((&MiddlewareBuilder{
	//	logFunc: func(accessLog string) {
	//		log.Println(accessLog)
	//	},
	//}).Build())

	s.Use(NewBuilder().LogFunc(func(accessLog string) {
		log.Println(accessLog)
	}).Build())

	s.Start(":8081")
}
