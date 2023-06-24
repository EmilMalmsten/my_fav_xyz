package database

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func createRandomEmail() string {
	seed := time.Now().UnixNano()
	rng := rand.New(rand.NewSource(seed))
	randomNumber := rng.Intn(100000) + 1
	email := fmt.Sprintf("testuser%d@mail.com", randomNumber)
	return email
}

func insertUser(t *testing.T) User {
	randomEmail := createRandomEmail()
	user := User{
		Email:          randomEmail,
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
