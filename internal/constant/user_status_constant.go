package constant

type UserStatusEnum string

const (
	UserStatusActive   UserStatusEnum = "active"
	UserStatusInactive UserStatusEnum = "inactive"
	UserStatusBanned   UserStatusEnum = "banned"
)
