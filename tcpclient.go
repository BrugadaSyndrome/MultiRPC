package rpc

import (
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"os"
)

type TcpClient struct {
	client        *rpc.Client
	serverAddress string

	Logger *log.Logger
	Name   string
}

func NewTcpClient(serverAddress string, name string) TcpClient {
	return TcpClient{
		serverAddress: serverAddress,
		Name:          name,
		Logger:        log.New(os.Stdout, name, log.Ldate|log.Ltime|log.Lmsgprefix),
	}
}

func (tc *TcpClient) Connect() error {
	if tc.client != nil {
		message := fmt.Sprintf("Already connected to server at address %s", tc.serverAddress)
		tc.Logger.Println(message)
		return nil
	}

	var err error
	tc.client, err = rpc.Dial("tcp", tc.serverAddress)
	if err != nil {
		tc.Logger.Println("Error connecting to server at address %s", tc.serverAddress)
		return err
	}
	tc.Logger.Println("Connected to server at: %s", tc.serverAddress)
	return nil
}

func (tc *TcpClient) Call(method string, request interface{}, reply interface{}) error {
	if tc.client == nil {
		message := fmt.Sprintf("Not connected to server at address %s : method %s", tc.serverAddress, method)
		tc.Logger.Println(message)
		return errors.New(message)
	}

	err := tc.client.Call(method, request, reply)
	if err != nil {
		tc.Logger.Println("Error calling server at address: %s, method: %s", tc.serverAddress, method)
		return err
	}
	tc.Logger.Println("Calling server [%s] %s", tc.serverAddress, method)
	return nil
}

func (tc *TcpClient) Disconnect() error {
	if tc.client == nil {
		message := fmt.Sprintf("Already disconnected from server at address %s", tc.serverAddress)
		tc.Logger.Println(message)
		return errors.New(message)
	}

	err := tc.client.Close()
	if err != nil {
		tc.Logger.Println("Error disconnecting from server at serverAddress %s", tc.serverAddress)
		return err
	}
	tc.Logger.Println("Disconnected from server at %s", tc.serverAddress)
	return nil
}
