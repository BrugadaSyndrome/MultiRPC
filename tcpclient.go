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

// NewTcpClient will return a new TcpClient object
func NewTcpClient(serverAddress string, name string) TcpClient {
	return TcpClient{
		serverAddress: serverAddress,
		Name:          name,
		Logger:        log.New(os.Stdout, name, log.Ldate|log.Ltime|log.Lmsgprefix),
	}
}

// Connect will attempt to connect this RPC client to the RPC server specified when this object was created
func (tc *TcpClient) Connect() error {
	// Check to make sure that TcpClient.Connect has not already been called
	if tc.client != nil {
		message := fmt.Sprintf("Already connected to server at address %s", tc.serverAddress)
		tc.Logger.Println(message)
		return errors.New(message)
	}

	// Make the initial connection to the RPC server
	var err error
	tc.client, err = rpc.Dial("tcp", tc.serverAddress)
	if err != nil {
		tc.Logger.Println("Error connecting to server at address %s", tc.serverAddress)
		return err
	}
	tc.Logger.Println("Connected to server at: %s", tc.serverAddress)
	return nil
}

// Call will execute the specified method on the RPC server
func (tc *TcpClient) Call(method string, request interface{}, reply interface{}) error {
	// Check to make sure that TcpClient.Connect has already been called
	if tc.client == nil {
		message := fmt.Sprintf("Not connected to server at address %s : method %s", tc.serverAddress, method)
		tc.Logger.Println(message)
		return errors.New(message)
	}

	// Make the call to the RPC server with the specified method and associated data
	err := tc.client.Call(method, request, reply)
	if err != nil {
		tc.Logger.Println("Error calling server at address: %s, method: %s", tc.serverAddress, method)
		return err
	}
	tc.Logger.Println("Calling server [%s] %s", tc.serverAddress, method)
	return nil
}

// Disconnect will close the TCP connection to the RPC Server
func (tc *TcpClient) Disconnect() error {
	// Check to make sure that TcpClient.Connect has already been called
	if tc.client == nil {
		message := fmt.Sprintf("Already disconnected from server at address %s", tc.serverAddress)
		tc.Logger.Println(message)
		return errors.New(message)
	}

	// Close the connection to the RPC server
	err := tc.client.Close()
	if err != nil {
		tc.Logger.Println("Error disconnecting from server at serverAddress %s", tc.serverAddress)
		return err
	}
	tc.Logger.Println("Disconnected from server at %s", tc.serverAddress)
	return nil
}
