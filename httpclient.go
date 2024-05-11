package multirpc

import (
	"errors"
	"fmt"
	"github.com/BrugadaSyndrome/bslogger"
	"net/rpc"
)

type HttpClient struct {
	client        *rpc.Client
	logger        bslogger.Logger
	serverAddress string
}

// NewHttpClient will return a new HttpClient object
func NewHttpClient(serverAddress string) HttpClient {
	hc := HttpClient{
		logger:        bslogger.NewLogger(),
		serverAddress: serverAddress,
	}
	hc.logger.Name = fmt.Sprintf("[HttpClient %s]", serverAddress)
	return hc
}

// Connect will attempt to connect this RPC client to the RPC server specified when this object was created
func (hc *HttpClient) Connect() error {
	// Check to make sure that HttpClient.Connect has not already been called
	if hc.client != nil {
		message := fmt.Sprintf("Already connected to server at address %s", hc.serverAddress)
		hc.logger.Error(message)
		return errors.New(message)
	}

	// Make the initial connection to the RPC server
	var err error
	hc.client, err = rpc.DialHTTP("tcp", hc.serverAddress)
	if err != nil {
		hc.logger.Errorf("Connecting to server at address %s : %s", hc.serverAddress, err)
		return err
	}
	hc.logger.Infof("Connected to server at %s", hc.serverAddress)
	return nil
}

// Call will execute the specified method on the RPC server
func (hc *HttpClient) Call(method string, request interface{}, reply interface{}) error {
	// Check to make sure that HttpClient.Connect has already been called
	if hc.client == nil {
		message := fmt.Sprintf("Not connected to server at address: %s, method: %s", hc.serverAddress, method)
		hc.logger.Error(message)
		return errors.New(message)
	}

	// Make the call to the RPC server with the specified method and associated data
	err := hc.client.Call(method, request, reply)
	if err != nil {
		hc.logger.Errorf("Calling server at address %s : method %s", hc.serverAddress, method)
		return err
	}
	hc.logger.Debugf("Calling server [%s] %s", hc.serverAddress, method)
	return nil
}

// Disconnect will close the HTTP connection to the RPC Server
func (hc *HttpClient) Disconnect() error {
	// Check to make sure that HttpClient.Connect has already been called
	if hc.client == nil {
		message := fmt.Sprintf("Already disconnected from server at address %s", hc.serverAddress)
		hc.logger.Error(message)
		return errors.New(message)
	}

	// Close the connection to the RPC server
	err := hc.client.Close()
	if err != nil {
		hc.logger.Errorf("Disconnecting from server at address %s", hc.serverAddress)
		return err
	}
	hc.logger.Infof("Disconnected from server at %s", hc.serverAddress)
	return nil
}

// LoggerVerbosity exposes the logger.Verbosity field
func (hc *HttpClient) LoggerVerbosity(verbosity bslogger.Verbosity) {
	hc.logger.Verbosity = verbosity
}

// LoggerName exposes the logger.Name field
func (hc *HttpClient) LoggerName(name string) {
	hc.logger.Name = name
}
