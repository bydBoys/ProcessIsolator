package server

import (
	"ProcessIsolator/constants"
	"github.com/fatih/color"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"runtime"
)

func StartRPCServer() error {
	runtime.LockOSThread()

	api := new(ProcServerImpl)
	newServer := rpc.NewServer()

	if err := newServer.Register(api); err != nil {
		return err
	}
	//newServer.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)

	http.Handle(rpc.DefaultRPCPath, newServer)
	//http.Handle(rpc.DefaultDebugPath, debugHTTP{newServer})

	color.Yellow("tcp listen: %s", constants.Port)

	listener, err := net.Listen("tcp", constants.Port)
	if err != nil {
		log.Fatal("StartRPCServer listen error:", err)
		return err
	}

	if err := http.Serve(listener, nil); err != nil {
		return err
	}
	return nil
}
