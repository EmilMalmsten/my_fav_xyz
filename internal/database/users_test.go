package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func insertUser(t *testing.T) User {
	user := User{
		Email:          "test@mail.com",
		HashedPassword: "asd123123123hjerwehr",
	}

	insertedUser, err := dbTestConfig.InsertUser(user)
	require.NoError(t, err)
	require.NotZero(t, insertedUser)
	require.Equal(t, user.Email, insertedUser.Email)
	require.Equal(t, user.HashedPassword, insertedUser.HashedPassword)
	return insertedUser
}

func TestInsertUser(t *testing.T) {
	insertUser(t)
}
