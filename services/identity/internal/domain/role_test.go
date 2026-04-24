package domain

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoleValid(t *testing.T) {
	require.Equal(t, true, RoleAdmin.Valid())
	require.Equal(t, true, RoleInstructor.Valid())
	require.Equal(t, true, RoleStudent.Valid())

	require.Equal(t, false, Role("omg").Valid())
}

func TestRoleCanCreateCourse(t *testing.T) {
	require.Equal(t, true, RoleInstructor.CanCreateCourse())
	require.Equal(t, true, RoleAdmin.CanCreateCourse())
	require.Equal(t, false, RoleStudent.CanCreateCourse())
}
