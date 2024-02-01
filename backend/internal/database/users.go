package database

import (
	"database/sql"
	"errors"
	"log"
)

func (dbCfg *DbConfig) InsertUser(user User) (User, error) {
	insertQuery := `
		INSERT INTO users (email, username, hashed_password)
		VALUES ($1, $2, $3) 
		RETURNING id, email, username, hashed_password
	`
	var insertedUser User
	err := dbCfg.database.QueryRow(insertQuery, user.Email, user.Username, user.HashedPassword).Scan(
		&insertedUser.ID,
		&insertedUser.Email,
		&insertedUser.Username,
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

func (dbCfg *DbConfig) UpdateUserPassword(userID int, hashedPassword string) (User, error) {
	updateQuery := `
		UPDATE users SET hashed_password = $1
		WHERE id = $2
		RETURNING id, email, hashed_password, created_at
	`

	var updatedUser User

	err := dbCfg.database.QueryRow(updateQuery, hashedPassword, userID).Scan(
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
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrNotExist
		}
		return User{}, err
	}

	return user, nil
}

func (dbCfg *DbConfig) GetUserByID(userID int) (User, error) {
	query := `
		SELECT id, email, username, hashed_password, created_at 
		FROM users WHERE id = $1
	`

	var user User
	err := dbCfg.database.QueryRow(query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.HashedPassword,
		&user.CreatedAt,
	)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (dbCfg *DbConfig) UserWithEmailExists(email string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM users 
			WHERE email = $1
		)
	`

	var exists bool
	err := dbCfg.database.QueryRow(query, email).Scan(&exists)
	if err != nil {
		log.Printf("Failed to check if email exists: %v", err)
		return false, err
	}

	return exists, nil
}

func (dbCfg *DbConfig) UserWithUsernameExists(username string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1 
			FROM users 
			WHERE username = $1
		)
	`

	var exists bool
	err := dbCfg.database.QueryRow(query, username).Scan(&exists)
	if err != nil {
		log.Printf("Failed to check if username exists: %v", err)
		return false, err
	}

	return exists, nil
}

func (dbCfg *DbConfig) DeleteUser(userID int) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := dbCfg.database.Exec(query, userID)
	if err != nil {
		return err
	}
	return nil
}

func (dbCfg *DbConfig) InsertPasswordResetToken(user User) error {
	query := `
		UPDATE users SET password_reset_token = $1,
		password_reset_token_expire_at = $2
		WHERE id = $3
	`
	_, err := dbCfg.database.Exec(query, user.PasswordResetToken, user.PasswordResetTokenExpireAt, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (dbCfg *DbConfig) ResetPassword(newPasswordHash, resetToken string) error {
	var tokenExists bool
	err := dbCfg.database.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE password_reset_token = $1)", resetToken).Scan(&tokenExists)
	if err != nil {
		return err
	}
	if !tokenExists {
		return ErrNotExist
	}

	query := `
		UPDATE users
		SET hashed_password = $1,
		password_reset_token = null
		WHERE password_reset_token = $2
		AND current_timestamp < password_reset_token_expire_at
	`

	result, err := dbCfg.database.Exec(query, newPasswordHash, resetToken)
	if err != nil {
		return err
	}

	rowCount, _ := result.RowsAffected()
	if rowCount == 0 {
		return ErrIsExpired
	}

	return nil
}
