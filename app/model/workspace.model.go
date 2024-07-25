package model

type CreateWorkspaceRequest struct {
	Name        string  `json:"name" binding:"required"`
	PhoneNumber *string `json:"phone_number"`
	Address     *string `json:"address"`
	Email       *string `json:"email"`
}
type CreateWorkspaceResponse struct{}

type UpdateWorkspaceRequest struct {
	ID string `query:"id" binding:"required"`
	CreateWorkspaceRequest
}
type UpdateWorkspaceResponse struct{}

type InactiveWorkspaceRequest struct {
	ID string `query:"id" binding:"required"`
}
type InactiveWorkspaceResponse struct{}
