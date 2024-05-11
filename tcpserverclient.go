package multirpc

type TcpServerClient struct {
	Client TcpClient
	Server TcpServer
}

func NewTcpServerClient(object interface{}, serverAddress string, clientAddress string) TcpServerClient {
	tsc := TcpServerClient{
		Client: NewTcpClient(clientAddress),
		Server: NewTcpServer(object, serverAddress),
	}
	return tsc
}
