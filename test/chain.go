package main

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

// 定义Handler类型
type HandlerFunc func(http.ResponseWriter, *http.Request, http.HandlerFunc)

// 中间件责任链结构
type MiddlewareChain struct {
	middlewares []HandlerFunc
	final       http.HandlerFunc
}

// Use 用于添加中间件的函数
func (mc *MiddlewareChain) Use(mw HandlerFunc) {
	mc.middlewares = append(mc.middlewares, mw)
}

// Build 构建最终处理请求的方法
func (mc *MiddlewareChain) Build() http.HandlerFunc {
	chain := mc.final
	for i := len(mc.middlewares) - 1; i >= 0; i-- {
		fn := runtime.FuncForPC(reflect.ValueOf(mc.middlewares[i]).Pointer())
		fmt.Printf("i: %d, %s\n", i, fn.Name())
		chain = buildNext(mc.middlewares[i], chain)
	}
	return chain
}

// buildNext 将每个中间件与下一个中间件连接
func buildNext(middleware HandlerFunc, next http.HandlerFunc) http.HandlerFunc {
	// 压入方法栈，栈特点：先进后出。 这里先压入authMiddleware, 在压入loggingMiddleware2
	// 所以先执行loggingMiddleware2, 在执行authMiddleware
	return func(w http.ResponseWriter, r *http.Request) {
		fn1 := runtime.FuncForPC(reflect.ValueOf(middleware).Pointer())
		fn2 := runtime.FuncForPC(reflect.ValueOf(next).Pointer())
		fmt.Printf("midFunc: %s, nextfunc: %s", fn1.Name(), fn2.Name())
		middleware(w, r, next)
	}
}

// loggingMiddleware 示例中间件1：日志中间件
func loggingMiddleware2(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("Logging request: ", r.URL.Path)
	next(w, r)
}

// authMiddleware 示例中间件2：认证中间件
func authMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	fmt.Println("checking authentication")
	next(w, r)
}

// finalHandler 最终处理请求的handler
func finalHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "final handler reached")
}

func main() {
	// 构建中间件链
	mc := &MiddlewareChain{
		final: finalHandler,
	}

	// 注册中间件
	mc.Use(loggingMiddleware2)
	mc.Use(authMiddleware)

	// 启动服务器并使用责任链
	//http.HandleFunc 将特定的 URL 路径与处理该路径的函数关联起来
	http.HandleFunc("/", mc.Build())
	fmt.Println("server start at :8080")
	http.ListenAndServe(":8080", nil)
}
