package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/emilmalmsten/my_top_xyz/backend/internal/auth"
	"github.com/emilmalmsten/my_top_xyz/backend/internal/database"
	"github.com/go-chi/chi"
	"github.com/thanhpk/randstr"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token        string        `json:"token"`
	RefreshToken string        `json:"refresh_token"`
	User         database.User `json:"user"`
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	loginRequest := LoginRequest{}
	err := decoder.Decode(&loginRequest)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	dbUser, err := cfg.DB.GetUserByEmail(loginRequest.Email)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user")
		return
	}

	err = auth.CheckPasswordHash(loginRequest.Password, dbUser.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password")
		return
	}

	accessToken, err := auth.MakeJWT(
		dbUser.ID,
		cfg.jwtSecret,
		time.Hour,
		auth.TokenTypeAccess,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create access JWT")
		return
	}

	refreshToken, err := auth.MakeJWT(
		dbUser.ID,
		cfg.jwtSecret,
		time.Hour*24*7,
		auth.TokenTypeRefresh,
	)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh JWT")
		return
	}

	respondWithJSON(w, http.StatusOK, LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         dbUser,
	})
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

func (cfg *apiConfig) handlerForgotPassword(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var forgotPasswordRequest ForgotPasswordRequest
	err := decoder.Decode(&forgotPasswordRequest)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	dbUser, err := cfg.DB.GetUserByEmail(forgotPasswordRequest.Email)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusNotFound, "Couldn't find user")
		return
	}

	resetToken := randstr.String(20)
	passwordResetToken := auth.Encode(resetToken)
	dbUser.PasswordResetToken = passwordResetToken
	dbUser.PasswordResetTokenExpireAt = time.Now().Add(time.Minute * 15)

	err = cfg.DB.InsertPasswordResetToken(dbUser)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Failed to save pw reset token")
		return
	}

	emailData := EmailData{
		URL:     "/api/resetpassword/" + resetToken,
		Subject: "Your password reset token (valid for 15min)",
	}

	err = cfg.SendEmail(&dbUser, &emailData, "resetPassword.html")
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Failed to send email")
		return
	}

	message := "You will receive a reset email if user with that email exist"
	respondWithJSON(w, http.StatusOK, message)
}

type ResetPasswordRequest struct {
	Password string `json:"password"`
}

func (cfg *apiConfig) handlerResetPassword(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var resetPasswordRequest ResetPasswordRequest
	err := decoder.Decode(&resetPasswordRequest)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	resetToken := chi.URLParam(r, "resetToken")
	encodedResetToken := auth.Encode(resetToken)

	hashedPassword, err := auth.HashPassword(resetPasswordRequest.Password)
	if err != nil {
		fmt.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash new password")
		return
	}

	err = cfg.DB.ResetPassword(hashedPassword, encodedResetToken)
	if err != nil {
		fmt.Println(err)
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Incorrect reset token")
			return
		} else if errors.Is(err, database.ErrIsExpired) {
			respondWithError(w, http.StatusBadRequest, "Reset token is already expired")
			return
		}
		respondWithError(w, http.StatusInternalServerError, "Failed to reset password")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}
