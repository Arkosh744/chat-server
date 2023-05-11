package model

import (
	"strings"
)

type Role int

const (
	RoleUnknown Role = iota // 0
	RoleAdmin               // 1
	RoleUser                // 2
)

type UserIdentifier struct {
	Username string
	Email    string
}

type User struct {
	UserIdentifier

	Password        string
	PasswordConfirm string
	Role            Role
}

func (r Role) String() string {
	switch r {
	case RoleAdmin:
		return "admin"
	case RoleUser:
		return "user"
	default:
		return ""
	}
}

func StringToRole(roleStr string) Role {
	switch strings.ToLower(roleStr) {
	case "admin":
		return RoleAdmin
	case "user":
		return RoleUser
	default:
		return RoleUnknown
	}
}
