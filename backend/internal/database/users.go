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
		if isUniqueConstraintError(err) {
			return insertedUser, ErrAlreadyExist
		}
		return insertedUser, err
	}

	return insertedUser, nil
}

func (dbCfg *DbConfig) UpdateUserEmail(userID int, email string) (User, error) {
	updateQuery := `
		UPDATE users SET email = $1
		WHERE id = $2
		RETURNING id, email, hashed_password, created_at
	`

	var updatedUser User

	err := dbCfg.database.QueryRow(updateQuery, email, userID).Scan(
		&updatedUser.ID,
		&updatedUser.Email,
		&updatedUser.HashedPassword,
		&updatedUser.CreatedAt,
	)
	if err != nil {
		return User{}, err
	}

	return updatedUser, nil
}

func (dbCfg *DbConfig) GetUserByEmail(email string) (User, error) {
	query := `
		SELECT id, email, hashed_password, created_at 
		FROM users WHERE email = $1
	`

	var user User
	err := dbCfg.database.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.HashedPassword,
		&user.CreatedAt,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (dbCfg *DbConfig) DeleteUser(userID int) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := dbCfg.database.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}
