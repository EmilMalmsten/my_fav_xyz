package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"

	"github.com/emilmalmsten/my_top_xyz/backend/internal/auth"
	"github.com/emilmalmsten/my_top_xyz/backend/internal/database"
)

type createUserRequest struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type resp struct {
		Id int `json:"id"`
	}

	decoder := json.NewDecoder(r.Body)
	createUserRequest := createUserRequest{}
	err := decoder.Decode(&createUserRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	if !validEmail(createUserRequest.Email) {
		respondWithError(w, http.StatusBadRequest, "Incorrect email format")
		return
	}

	if len(createUserRequest.Password) < 8 {
		respondWithError(w, http.StatusBadRequest, "Password needs to be minimum 8 characters")
		return
	}

	if len(createUserRequest.Username) < 3 {
		respondWithError(w, http.StatusBadRequest, "Username needs to be minimum 3 characters")
		return
	}

	emailExists, err := cfg.DB.UserWithEmailExists(createUserRequest.Email)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Error occurred when creating new user")
		return
	}
	if emailExists {
		respondWithError(w, http.StatusBadRequest, "A user with that email already exists")
		return
	}

	usernameExists, err := cfg.DB.UserWithUsernameExists(createUserRequest.Username)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Error occurred when creating new user")
		return
	}
	if usernameExists {
		respondWithError(w, http.StatusBadRequest, "A user with that username already exists")
		return
	}

	hashedPassword, err := auth.HashPassword(createUserRequest.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password")
		return
	}

	createdUser, err := cfg.DB.InsertUser(database.User{
		Email:          createUserRequest.Email,
		Username:       createUserRequest.Username,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, database.ErrAlreadyExist) {
			respondWithError(w, http.StatusBadRequest, "Unique constraint creating user")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Error occurred when creating new user")
		return
	}

	respondWithJSON(w, http.StatusCreated, resp{
		Id: createdUser.ID,
	})
}

func validEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
