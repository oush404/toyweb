package main

import (
	"fmt"
	"net/http"
	"time"
)

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		fmt.Printf("started request for %s at %v\n", r.URL.Path, start)

		next.ServeHTTP(w, r)

		fmt.Printf("completed request for %s in %v\n", r.URL.Path, time.Since(start))
	})
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("11111")
	fmt.Fprintf(w, "hello, AOP in Go!")
	fmt.Println("2222")
}

//func main() {
//	mux := http.NewServeMux()
//	mux.Handle("/user", loggingMiddleware(http.HandlerFunc(helloHandler)))
//	http.ListenAndServe(":8082", mux)
//}
