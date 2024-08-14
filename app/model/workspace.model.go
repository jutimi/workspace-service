package model

type CreateWorkspaceRequest struct {
	Name        string  `json:"name" validate:"required"`
	PhoneNumber string  `json:"phone_number" validate:"required,phone_number"`
	Address     *string `json:"address"`
	Email       string  `json:"email" validate:"required,email"`
}
type CreateWorkspaceResponse struct{}

type UpdateWorkspaceRequest struct {
	Id string `query:"id" validate:"required"`
	CreateWorkspaceRequest
}
type UpdateWorkspaceResponse struct{}

type InactiveWorkspaceRequest struct {
	Id string `query:"id" validate:"required"`
}
type InactiveWorkspaceResponse struct{}
