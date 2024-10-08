package client

import (
	"context"
	"log"

	"workspace-server/config"
	"workspace-server/package/errors"

	"github.com/jutimi/grpc-service/oauth"
	"github.com/jutimi/grpc-service/utils"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type oauthClient struct {
	conn       *grpc.ClientConn
	userClient oauth.UserRouteClient
}

type OAuthClient interface {
	CloseConn()
	GetUserByFilter(ctx context.Context, data *oauth.GetUserByFilterParams) (*oauth.UserResponse, error)
	GetUsersByFilter(ctx context.Context, data *oauth.GetUserByFilterParams) (*oauth.UsersResponse, error)
	CreateUser(ctx context.Context, data *oauth.CreateUserParams) (*oauth.UserResponse, error)
	BulkCreateUsers(ctx context.Context, data *oauth.CreateUsersParams) (*oauth.UsersResponse, error)
}

func NewOAuthClient() OAuthClient {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	}
	conf := config.GetConfiguration().GRPC

	// Connect to Workspace grpc server
	conn, err := grpc.NewClient(conf.OAuthUrl, opts...)
	if err != nil {
		log.Fatalf("Error connect to OAuth grpc server: %s", err.Error())
	}

	return &oauthClient{
		conn:       conn,
		userClient: oauth.NewUserRouteClient(conn),
	}
}

func (c *oauthClient) GetUserByFilter(ctx context.Context, data *oauth.GetUserByFilterParams) (*oauth.UserResponse, error) {
	conf := config.GetConfiguration().Server
	ctx = utils.GenerateContext(utils.Metadata{
		Ctx:         ctx,
		ServiceName: conf.ServiceName,
	})

	resp, err := c.userClient.GetUserByFilter(ctx, data)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	if resp.Error != nil {
		return nil, errors.NewCustomError(int(resp.Error.ErrorCode), resp.Error.ErrorMessage)
	}
	if !resp.Success {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return resp, nil
}

func (c *oauthClient) GetUsersByFilter(ctx context.Context, data *oauth.GetUserByFilterParams) (*oauth.UsersResponse, error) {
	conf := config.GetConfiguration().Server
	ctx = utils.GenerateContext(utils.Metadata{
		Ctx:         ctx,
		ServiceName: conf.ServiceName,
	})

	resp, err := c.userClient.GetUsersByFilter(ctx, data)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	if resp.Error != nil {
		return nil, errors.NewCustomError(int(resp.Error.ErrorCode), resp.Error.ErrorMessage)
	}
	if !resp.Success {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return resp, nil
}

func (c *oauthClient) CreateUser(ctx context.Context, data *oauth.CreateUserParams) (*oauth.UserResponse, error) {
	conf := config.GetConfiguration().Server
	ctx = utils.GenerateContext(utils.Metadata{
		Ctx:         ctx,
		ServiceName: conf.ServiceName,
	})

	resp, err := c.userClient.CreateUser(ctx, data)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	if resp.Error != nil {
		return nil, errors.NewCustomError(int(resp.Error.ErrorCode), resp.Error.ErrorMessage)
	}
	if !resp.Success {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return resp, nil
}

func (c *oauthClient) BulkCreateUsers(ctx context.Context, data *oauth.CreateUsersParams) (*oauth.UsersResponse, error) {
	conf := config.GetConfiguration().Server
	ctx = utils.GenerateContext(utils.Metadata{
		Ctx:         ctx,
		ServiceName: conf.ServiceName,
	})

	resp, err := c.userClient.BulkCreateUsers(ctx, data)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}
	if resp.Error != nil {
		return nil, errors.NewCustomError(int(resp.Error.ErrorCode), resp.Error.ErrorMessage)
	}
	if !resp.Success {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return resp, nil
}

func (c *oauthClient) CloseConn() {
	c.conn.Close()
}
