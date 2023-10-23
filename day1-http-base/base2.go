package main

import (
	"fmt"
	"log"
	"net/http"
)

type Engine struct {
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		fmt.Fprintf(w, "URL.Path = %q\n", req.URL.Path)
	case "/test":
		for i, v := range req.Header {
			fmt.Fprintf(w, "Header[%q]=%q\n", i, v)
		}
	default:
		fmt.Fprintf(w, "404 NOT FOUND:%s\n", req.URL)
	}
}

// 第二个参数
//
//	type Handler interface {
//	   ServeHTTP(ResponseWriter, *Request)
//	}
func main() {
	engine := new(Engine)
	log.Fatal(http.ListenAndServe(":9088", engine))
}
