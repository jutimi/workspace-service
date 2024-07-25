package client_grpc

import (
	"log"
	"workspace-server/config"

	"github.com/jutimi/grpc-service/oauth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type oauthClient struct {
	conn        *grpc.ClientConn
	oauthClient oauth.OAuthRouteClient
	userClient  oauth.UserRouteClient
}

type OAuthClient interface{}

func NewOAuthClient() OAuthClient {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	conf := config.GetConfiguration().GRPC

	// Connect to Workspace grpc server
	conn, err := grpc.NewClient(conf.OAuthPort, opts...)
	if err != nil {
		log.Fatalf("Error connect to OAuth grpc server: %s", err.Error())
	}

	return &oauthClient{
		conn:        conn,
		oauthClient: oauth.NewOAuthRouteClient(conn),
		userClient:  oauth.NewUserRouteClient(conn),
	}
}
