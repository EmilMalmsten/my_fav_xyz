package database

func (dbCfg *DbConfig) InsertUser(user User) (User, error) {
	insertQuery := `
		INSERT INTO users (email, hashed_password)
		VALUES ($1, $2) 
		RETURNING id, email, hashed_password
	`
	var insertedUser User
	err := dbCfg.database.QueryRow(insertQuery, user.Email, user.HashedPassword).Scan(
		&insertedUser.ID,
		&insertedUser.Email,
		&insertedUser.HashedPassword,
	)
	if err != nil {
		return insertedUser, err
	}

	return insertedUser, nil
}
