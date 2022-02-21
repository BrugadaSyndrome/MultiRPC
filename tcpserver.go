package rpc

import (
	"log"
	"net"
	"net/rpc"
	"os"
	"sync"
	"time"
)

type TcpServer struct {
	address  string
	listener *net.TCPListener
	object   interface{}
	shutdown chan bool

	Logger *log.Logger
	Name   string
	WG     *sync.WaitGroup
}

func NewTcpServer(object interface{}, address string, name string) TcpServer {
	return TcpServer{
		address:  address,
		object:   object,
		shutdown: make(chan bool, 1),
		Logger:   log.New(os.Stdout, name, log.Ldate|log.Ltime|log.Lmsgprefix),
		Name:     name,
		WG:       &sync.WaitGroup{},
	}
}

func (ts *TcpServer) Run() error {
	ts.WG.Add(1)

	handler := rpc.NewServer()
	err := handler.Register(ts.object)
	if err != nil {
		ts.Logger.Println("Error registering object")
		return err
	}

	tcpAddress, err := net.ResolveTCPAddr("tcp", ts.address)
	if err != nil {
		ts.Logger.Println("Error resolving tcp address %s", ts.address)
		return err
	}

	ts.listener, err = net.ListenTCP("tcp", tcpAddress)
	if err != nil {
		ts.Logger.Println("Error listening at address %s", ts.address)
		return err
	}

	go func() {
		for {
			select {
			case <-ts.shutdown:
				// Server has been given the signal to shut down
				err := ts.listener.Close()
				if err != nil {
					ts.Logger.Println("Server closed connection to client - %s", err)
				}
				return
			default:
				// Poll this connection periodically
				err := ts.listener.SetDeadline(time.Now().Add(1 * time.Second))
				if err != nil {
					ts.Logger.Fatal(err.Error())
				}
			}

			conn, err := ts.listener.Accept()
			if err != nil {
				netErr, ok := err.(net.Error)
				if ok && netErr.Timeout() {
					// Deadline timeout has occurred
					continue
				}
				// There was an error listening
				ts.Logger.Printf("Error listening on connection at address %s - %s\n", conn.RemoteAddr(), err.Error())
				continue
			}

			ts.Logger.Println("Server opened connection to client at address %s", conn.RemoteAddr())
			go handler.ServeConn(conn)
		}
	}()

	ts.Logger.Println("Running server at address %s", ts.address)
	return nil
}

func (ts *TcpServer) Stop() error {
	ts.Logger.Println("Shutting down server at address %s", ts.address)
	close(ts.shutdown)
	ts.WG.Done()
	return nil
}
