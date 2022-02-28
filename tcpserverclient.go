package multirpc

type TcpServerClient struct {
	Client TcpClient
	Server TcpServer
}

func NewTcpServerClient(object interface{}, serverAddress string, serverName string, clientServerAddress string, clientName string) TcpServerClient {
	return TcpServerClient{
		Client: NewTcpClient(clientServerAddress, clientName),
		Server: NewTcpServer(object, serverAddress, serverName),
	}
}
