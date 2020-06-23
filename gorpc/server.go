package main

import (
	"log"
	"net"
	"net/http"
	"net/rpc"
	"rpc/gorpc/service"
)

func main() {
	stringService := new(service.StringService)

	registerError := rpc.Register(stringService)
	if registerError != nil {
		log.Fatal("register error: ", registerError)
	}
	rpc.HandleHTTP()
	listener, err := net.Listen("tcp", "127.0.0.1:1234")
	if err != nil {
		log.Fatal("listen error:", err)
	}
	http.Serve(listener, nil)

}
