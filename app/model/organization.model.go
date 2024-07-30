package model

type CreateOrganizationRequest struct {
	Name                 string   `json:"name" validate:"required"`
	ParentOrganizationId *string  `json:"parent_organization_id"`
	LeaderId             *string  `json:"leader_id"`
	SubLeaderIds         []string `json:"sub_leader_ids"`
}
type CreateOrganizationResponse struct{}

type UpdateOrganizationRequest struct {
	Id string `json:"id" validate:"required"`
	CreateOrganizationRequest
}
type UpdateOrganizationResponse struct{}

type RemoveOrganizationRequest struct {
	Id string `params:"id" validate:"required"`
}
type RemoveOrganizationResponse struct{}
