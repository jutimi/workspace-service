package helper

import (
	"context"

	"github.com/jutimi/workspace-server/app/entity"
	postgres_repository "github.com/jutimi/workspace-server/app/repository/postgres"
	"github.com/jutimi/workspace-server/package/errors"
)

type userWSHelper struct {
	postgresRepo postgres_repository.PostgresRepositoryCollections
}

func NewUserWSHelper(
	postgresRepo postgres_repository.PostgresRepositoryCollections,
) UserWSHelper {
	return &userWSHelper{
		postgresRepo: postgresRepo,
	}
}

func (h *userWSHelper) CreateUserWS(ctx context.Context, data *CreateUserWsParams) (*entity.UserWorkspace, error) {
	userWS := entity.NewUserWorkspace()
	userWS.BaseWorkspace.WorkspaceID = data.WorkspaceID
	userWS.UserID = data.UserID
	userWS.Email = data.Email
	userWS.PhoneNumber = data.PhoneNumber
	userWS.Role = data.Role

	if err := h.postgresRepo.UserWorkspaceRepo.Create(ctx, data.Tx, userWS); err != nil {
		return nil, errors.New(errors.ErrCodeInternalServerError)
	}

	return userWS, nil
}
