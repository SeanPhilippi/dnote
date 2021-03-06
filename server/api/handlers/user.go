/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote.
 *
 * Dnote is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Dnote.  If not, see <https://www.gnu.org/licenses/>.
 */

package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/dnote/dnote/server/api/crypt"
	"github.com/dnote/dnote/server/api/helpers"
	"github.com/dnote/dnote/server/api/operations"
	"github.com/dnote/dnote/server/database"
	"github.com/dnote/dnote/server/mailer"

	"github.com/pkg/errors"
)

type updateProfilePayload struct {
	Name string `json:"name"`
}

// updateProfile updates user
func (a *App) updateProfile(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn

	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		http.Error(w, "No authenticated user found", http.StatusInternalServerError)
		return
	}

	var params updateProfilePayload
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Validate
	if len(params.Name) > 50 {
		http.Error(w, "Name is too long", http.StatusBadRequest)
		return
	}

	var account database.Account
	err = db.Where("user_id = ?", user.ID).First(&account).Error
	if err != nil {
		http.Error(w, errors.Wrap(err, "finding account").Error(), http.StatusInternalServerError)
		return
	}

	tx := db.Begin()
	user.Name = params.Name
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		http.Error(w, errors.Wrap(err, "saving user").Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		http.Error(w, errors.Wrap(err, "saving user").Error(), http.StatusInternalServerError)
		return
	}
	tx.Commit()

	session := makeSession(user, account)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(session); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type updateEmailPayload struct {
	NewEmail        string `json:"new_email"`
	NewCipherKeyEnc string `json:"new_cipher_key_enc"`
	OldAuthKey      string `json:"old_auth_key"`
	NewAuthKey      string `json:"new_auth_key"`
}

// updateEmail updates user
func (a *App) updateEmail(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn

	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		http.Error(w, "No authenticated user found", http.StatusInternalServerError)
		return
	}

	var params updateEmailPayload
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Validate
	if len(params.NewEmail) > 100 {
		http.Error(w, "Email is too long", http.StatusBadRequest)
		return
	}

	var account database.Account
	err = db.Where("user_id = ?", user.ID).First(&account).Error
	if err != nil {
		http.Error(w, errors.Wrap(err, "finding account").Error(), http.StatusInternalServerError)
		return
	}

	authKeyHash := crypt.HashAuthKey(params.OldAuthKey, account.Salt, account.ServerKDFIteration)
	if account.AuthKeyHash != authKeyHash {
		http.Error(w, "wrong password", http.StatusUnauthorized)
		return
	}

	if account.Email.String == params.NewEmail {
		http.Error(w, "New email is the same as the old", http.StatusBadRequest)
		return
	}

	tx := db.Begin()

	account.Email = database.ToNullString(params.NewEmail)
	account.CipherKeyEnc = params.NewCipherKeyEnc
	account.AuthKeyHash = crypt.HashAuthKey(params.NewAuthKey, account.Salt, crypt.ServerKDFIteration)
	account.ServerKDFIteration = crypt.ServerKDFIteration
	account.EmailVerified = false

	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		http.Error(w, errors.Wrap(err, "saving account").Error(), http.StatusInternalServerError)
		return
	}
	tx.Commit()

	session := makeSession(user, account)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(session); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func respondWithCalendar(w http.ResponseWriter, userID int) {
	db := database.DBConn

	rows, err := db.Table("notes").Select("COUNT(id), date(to_timestamp(added_on/1000000000)) AS added_date").
		Where("user_id = ?", userID).
		Group("added_date").
		Order("added_date DESC").Rows()

	if err != nil {
		http.Error(w, errors.Wrap(err, "Failed to count lessons").Error(), http.StatusInternalServerError)
		return
	}

	payload := map[string]int{}

	for rows.Next() {
		var count int
		var d time.Time

		if err := rows.Scan(&count, &d); err != nil {
			http.Error(w, errors.Wrap(err, "counting notes").Error(), http.StatusInternalServerError)
		}
		payload[d.Format("2006-1-2")] = count
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) getCalendar(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		http.Error(w, "No authenticated user found", http.StatusInternalServerError)
		return
	}

	respondWithCalendar(w, user.ID)
}

func (a *App) getDemoCalendar(w http.ResponseWriter, r *http.Request) {
	userID, err := helpers.GetDemoUserID()
	if err != nil {
		http.Error(w, errors.Wrap(err, "finding demo user").Error(), http.StatusInternalServerError)
		return
	}

	respondWithCalendar(w, userID)
}

func (a *App) createVerificationToken(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn

	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		http.Error(w, "No authenticated user found", http.StatusInternalServerError)
		return
	}

	var account database.Account
	err := db.Where("user_id = ?", user.ID).First(&account).Error
	if err != nil {
		http.Error(w, errors.WithMessage(err, "finding account").Error(), http.StatusInternalServerError)
		return
	}

	if account.EmailVerified {
		http.Error(w, "Email already verified", http.StatusGone)
		return
	}
	if !account.Email.Valid {
		http.Error(w, "Email not set", http.StatusUnprocessableEntity)
		return
	}

	tokenValue, err := generateVerificationCode()
	if err != nil {
		http.Error(w, errors.Wrap(err, "generating verification code").Error(), http.StatusInternalServerError)
		return
	}

	token := database.Token{
		UserID: account.UserID,
		Value:  tokenValue,
		Type:   database.TokenTypeEmailVerification,
	}

	if err := db.Save(&token).Error; err != nil {
		http.Error(w, errors.Wrap(err, "saving token").Error(), http.StatusInternalServerError)
		return
	}

	subject := "Verify your email"
	data := struct {
		Subject string
		Token   string
	}{
		subject,
		tokenValue,
	}
	email := mailer.NewEmail("noreply@dnote.io", []string{account.Email.String}, subject)
	if err := email.ParseTemplate(mailer.EmailTypeEmailVerification, data); err != nil {
		http.Error(w, errors.Wrap(err, "parsing template").Error(), http.StatusInternalServerError)
		return
	}

	if err := email.Send(); err != nil {
		http.Error(w, errors.Wrap(err, "sending email").Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

type verifyEmailPayload struct {
	Token string `json:"token"`
}

func (a *App) verifyEmail(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn

	var params verifyEmailPayload
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, errors.Wrap(err, "decoding payload").Error(), http.StatusInternalServerError)
		return
	}

	var token database.Token
	if err := db.
		Where("value = ? AND type = ?", params.Token, database.TokenTypeEmailVerification).
		First(&token).Error; err != nil {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	if token.UsedAt != nil {
		http.Error(w, "invalid token", http.StatusBadRequest)
		return
	}

	// Expire after ttl
	if time.Since(token.CreatedAt).Minutes() > 30 {
		http.Error(w, "This link has been expired. Please request a new link.", http.StatusGone)
		return
	}

	var account database.Account
	if err := db.Where("user_id = ?", token.UserID).First(&account).Error; err != nil {
		http.Error(w, errors.WithMessage(err, "finding account").Error(), http.StatusInternalServerError)
		return
	}
	if account.EmailVerified {
		http.Error(w, "Already verified", http.StatusConflict)
		return
	}

	tx := db.Begin()
	account.EmailVerified = true
	if err := tx.Save(&account).Error; err != nil {
		tx.Rollback()
		http.Error(w, errors.WithMessage(err, "updating email_verified").Error(), http.StatusInternalServerError)
		return
	}
	if err := tx.Model(&token).Update("used_at", time.Now()).Error; err != nil {
		tx.Rollback()
		http.Error(w, errors.WithMessage(err, "updating reset token").Error(), http.StatusInternalServerError)
		return
	}
	tx.Commit()

	var user database.User
	if err := db.Where("id = ?", token.UserID).First(&user).Error; err != nil {
		http.Error(w, errors.WithMessage(err, "finding user").Error(), http.StatusInternalServerError)
		return
	}

	session := makeSession(user, account)
	setAuthCookie(w, user)
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(session); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type updateEmailPreferencePayload struct {
	DigestWeekly bool `json:"digest_weekly"`
}

func (a *App) updateEmailPreference(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn

	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		http.Error(w, "No authenticated user found", http.StatusInternalServerError)
		return
	}

	var params updateEmailPreferencePayload
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, errors.Wrap(err, "decoding payload").Error(), http.StatusInternalServerError)
		return
	}

	var frequency database.EmailPreference
	if err := db.Where(database.EmailPreference{UserID: user.ID}).FirstOrCreate(&frequency).Error; err != nil {
		http.Error(w, errors.Wrap(err, "finding frequency").Error(), http.StatusInternalServerError)
		return
	}

	tx := db.Begin()

	frequency.DigestWeekly = params.DigestWeekly
	if err := tx.Save(&frequency).Error; err != nil {
		tx.Rollback()
		http.Error(w, errors.Wrap(err, "saving frequency").Error(), http.StatusInternalServerError)
		return
	}

	token, ok := r.Context().Value(helpers.KeyToken).(database.Token)
	if ok {
		// Use token if the user was authenticated by token
		if err := tx.Model(&token).Update("used_at", time.Now()).Error; err != nil {
			tx.Rollback()
			http.Error(w, errors.WithMessage(err, "updating reset token").Error(), http.StatusInternalServerError)
			return
		}
	}

	tx.Commit()

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(frequency); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *App) getEmailPreference(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn

	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		http.Error(w, "No authenticated user found", http.StatusInternalServerError)
		return
	}

	var pref database.EmailPreference
	if err := db.Where(database.EmailPreference{UserID: user.ID}).First(&pref).Error; err != nil {
		http.Error(w, errors.Wrap(err, "finding pref").Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(pref); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

type updatePasswordPayload struct {
	OldAuthKey      string `json:"old_auth_key"`
	NewAuthKey      string `json:"new_auth_key"`
	NewCipherKeyEnc string `json:"new_cipher_key_enc"`
	NewKDFIteration int    `json:"new_kdf_iteration"`
}

func (a *App) updatePassword(w http.ResponseWriter, r *http.Request) {
	db := database.DBConn

	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		http.Error(w, "No authenticated user found", http.StatusInternalServerError)
		return
	}
	var account database.Account
	if err := db.Where("user_id = ?", user.ID).First(&account).Error; err != nil {
		http.Error(w, errors.Wrap(err, "getting account").Error(), http.StatusInternalServerError)
		return
	}
	var params updatePasswordPayload
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	oldAuthKeyHash := crypt.HashAuthKey(params.OldAuthKey, account.Salt, account.ServerKDFIteration)
	if oldAuthKeyHash != account.AuthKeyHash {
		http.Error(w, ErrLoginFailure.Error(), http.StatusUnauthorized)
		return
	}

	newAuthKeyHash := crypt.HashAuthKey(params.NewAuthKey, account.Salt, account.ServerKDFIteration)

	if err := db.
		Model(&account).
		Debug().
		Updates(map[string]interface{}{
			"auth_key_hash":        newAuthKeyHash,
			"client_kdf_iteration": params.NewKDFIteration,
			"server_kdf_iteration": account.ServerKDFIteration,
			"cipher_key_enc":       params.NewCipherKeyEnc,
		}).Error; err != nil {
		http.Error(w, errors.Wrap(err, "updating account").Error(), http.StatusInternalServerError)
		return
	}

	if err := operations.DeleteUserSessions(db, user.ID); err != nil {
		http.Error(w, "deleting user sessions", http.StatusBadRequest)
		return
	}

	respondWithSession(w, user.ID, account.CipherKeyEnc)
}
