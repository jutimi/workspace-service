package model

import "github.com/google/uuid"

type CreateOrganizationRequest struct {
	Name                       string          `json:"name" validate:"required"`
	ParentOrganizationId       *uuid.UUID      `json:"parent_organization_id"`
	ParentOrganizationLeaderId *uuid.UUID      `json:"parent_organization_leader_id"`
	LeaderId                   *uuid.UUID      `json:"leader_id"`
	SubLeaders                 []SubLeaderData `json:"sub_leaders"`
}
type CreateOrganizationResponse struct{}
type SubLeaderData struct {
	SubLeaderId uuid.UUID   `json:"sub_leader_id"` // Sub leader of organization (user workspace id)
	MemberIds   []uuid.UUID `json:"member_ids"`    // Member of sub leader (user workspace ids)
}

type UpdateOrganizationRequest struct {
	Id string `json:"id" validate:"required"`
	CreateOrganizationRequest
}
type UpdateOrganizationResponse struct{}

type RemoveOrganizationRequest struct {
	Id string `params:"id" validate:"required"`
}
type RemoveOrganizationResponse struct{}
