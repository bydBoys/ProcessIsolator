package server

import (
	"ProcessIsolator/constants"
	"ProcessIsolator/impl/server/internal/impl"
	"fmt"
	"net"
	"net/http"
	"net/rpc"
	"runtime"
)

var listener net.Listener

func StartServer(errorChan chan<- error, msgChan chan<- string) {
	var err error
	if listener, err = net.Listen("tcp", constants.Port); err != nil {
		errorChan <- err
		return
	}
	msgChan <- fmt.Sprintf("listen(%s) on %s", "tcp", constants.Port)
	startFileServer(errorChan, msgChan)
	go startProcessServer(errorChan, msgChan)
	runtime.Gosched()
	runtime.GC()
}

func startProcessServer(errorChan chan<- error, msgChan chan<- string) {
	runtime.LockOSThread()

	api := new(impl.ProcessServerImpl)
	newServer := rpc.NewServer()

	if err := newServer.Register(api); err != nil {
		errorChan <- err
		return
	}
	//newServer.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	http.Handle(rpc.DefaultRPCPath, newServer)
	//http.Handle(rpc.DefaultDebugPath, debugHTTP{newServer})
	msgChan <- "ProcessServer has started"
	if err := http.Serve(listener, nil); err != nil {
		errorChan <- err
		return
	}
}
func startFileServer(errorChan chan<- error, msgChan chan<- string) {
	api := new(impl.FileServerImpl)
	api.Init(errorChan, msgChan)
	http.HandleFunc("/upload", api.UploadFile)
	http.HandleFunc("/list", api.ListFile)
	http.HandleFunc("/delete", api.DeleteFile)
	msgChan <- "FileServer has started"
}
