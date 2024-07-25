package client_grpc

type ClientGRPCCollection struct {
	OAuthClient OAuthClient
}

func RegisterClientGRPC() ClientGRPCCollection {
	return ClientGRPCCollection{
		OAuthClient: NewOAuthClient(),
	}
}
