package database

func (dbCfg *DbConfig) RevokeToken(token string) error {
	query := `
		INSERT INTO token_revocations (token) VALUES ($1)
		RETURNING token, revoked_at
	`
	var revocation Revocation
	err := dbCfg.database.QueryRow(query, token).Scan(
		&revocation.Token,
		&revocation.RevokedAt,
	)
	if err != nil {
		return err
	}
	return nil
}

func (dbCfg *DbConfig) IsTokenRevoked(token string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 from token_revocations where token = $1
		)
	`
	var tokenIsRevoked bool
	err := dbCfg.database.QueryRow(query, token).Scan(&tokenIsRevoked)
	if err != nil {
		return true, err
	}

	return tokenIsRevoked, nil
}
