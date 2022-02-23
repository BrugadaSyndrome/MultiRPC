package rpc

import (
	"context"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"sync"
)

type HttpServer struct {
	address  string
	listener net.Listener
	mux      *http.ServeMux
	object   interface{}
	server   *http.Server
	wg       *sync.WaitGroup

	Logger *log.Logger
	Name   string
}

// NewHttpServer will return a new HttpServer object
func NewHttpServer(object interface{}, address string, name string) HttpServer {
	return HttpServer{
		address: address,
		mux:     http.NewServeMux(),
		object:  object,
		wg:      &sync.WaitGroup{},
		Logger:  log.New(os.Stdout, name, log.Ldate|log.Ltime|log.Lmsgprefix),
		Name:    name,
	}
}

// Run will start serving the object via RPC over HTTP
func (hs *HttpServer) Run() error {
	handler := rpc.NewServer()
	err := handler.Register(hs.object)
	if err != nil {
		hs.Logger.Println("Error registering object")
		return err
	}

	// Make a new http request multiplexer for this object
	// https://github.com/golang/go/issues/13395
	oldMux := http.DefaultServeMux
	http.DefaultServeMux = hs.mux
	handler.HandleHTTP(rpc.DefaultRPCPath, rpc.DefaultDebugPath)
	http.DefaultServeMux = oldMux

	// Make a new listener for this object
	hs.listener, err = net.Listen("tcp", hs.address)
	if err != nil {
		hs.Logger.Println("Error listening at address %s", hs.address)
		return err
	}

	// Start the server until a stop signal is received
	hs.server = &http.Server{Addr: hs.address, Handler: hs.mux}
	go func() {
		hs.wg.Add(1)

		if err := hs.server.Serve(hs.listener); err != http.ErrServerClosed {
			hs.Logger.Println("Error serving at address %s", hs.address)
			hs.Logger.Fatal(err.Error())
		}
	}()

	hs.Logger.Println("Running server at address %s", hs.address)
	return nil
}

// Stop is called to shut down the server by decrementing the wait group
func (hs *HttpServer) Stop() error {
	if err := hs.server.Shutdown(context.Background()); err != nil {
		hs.Logger.Println("Error shutting down server at address %s", hs.address)
		return err
	}
	hs.Logger.Println("Shutting down server at address %s", hs.address)
	hs.wg.Done()
	return nil
}

// Wait can be called to have the code wait for the server to shut down before continuing
func (hs *HttpServer) Wait() {
	hs.wg.Wait()
}
