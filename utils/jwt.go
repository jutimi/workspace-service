package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserPayload struct {
	Id    uuid.UUID `json:"id"`
	Scope string    `json:"scopes"`
	jwt.RegisteredClaims
}

type WorkspacePayload struct {
	Id              uuid.UUID `json:"id"`
	Scope           string    `json:"scopes"`
	WorkspaceId     uuid.UUID `json:"workspace_id"`
	UserWorkspaceId uuid.UUID `json:"user_workspace_id"`
	jwt.RegisteredClaims
}

// Get payload from user token
func ParseUserToken(tokenString string) (*UserPayload, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &UserPayload{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserPayload); ok {
		return claims, nil
	}

	return nil, jwt.ErrTokenMalformed
}

// Get payload from workspace token
func ParseWorkspaceToken(tokenString string) (*WorkspacePayload, error) {
	token, _, err := new(jwt.Parser).ParseUnverified(tokenString, &WorkspacePayload{})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*WorkspacePayload); ok {
		return claims, nil
	}

	return nil, jwt.ErrTokenMalformed
}
