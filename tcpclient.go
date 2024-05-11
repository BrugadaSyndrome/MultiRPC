package multirpc

import (
	"errors"
	"fmt"
	"github.com/BrugadaSyndrome/bslogger"
	"net/rpc"
)

type TcpClient struct {
	client        *rpc.Client
	logger        bslogger.Logger
	serverAddress string
}

// NewTcpClient will return a new TcpClient object
func NewTcpClient(serverAddress string) TcpClient {
	tc := TcpClient{
		logger:        bslogger.NewLogger(),
		serverAddress: serverAddress,
	}
	tc.logger.Name = fmt.Sprintf("[Tcpclient %s]", serverAddress)
	return tc
}

// Connect will attempt to connect this RPC client to the RPC server specified when this object was created
func (tc *TcpClient) Connect() error {
	// Check to make sure that TcpClient.Connect has not already been called
	if tc.client != nil {
		message := fmt.Sprintf("Already connected to server at address %s", tc.serverAddress)
		tc.logger.Errorf(message)
		return errors.New(message)
	}

	// Make the initial connection to the RPC server
	var err error
	tc.client, err = rpc.Dial("tcp", tc.serverAddress)
	if err != nil {
		tc.logger.Errorf("Connecting to server at address %s", tc.serverAddress)
		return err
	}
	tc.logger.Infof("Connected to server at: %s", tc.serverAddress)
	return nil
}

// Call will execute the specified method on the RPC server
func (tc *TcpClient) Call(method string, request interface{}, reply interface{}) error {
	// Check to make sure that TcpClient.Connect has already been called
	if tc.client == nil {
		message := fmt.Sprintf("Not connected to server at address %s : method %s", tc.serverAddress, method)
		tc.logger.Error(message)
		return errors.New(message)
	}

	// Make the call to the RPC server with the specified method and associated data
	err := tc.client.Call(method, request, reply)
	if err != nil {
		tc.logger.Errorf("Calling server at address: %s, method: %s", tc.serverAddress, method)
		return err
	}
	tc.logger.Debugf("Calling server [%s] %s", tc.serverAddress, method)
	return nil
}

// Disconnect will close the TCP connection to the RPC Server
func (tc *TcpClient) Disconnect() error {
	// Check to make sure that TcpClient.Connect has already been called
	if tc.client == nil {
		message := fmt.Sprintf("Already disconnected from server at address %s", tc.serverAddress)
		tc.logger.Error(message)
		return errors.New(message)
	}

	// Close the connection to the RPC server
	err := tc.client.Close()
	if err != nil {
		tc.logger.Errorf("Disconnecting from server at serverAddress %s", tc.serverAddress)
		return err
	}
	tc.logger.Infof("Disconnected from server at %s", tc.serverAddress)
	return nil
}

// LoggerVerbosity exposes the logger.Verbosity field
func (tc *TcpClient) LoggerVerbosity(verbosity bslogger.Verbosity) {
	tc.logger.Verbosity = verbosity
}

// LoggerName exposes the logger.Name field
func (tc *TcpClient) LoggerName(name string) {
	tc.logger.Name = name
}
