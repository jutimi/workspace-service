package helper

import (
	"context"

	"workspace-server/app/entity"
	postgres_repository "workspace-server/app/repository/postgres"
	"workspace-server/package/errors"
)

type userWorkspaceHelper struct {
	postgresRepo postgres_repository.PostgresRepositoryCollections
}

func NewUserWorkspaceHelper(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) UserWorkspaceHelper {
	return &userWorkspaceHelper{
		postgresRepo: postgresRepo,
	}
}

func (h *userWorkspaceHelper) CreateUserWorkspace(ctx context.Context, data *CreateUserWorkspaceParams) (*entity.UserWorkspace, error) {
	userWorkspace := entity.NewUserWorkspace()
	userWorkspace.WorkspaceId = data.WorkspaceId
	userWorkspace.UserId = data.UserId
	userWorkspace.Email = data.Email
	userWorkspace.PhoneNumber = data.PhoneNumber
	userWorkspace.Role = data.Role

	if err := h.postgresRepo.UserWorkspaceRepo.Create(ctx, data.Tx, userWorkspace); err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return userWorkspace, nil
}
