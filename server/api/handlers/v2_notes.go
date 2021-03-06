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

	"github.com/dnote/dnote/server/api/helpers"
	"github.com/dnote/dnote/server/api/operations"
	"github.com/dnote/dnote/server/api/presenters"
	"github.com/dnote/dnote/server/database"
	"github.com/pkg/errors"
)

type createNoteV2Payload struct {
	BookUUID string `json:"book_uuid"`
	Content  string `json:"content"`
	AddedOn  *int64 `json:"added_on"`
	EditedOn *int64 `json:"edited_on"`
}

func validateCreateNoteV2Payload(p createNoteV2Payload) error {
	if p.BookUUID == "" {
		return errors.New("bookUUID is required")
	}

	return nil
}

// CreateNoteV2Resp is a response for creating a note
type CreateNoteV2Resp struct {
	Result presenters.Note `json:"result"`
}

// CreateNoteV2 creates a note
func (a *App) CreateNoteV2(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value(helpers.KeyUser).(database.User)
	if !ok {
		http.Error(w, "Failed to get user", http.StatusInternalServerError)
		return
	}

	var params createNoteV2Payload
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, errors.Wrap(err, "decoding payload").Error(), http.StatusInternalServerError)
		return
	}

	err = validateCreateNoteV2Payload(params)
	if err != nil {
		http.Error(w, errors.Wrap(err, "validating payload").Error(), http.StatusBadRequest)
		return
	}

	var book database.Book
	db := database.DBConn
	if err := db.Where("uuid = ? AND user_id = ?", params.BookUUID, user.ID).First(&book).Error; err != nil {
		http.Error(w, errors.Wrap(err, "finding book").Error(), http.StatusInternalServerError)
		return
	}

	note, err := operations.CreateNote(user, a.Clock, params.BookUUID, params.Content, params.AddedOn, params.EditedOn, false)
	if err != nil {
		http.Error(w, errors.Wrap(err, "creating note").Error(), http.StatusInternalServerError)
		return
	}

	// preload associations
	note.User = user
	note.Book = book

	resp := CreateNoteV2Resp{
		Result: presenters.PresentNote(note),
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// NotesOptionsV2 is a handler for OPTIONS endpoint for notes
func (a *App) NotesOptionsV2(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Authorization, Version")
}
