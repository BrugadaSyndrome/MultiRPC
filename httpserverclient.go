package multirpc

type HttpServerClient struct {
	Client HttpClient
	Server HttpServer
}

func NewHttpServerClient(object interface{}, serverAddress string, serverName string, clientServerAddress string, clientName string) HttpServerClient {
	return HttpServerClient{
		Client: NewHttpClient(clientServerAddress, clientName),
		Server: NewHttpServer(object, serverAddress, serverName),
	}
}
