package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", gee0Handler)
	http.HandleFunc("/test", gee1Handler)

	//通过调用log.Fatal()函数，程序会在启动失败时打印错误信息并退出
	log.Fatal(http.ListenAndServe(":9080", nil))
}

// 这是名为indexHandler的处理器函数。它接受两个参数：
// w是一个http.ResponseWriter类型的对象，用于将响应发送给客户端；
// req是一个http.Request类型的对象，包含了从客户端发出的请求信息。
func gee0Handler(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "URL.path= %q\n", req.URL.Path)
}

func gee1Handler(w http.ResponseWriter, req *http.Request) {
	for i, v := range req.Header {
		fmt.Fprintf(w, "Header[%q]=%q\n", i, v)
	}
}
