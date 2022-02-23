package rpc

import (
	"errors"
	"fmt"
	"log"
	"net/rpc"
	"os"
)

type HttpClient struct {
	serverAddress string
	client        *rpc.Client

	Logger *log.Logger
	Name   string
}

// NewHttpClient will return a new HttpClient object
func NewHttpClient(serverAddress string, name string) HttpClient {
	return HttpClient{
		serverAddress: serverAddress,
		Logger:        log.New(os.Stdout, name, log.Ldate|log.Ltime|log.Lmsgprefix),
		Name:          name,
	}
}

// Connect will attempt to connect this RPC client to the RPC server specified when this object was created
func (hc *HttpClient) Connect() error {
	// Check to make sure that HttpClient.Connect has not already been called
	if hc.client != nil {
		message := fmt.Sprintf("Already connected to server at address %s", hc.serverAddress)
		hc.Logger.Println(message)
		return errors.New(message)
	}

	// Make the initial connection to the RPC server
	var err error
	hc.client, err = rpc.DialHTTP("tcp", hc.serverAddress)
	if err != nil {
		hc.Logger.Println("Error connecting to server at address %s : %s", hc.serverAddress, err)
		return err
	}
	hc.Logger.Println("Connected to server at %s", hc.serverAddress)
	return nil
}

// Call will execute the specified method on the RPC server
func (hc *HttpClient) Call(method string, request interface{}, reply interface{}) error {
	// Check to make sure that HttpClient.Connect has already been called
	if hc.client == nil {
		message := fmt.Sprintf("Not connected to server at address: %s, method: %s", hc.serverAddress, method)
		hc.Logger.Println(message)
		return errors.New(message)
	}

	// Make the call to the RPC server with the specified method and associated data
	err := hc.client.Call(method, request, reply)
	if err != nil {
		hc.Logger.Println("Error calling server at address %s : method %s", hc.serverAddress, method)
		return err
	}
	hc.Logger.Println("Calling server %s", method)
	return nil
}

// Disconnect will close the HTTP connection to the RPC Server
func (hc *HttpClient) Disconnect() error {
	// Check to make sure that HttpClient.Connect has already been called
	if hc.client == nil {
		message := fmt.Sprintf("Already disconnected from server at address %s", hc.serverAddress)
		hc.Logger.Println(message)
		return errors.New(message)
	}

	// Close the connection to the RPC server
	err := hc.client.Close()
	if err != nil {
		hc.Logger.Println("Error disconnecting from server at address %s", hc.serverAddress)
		return err
	}
	hc.Logger.Println("Disconnected from server at %s", hc.serverAddress)
	return nil
}
