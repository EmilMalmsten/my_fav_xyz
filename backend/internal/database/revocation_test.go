package database

import (
	"testing"
	"time"

	"github.com/emilmalmsten/my_top_xyz/backend/internal/auth"
	"github.com/stretchr/testify/require"
)

func revokeToken(t *testing.T, refreshToken string) {
	revocation, err := dbTestConfig.RevokeToken(refreshToken)
	require.NoError(t, err)
	require.NotEmpty(t, revocation)

	require.Equal(t, refreshToken, revocation.Token)
	require.WithinDuration(t, time.Now(), revocation.RevokedAt, time.Second)
}

func makeJWT(t *testing.T, userID int) string {
	tokenTestSecret := "Xme6qJa1XiPxagxlXs+CuRm2Nam7fUaTe95igkc66mARBNE0DA3dfRRws17B4WTEJlpWWmmpOL+aVPPfebSung=="

	refreshToken, err := auth.MakeJWT(userID, tokenTestSecret, time.Hour*24*7, auth.TokenTypeRefresh)
	require.NoError(t, err)
	require.NotEmpty(t, refreshToken)
	return refreshToken
}

func TestRevokeToken(t *testing.T) {
	userID := 1
	refreshToken := makeJWT(t, userID)
	revokeToken(t, refreshToken)
}

func TestIsTokenRevoked(t *testing.T) {
	userID := 2
	refreshToken := makeJWT(t, userID)

	isRevoked, err := dbTestConfig.IsTokenRevoked(refreshToken)
	require.NoError(t, err)
	require.False(t, isRevoked)

	revokeToken(t, refreshToken)
	isRevoked, err = dbTestConfig.IsTokenRevoked(refreshToken)
	require.NoError(t, err)
	require.True(t, isRevoked)
}
