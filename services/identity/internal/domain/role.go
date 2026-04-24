package domain

type Role string

const (
	RoleStudent    Role = "student"
	RoleInstructor Role = "instructor"
	RoleAdmin      Role = "admin"
)

func (r Role) Valid() bool {
	switch r {
	case RoleStudent, RoleInstructor, RoleAdmin:
		return true
	}
	return false
}

func (r Role) CanCreateCourse() bool {
	return r == RoleInstructor || r == RoleAdmin
}
