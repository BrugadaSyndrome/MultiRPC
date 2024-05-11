package multirpc

type HttpServerClient struct {
	Client HttpClient
	Server HttpServer
}

func NewHttpServerClient(object interface{}, serverAddress string, clientAddress string) HttpServerClient {
	hsc := HttpServerClient{
		Client: NewHttpClient(clientAddress),
		Server: NewHttpServer(object, serverAddress),
	}
	return hsc
}
