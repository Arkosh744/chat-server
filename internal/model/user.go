package model

import desc "github.com/Arkosh744/auth-service-api/pkg/user_v1"

const (
	RoleUnknown = "unknown"
	RoleAdmin   = "admin"
	RoleUser    = "user"
)

type UserIdentifier struct {
	Username string
	Email    string
}

type User struct {
	UserIdentifier

	Role string
}

func ToRole(role desc.Role) string {
	switch role {
	case desc.Role_ADMIN:
		return RoleAdmin
	case desc.Role_USER:
		return RoleUser
	default:
		return RoleUnknown
	}
}
