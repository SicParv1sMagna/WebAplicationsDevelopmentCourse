package roles

type Role int

const (
	User Role = iota
	Moderator
	Admin
)
