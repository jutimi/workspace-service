package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserPayload struct {
	ID uuid.UUID `json:"id"`
	jwt.RegisteredClaims
}

type WorkspacePayload struct {
	UserPayload
	WorkspaceID uuid.UUID `json:"workspace_id"`
	jwt.RegisteredClaims
}

type UserWorkspacePayload struct {
	WorkspacePayload
	UserWorkspaceID uuid.UUID `json:"user_workspace_id"`
	jwt.RegisteredClaims
}

// Get payload from user token
func ParseUserToken(tokenString string) (*UserPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserPayload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserPayload); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenMalformed
}

// Get payload from workspace token
func ParseWSToken(tokenString string) (*WorkspacePayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &WorkspacePayload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*WorkspacePayload); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenMalformed
}

// Get payload from user workspace token
func ParseUserWSToken(tokenString string) (*UserWorkspacePayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &UserWorkspacePayload{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(""), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*UserWorkspacePayload); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenMalformed
}
