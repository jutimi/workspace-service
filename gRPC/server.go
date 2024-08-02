package server_grpc

import (
	"context"

	"github.com/jutimi/workspace-server/app/entity"
	"github.com/jutimi/workspace-server/app/helper"
	"github.com/jutimi/workspace-server/app/repository"
	postgres_repository "github.com/jutimi/workspace-server/app/repository/postgres"
	"github.com/jutimi/workspace-server/package/errors"
	"github.com/jutimi/workspace-server/utils"

	"github.com/google/uuid"
	grpc_utils "github.com/jutimi/grpc-service/utils"
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

func (s *grpcServer) GetWorkspaceByFilter(ctx context.Context, data *workspace.GetWorkspaceByFilterParams) (*workspace.WorkspaceResponse, error) {
	customErr := errors.New(errors.ErrCodeInternalServerError)
	ids := make([]uuid.UUID, 0)

	id, err := utils.ConvertStringToUUID(*data.Id)
	if err != nil {
		return &workspace.WorkspaceResponse{
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
			return &workspace.WorkspaceResponse{
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

	ws, err := s.postgresRepo.WorkspaceRepo.FindOneByFilter(ctx, &repository.FindWorkspaceByFilter{
		ID:       &id,
		IsActive: data.IsActive,
		IDs:      ids,
	})
	if err != nil {
		customErr = errors.New(errors.ErrCodeWorkspaceNotFound)
		return &workspace.WorkspaceResponse{
			Success: false,
			Data:    nil,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	return &workspace.WorkspaceResponse{
		Success: true,
		Data: &workspace.WorkspaceDetail{
			Id:   ws.ID.String(),
			Name: ws.Name,
		},
		Error: nil,
	}, nil
}

func (s *grpcServer) GetUserWorkspaceByFilter(ctx context.Context, data *workspace.GetUserWorkspaceByFilterParams) (*workspace.UserWorkspaceResponse, error) {
	customErr := errors.New(errors.ErrCodeInternalServerError)

	filter, err := convertUserParamsToFilter(data)
	if err != nil {
		return &workspace.UserWorkspaceResponse{
			Success: false,
			Data:    nil,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	userWs, err := s.postgresRepo.UserWorkspaceRepo.FindOneByFilter(ctx, filter)
	if err != nil {
		customErr = errors.New(errors.ErrCodeUserWorkspaceNotFound)
		return &workspace.UserWorkspaceResponse{
			Success: false,
			Data:    nil,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	return &workspace.UserWorkspaceResponse{
		Success: true,
		Data: &workspace.UserWorkspaceDetail{
			Id:          userWs.ID.String(),
			Name:        userWs.UserWorkspaceDetail.Name,
			Email:       *userWs.Email,
			PhoneNumber: *userWs.PhoneNumber,
			IsActive:    userWs.IsActive,
			Role:        convertUserWSRole(userWs.Role),
			WorkspaceId: userWs.WorkspaceID.String(),
			UserId:      userWs.UserID.String(),
		},
		Error: nil,
	}, nil
}

func (s *grpcServer) GetUserWorkspacesByFilter(ctx context.Context, data *workspace.GetUserWorkspaceByFilterParams) (*workspace.UserWorkspacesResponse, error) {
	var userWSRes []*workspace.UserWorkspaceDetail
	customErr := errors.New(errors.ErrCodeInternalServerError)

	filter, err := convertUserParamsToFilter(data)
	if err != nil {
		return &workspace.UserWorkspacesResponse{
			Success: false,
			Data:    nil,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	userWs, err := s.postgresRepo.UserWorkspaceRepo.FindByFilter(ctx, filter)
	if err != nil {
		customErr = errors.New(errors.ErrCodeUserWorkspaceNotFound)
		return &workspace.UserWorkspacesResponse{
			Success: false,
			Data:    nil,
			Error: grpc_utils.FormatErrorResponse(
				int32(customErr.GetCode()),
				customErr.Error(),
			),
		}, nil
	}

	for _, user := range userWs {
		userWSRes = append(userWSRes, &workspace.UserWorkspaceDetail{
			Id:          user.ID.String(),
			Name:        user.UserWorkspaceDetail.Name,
			Email:       *user.Email,
			PhoneNumber: *user.PhoneNumber,
			IsActive:    user.IsActive,
			Role:        convertUserWSRole(user.Role),
			WorkspaceId: user.WorkspaceID.String(),
			UserId:      user.UserID.String(),
		})
	}

	return &workspace.UserWorkspacesResponse{
		Success: true,
		Data:    userWSRes,
		Error:   nil,
	}, nil
}

// ------------------------ Helper ------------------------
func convertUserParamsToFilter(data *workspace.GetUserWorkspaceByFilterParams) (*repository.FindUserWorkspaceByFilter, error) {
	var id, workspaceId, userId uuid.UUID
	var ids, workspaceIds, userIds []uuid.UUID
	var err error
	limit := int(*data.Limit)
	offset := int(*data.Offset)

	if data.Id != nil {
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

	if data.WorkspaceId != nil {
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

	if data.UserId != nil {
		userId, err = utils.ConvertStringToUUID(*data.WorkspaceId)
		if err != nil {
			return nil, err
		}
	}
	if data.UserIds != nil {
		for _, id := range data.WorkspaceIds {
			convertId, err := utils.ConvertStringToUUID(id)
			if err != nil {
				return nil, err
			}

			userIds = append(userIds, convertId)
		}
	}

	return &repository.FindUserWorkspaceByFilter{
		ID:              &id,
		IDs:             ids,
		WorkspaceID:     &workspaceId,
		WorkspaceIDs:    workspaceIds,
		UserID:          &userId,
		UserIDs:         userIds,
		IsActive:        data.IsActive,
		Limit:           &limit,
		Offset:          &offset,
		IsIncludeDetail: true,
	}, nil
}

func convertUserWSRole(role string) workspace.UserWorkspaceRole {
	switch role {
	case entity.ROLE_OWNER:
		return workspace.UserWorkspaceRole_OWNER
	case entity.ROLE_ADMIN:
		return workspace.UserWorkspaceRole_ADMIN
	default:
		return workspace.UserWorkspaceRole_USER
	}
}
