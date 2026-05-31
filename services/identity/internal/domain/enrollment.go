package domain

import "time"

type Enrollment struct {
	ID         string
	UserID     string
	CourseID   string
	Role       Role
	EnrolledAt time.Time
}
