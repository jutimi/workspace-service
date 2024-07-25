package server_grpc

import (
	"workspace-server/app/helper"
	postgres_repository "workspace-server/app/repository/postgres"

	"github.com/jutimi/grpc-service/workspace"
)

type grpcServer struct {
	workspace.UnimplementedWorkspaceRouteServer
	workspace.UnimplementedUserWorkspaceRouteServer

	postgresRepo postgres_repository.PostgresRepositoryCollections
	helper       helper.HelperCollections
}

type WSServer interface {
	workspace.WorkspaceRouteServer
	workspace.UserWorkspaceRouteServer
}

func NewGRPCServer(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
	helper helper.HelperCollections,
) WSServer {
	return &grpcServer{
		postgresRepo: postgresRepo,
		helper:       helper,
	}
}
