package models

import (
	"encoding/json"
	"fmt"
)

type Role int

const (
	RoleUser Role = iota
	RoleAdmin
)

func (r Role) String() string {
	switch r {
	case RoleUser:
		return "user"
	case RoleAdmin:
		return "admin"
	default:
		return "unknown"
	}
}

func (r Role) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *Role) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	switch s {
	case "user":
		*r = RoleUser
	case "admin":
		*r = RoleAdmin
	default:
		return fmt.Errorf("invalid role: %s", s)
	}
	return nil
}

func ParseRole(s string) (Role, error) {
	switch s {
	case "user":
		return RoleUser, nil
	case "admin":
		return RoleAdmin, nil
	default:
		return RoleUser, fmt.Errorf("invalid role: %s", s)
	}
}
