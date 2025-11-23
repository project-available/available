package utils

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

type Status string

const (
	Booked    Status = "booked"
	Available Status = "available"
	Pending   Status = "pending"
)
