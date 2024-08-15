package server

import (
	"context"

	"workspace-server/app/entity"
	"workspace-server/app/helper"
	"workspace-server/app/repository"
	postgres_repository "workspace-server/app/repository/postgres"
	"workspace-server/package/errors"
	logger "workspace-server/package/log"
	"workspace-server/utils"

	"github.com/google/uuid"
	grpc_utils "github.com/jutimi/grpc-service/utils"
	grpc_workspace "github.com/jutimi/grpc-service/workspace"
	"google.golang.org/grpc/metadata"
)

type grpcServer struct {
	grpc_workspace.UnimplementedWorkspaceRouteServer
	grpc_workspace.UnimplementedUserWorkspaceRouteServer

	postgresRepo postgres_repository.PostgresRepositoryCollections
	helper       helper.HelperCollections
}

type WorkspaceServer interface {
	grpc_workspace.WorkspaceRouteServer
	grpc_workspace.UserWorkspaceRouteServer
}

func NewGRPCServer(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
	helper helper.HelperCollections,
) WorkspaceServer {
	return &grpcServer{
		postgresRepo: postgresRepo,
		helper:       helper,
	}
}

func (s *grpcServer) GetWorkspaceByFilter(ctx context.Context, data *grpc_workspace.GetWorkspaceByFilterParams) (*grpc_workspace.WorkspaceResponse, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		logger.Println(logger.LogPrintln{
			Ctx:       ctx,
			FileName:  "external/server/server.go",
			FuncName:  "GetWorkspaceByFilter",
			TraceData: utils.ConvertStructToString(data),
			Msg:       "InComing Metadata: " + utils.ConvertStructToString(md),
		})
	}

	customErr := errors.New(errors.ErrCodeInternalServerError)
	ids := make([]uuid.UUID, 0)

	id, err := utils.ConvertStringToUUID(*data.Id)
	if err != nil {
		return &grpc_workspace.WorkspaceResponse{
			Success: false,
			Data:    nil,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}
	for _, id := range data.Ids {
		convertId, err := utils.ConvertStringToUUID(id)
		if err != nil {
			return &grpc_workspace.WorkspaceResponse{
				Success: false,
				Data:    nil,
				Error: grpc_utils.FormatErrorResponse(
					int32(customErr.GetCode()),
					customErr.Error(),
				),
			}, nil
		}
		ids = append(ids, convertId)
	}

	workspace, err := s.postgresRepo.WorkspaceRepo.FindOneByFilter(ctx, nil, &repository.FindWorkspaceByFilter{
		Id:       &id,
		IsActive: data.IsActive,
		Ids:      ids,
	})
	if err != nil {
		customErr = errors.New(errors.ErrCodeWorkspaceNotFound)
		return &grpc_workspace.WorkspaceResponse{
			Success: false,
			Data:    nil,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	return &grpc_workspace.WorkspaceResponse{
		Success: true,
		Data: &grpc_workspace.WorkspaceDetail{
			Id:   workspace.Id.String(),
			Name: workspace.Name,
		},
		Error: nil,
	}, nil
}

func (s *grpcServer) GetUserWorkspaceByFilter(ctx context.Context, data *grpc_workspace.GetUserWorkspaceByFilterParams) (*grpc_workspace.UserWorkspaceResponse, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		logger.Println(logger.LogPrintln{
			Ctx:       ctx,
			FileName:  "external/server/server.go",
			FuncName:  "GetUserWorkspaceByFilter",
			TraceData: utils.ConvertStructToString(data),
			Msg:       "InComing Metadata: " + utils.ConvertStructToString(md),
		})
	}

	var email, phoneNumber string
	customErr := errors.New(errors.ErrCodeInternalServerError)
	filter, err := convertUserParamsToFilter(data)
	if err != nil {
		return &grpc_workspace.UserWorkspaceResponse{
			Success: false,
			Data:    nil,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	userWorkspace, err := s.postgresRepo.UserWorkspaceRepo.FindOneByFilter(ctx, nil, filter)
	if err != nil {
		customErr = errors.New(errors.ErrCodeUserWorkspaceNotFound)
		return &grpc_workspace.UserWorkspaceResponse{
			Success: false,
			Data:    nil,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	if userWorkspace.Email != nil {
		email = *userWorkspace.Email
	}
	if userWorkspace.PhoneNumber != nil {
		phoneNumber = *userWorkspace.PhoneNumber
	}
	return &grpc_workspace.UserWorkspaceResponse{
		Success: true,
		Data: &grpc_workspace.UserWorkspaceDetail{
			Id:          userWorkspace.Id.String(),
			Name:        userWorkspace.UserWorkspaceDetail.Name,
			Email:       email,
			PhoneNumber: phoneNumber,
			IsActive:    userWorkspace.IsActive,
			Role:        convertUserWorkspaceRole(userWorkspace.Role),
			WorkspaceId: userWorkspace.WorkspaceId.String(),
			UserId:      userWorkspace.UserId.String(),
		},
		Error: nil,
	}, nil
}

func (s *grpcServer) GetUserWorkspacesByFilter(ctx context.Context, data *grpc_workspace.GetUserWorkspaceByFilterParams) (*grpc_workspace.UserWorkspacesResponse, error) {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		logger.Println(logger.LogPrintln{
			Ctx:       ctx,
			FileName:  "external/server/server.go",
			FuncName:  "GetUserWorkspacesByFilter",
			TraceData: utils.ConvertStructToString(data),
			Msg:       "InComing Metadata: " + utils.ConvertStructToString(md),
		})
	}

	var email, phoneNumber string
	var userWorkspaceRes []*grpc_workspace.UserWorkspaceDetail
	customErr := errors.New(errors.ErrCodeInternalServerError)

	filter, err := convertUserParamsToFilter(data)
	if err != nil {
		return &grpc_workspace.UserWorkspacesResponse{
			Success: false,
			Data:    nil,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	userWorkspace, err := s.postgresRepo.UserWorkspaceRepo.FindByFilter(ctx, nil, filter)
	if err != nil {
		customErr = errors.New(errors.ErrCodeUserWorkspaceNotFound)
		return &grpc_workspace.UserWorkspacesResponse{
			Success: false,
			Data:    nil,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	for _, user := range userWorkspace {
		if user.Email != nil {
			email = *user.Email
		}
		if user.PhoneNumber != nil {
			phoneNumber = *user.PhoneNumber
		}

		userWorkspaceRes = append(userWorkspaceRes, &grpc_workspace.UserWorkspaceDetail{
			Id:          user.Id.String(),
			Name:        user.UserWorkspaceDetail.Name,
			Email:       email,
			PhoneNumber: phoneNumber,
			IsActive:    user.IsActive,
			Role:        convertUserWorkspaceRole(user.Role),
			WorkspaceId: user.WorkspaceId.String(),
			UserId:      user.UserId.String(),
		})
	}

	return &grpc_workspace.UserWorkspacesResponse{
		Success: true,
		Data:    userWorkspaceRes,
		Error:   nil,
	}, nil
}

// ------------------------ Helper ------------------------
func convertUserParamsToFilter(data *grpc_workspace.GetUserWorkspaceByFilterParams) (*repository.FindUserWorkspaceByFilter, error) {
	var id, workspaceId, userId uuid.UUID
	var ids, workspaceIds, userIds []uuid.UUID
	var err error
	var limit, offset int

	if data.Limit != nil {
		limit = int(*data.Limit)
	}
	if data.Offset != nil {
		offset = int(*data.Offset)
	}

	if data.Id != nil && *data.Id != "" {
		id, err = utils.ConvertStringToUUID(*data.Id)
		if err != nil {
			return nil, err
		}
	}
	if data.Ids != nil {
		for _, id := range data.Ids {
			convertId, err := utils.ConvertStringToUUID(id)
			if err != nil {
				return nil, err
			}

			ids = append(ids, convertId)
		}
	}

	if data.WorkspaceId != nil && *data.WorkspaceId != "" {
		workspaceId, err = utils.ConvertStringToUUID(*data.WorkspaceId)
		if err != nil {
			return nil, err
		}
	}
	if data.WorkspaceIds != nil {
		for _, id := range data.WorkspaceIds {
			convertId, err := utils.ConvertStringToUUID(id)
			if err != nil {
				return nil, err
			}

			workspaceIds = append(workspaceIds, convertId)
		}
	}

	if data.UserId != nil && *data.UserId != "" {
		userId, err = utils.ConvertStringToUUID(*data.UserId)
		if err != nil {
			return nil, err
		}
	}
	if data.UserIds != nil {
		for _, id := range data.UserIds {
			convertId, err := utils.ConvertStringToUUID(id)
			if err != nil {
				return nil, err
			}

			userIds = append(userIds, convertId)
		}
	}

	return &repository.FindUserWorkspaceByFilter{
		Id:              &id,
		Ids:             ids,
		WorkspaceId:     &workspaceId,
		WorkspaceIds:    workspaceIds,
		UserId:          &userId,
		UserIds:         userIds,
		IsActive:        data.IsActive,
		Limit:           &limit,
		Offset:          &offset,
		IsIncludeDetail: true,
	}, nil
}

func convertUserWorkspaceRole(role string) grpc_workspace.UserWorkspaceRole {
	switch role {
	case entity.ROLE_OWNER:
		return grpc_workspace.UserWorkspaceRole_OWNER
	case entity.ROLE_ADMIN:
		return grpc_workspace.UserWorkspaceRole_ADMIN
	default:
		return grpc_workspace.UserWorkspaceRole_USER
	}
}
