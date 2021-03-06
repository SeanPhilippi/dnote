/* Copyright (C) 2019 Monomax Software Pty Ltd
 *
 * This file is part of Dnote CLI.
 *
 * Dnote CLI is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * Dnote CLI is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with Dnote CLI.  If not, see <https://www.gnu.org/licenses/>.
 */

// Package client provides interfaces for interacting with the Dnote server
// and the data structures for responses
package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dnote/dnote/cli/crypt"
	"github.com/dnote/dnote/cli/infra"
	"github.com/dnote/dnote/cli/utils"
	"github.com/pkg/errors"
)

// ErrInvalidLogin is an error for invalid credentials for login
var ErrInvalidLogin = errors.New("wrong credentials")

// GetSyncStateResp is the response get sync state endpoint
type GetSyncStateResp struct {
	FullSyncBefore int   `json:"full_sync_before"`
	MaxUSN         int   `json:"max_usn"`
	CurrentTime    int64 `json:"current_time"`
}

// GetSyncState gets the sync state response from the server
func GetSyncState(ctx infra.DnoteCtx) (GetSyncStateResp, error) {
	var ret GetSyncStateResp

	hc := http.Client{}
	res, err := utils.DoAuthorizedReq(ctx, hc, "GET", "/v1/sync/state", "")
	if err != nil {
		return ret, errors.Wrap(err, "constructing http request")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ret, errors.Wrap(err, "reading the response body")
	}

	if err = json.Unmarshal(body, &ret); err != nil {
		return ret, errors.Wrap(err, "unmarshalling the payload")
	}

	return ret, nil
}

// SyncFragNote represents a note in a sync fragment and contains only the necessary information
// for the client to sync the note locally
type SyncFragNote struct {
	UUID      string    `json:"uuid"`
	BookUUID  string    `json:"book_uuid"`
	USN       int       `json:"usn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AddedOn   int64     `json:"added_on"`
	EditedOn  int64     `json:"edited_on"`
	Body      string    `json:"content"`
	Public    bool      `json:"public"`
	Deleted   bool      `json:"deleted"`
}

// SyncFragBook represents a book in a sync fragment and contains only the necessary information
// for the client to sync the note locally
type SyncFragBook struct {
	UUID      string    `json:"uuid"`
	USN       int       `json:"usn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	AddedOn   int64     `json:"added_on"`
	Label     string    `json:"label"`
	Deleted   bool      `json:"deleted"`
}

// SyncFragment contains a piece of information about the server's state.
type SyncFragment struct {
	FragMaxUSN    int            `json:"frag_max_usn"`
	UserMaxUSN    int            `json:"user_max_usn"`
	CurrentTime   int64          `json:"current_time"`
	Notes         []SyncFragNote `json:"notes"`
	Books         []SyncFragBook `json:"books"`
	ExpungedNotes []string       `json:"expunged_notes"`
	ExpungedBooks []string       `json:"expunged_books"`
}

// GetSyncFragmentResp is the response from the get sync fragment endpoint
type GetSyncFragmentResp struct {
	Fragment SyncFragment `json:"fragment"`
}

// GetSyncFragment gets a sync fragment response from the server
func GetSyncFragment(ctx infra.DnoteCtx, afterUSN int) (GetSyncFragmentResp, error) {
	v := url.Values{}
	v.Set("after_usn", strconv.Itoa(afterUSN))
	queryStr := v.Encode()

	path := fmt.Sprintf("/v1/sync/fragment?%s", queryStr)
	hc := http.Client{}
	res, err := utils.DoAuthorizedReq(ctx, hc, "GET", path, "")

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return GetSyncFragmentResp{}, errors.Wrap(err, "reading the response body")
	}

	var resp GetSyncFragmentResp
	if err = json.Unmarshal(body, &resp); err != nil {
		return resp, errors.Wrap(err, "unmarshalling the payload")
	}

	return resp, nil
}

// RespBook is the book in the response from the create book api
type RespBook struct {
	ID        int       `json:"id"`
	UUID      string    `json:"uuid"`
	USN       int       `json:"usn"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Label     string    `json:"label"`
}

// CreateBookPayload is a payload for creating a book
type CreateBookPayload struct {
	Name string `json:"name"`
}

// CreateBookResp is the response from create book api
type CreateBookResp struct {
	Book RespBook `json:"book"`
}

// checkRespErr checks if the given http response indicates an error. It returns a boolean indicating
// if the response is an error, and a decoded error message.
func checkRespErr(res *http.Response) (bool, string, error) {
	if res.StatusCode < 400 {
		return false, "", nil
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return true, "", errors.Wrapf(err, "server responded with %d but could not read the response body", res.StatusCode)
	}

	bodyStr := string(body)
	message := fmt.Sprintf(`response %d "%s"`, res.StatusCode, strings.TrimRight(bodyStr, "\n"))
	return true, message, nil
}

// CreateBook creates a new book in the server
func CreateBook(ctx infra.DnoteCtx, label string) (CreateBookResp, error) {
	encLabel, err := crypt.AesGcmEncrypt(ctx.CipherKey, []byte(label))
	if err != nil {
		return CreateBookResp{}, errors.Wrap(err, "encrypting the label")
	}

	payload := CreateBookPayload{
		Name: encLabel,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return CreateBookResp{}, errors.Wrap(err, "marshaling payload")
	}

	hc := http.Client{}
	res, err := utils.DoAuthorizedReq(ctx, hc, "POST", "/v2/books", string(b))
	if err != nil {
		return CreateBookResp{}, errors.Wrap(err, "posting a book to the server")
	}

	hasErr, message, err := checkRespErr(res)
	if err != nil {
		return CreateBookResp{}, errors.Wrap(err, "checking repsonse error")
	}
	if hasErr {
		return CreateBookResp{}, errors.New(message)
	}

	var resp CreateBookResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return resp, errors.Wrap(err, "decoding response payload")
	}

	return resp, nil
}

type updateBookPayload struct {
	Name *string `json:"name"`
}

// UpdateBookResp is the response from create book api
type UpdateBookResp struct {
	Book RespBook `json:"book"`
}

// UpdateBook updates a book in the server
func UpdateBook(ctx infra.DnoteCtx, label, uuid string) (UpdateBookResp, error) {
	encName, err := crypt.AesGcmEncrypt(ctx.CipherKey, []byte(label))
	if err != nil {
		return UpdateBookResp{}, errors.Wrap(err, "encrypting the content")
	}

	payload := updateBookPayload{
		Name: &encName,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return UpdateBookResp{}, errors.Wrap(err, "marshaling payload")
	}

	hc := http.Client{}
	endpoint := fmt.Sprintf("/v1/books/%s", uuid)
	res, err := utils.DoAuthorizedReq(ctx, hc, "PATCH", endpoint, string(b))
	if err != nil {
		return UpdateBookResp{}, errors.Wrap(err, "posting a book to the server")
	}

	hasErr, message, err := checkRespErr(res)
	if err != nil {
		return UpdateBookResp{}, errors.Wrap(err, "checking repsonse error")
	}
	if hasErr {
		return UpdateBookResp{}, errors.New(message)
	}

	var resp UpdateBookResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return resp, errors.Wrap(err, "decoding payload")
	}

	return resp, nil
}

// DeleteBookResp is the response from create book api
type DeleteBookResp struct {
	Status int      `json:"status"`
	Book   RespBook `json:"book"`
}

// DeleteBook deletes a book in the server
func DeleteBook(ctx infra.DnoteCtx, uuid string) (DeleteBookResp, error) {
	hc := http.Client{}
	endpoint := fmt.Sprintf("/v1/books/%s", uuid)
	res, err := utils.DoAuthorizedReq(ctx, hc, "DELETE", endpoint, "")
	if err != nil {
		return DeleteBookResp{}, errors.Wrap(err, "deleting a book in the server")
	}

	hasErr, message, err := checkRespErr(res)
	if err != nil {
		return DeleteBookResp{}, errors.Wrap(err, "checking repsonse error")
	}
	if hasErr {
		return DeleteBookResp{}, errors.New(message)
	}

	var resp DeleteBookResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return resp, errors.Wrap(err, "decoding the response")
	}

	return resp, nil
}

// CreateNotePayload is a payload for creating a note
type CreateNotePayload struct {
	BookUUID string `json:"book_uuid"`
	Body     string `json:"content"`
}

// CreateNoteResp is the response from create note endpoint
type CreateNoteResp struct {
	Result RespNote `json:"result"`
}

type respNoteBook struct {
	UUID  string `json:"uuid"`
	Label string `json:"label"`
}

type respNoteUser struct {
	Name string `json:"name"`
}

// RespNote is a note in the response
type RespNote struct {
	UUID      string       `json:"uuid"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
	Body      string       `json:"content"`
	AddedOn   int64        `json:"added_on"`
	Public    bool         `json:"public"`
	USN       int          `json:"usn"`
	Book      respNoteBook `json:"book"`
	User      respNoteUser `json:"user"`
}

// CreateNote creates a note in the server
func CreateNote(ctx infra.DnoteCtx, bookUUID, content string) (CreateNoteResp, error) {
	encBody, err := crypt.AesGcmEncrypt(ctx.CipherKey, []byte(content))
	if err != nil {
		return CreateNoteResp{}, errors.Wrap(err, "encrypting the content")
	}

	payload := CreateNotePayload{
		BookUUID: bookUUID,
		Body:     encBody,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return CreateNoteResp{}, errors.Wrap(err, "marshaling payload")
	}

	hc := http.Client{}
	res, err := utils.DoAuthorizedReq(ctx, hc, "POST", "/v2/notes", string(b))
	if err != nil {
		return CreateNoteResp{}, errors.Wrap(err, "posting a book to the server")
	}

	hasErr, message, err := checkRespErr(res)
	if err != nil {
		return CreateNoteResp{}, errors.Wrap(err, "checking repsonse error")
	}
	if hasErr {
		return CreateNoteResp{}, errors.New(message)
	}

	var resp CreateNoteResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return CreateNoteResp{}, errors.Wrap(err, "decoding payload")
	}

	return resp, nil
}

type updateNotePayload struct {
	BookUUID *string `json:"book_uuid"`
	Body     *string `json:"content"`
	Public   *bool   `json:"public"`
}

// UpdateNoteResp is the response from create book api
type UpdateNoteResp struct {
	Status int      `json:"status"`
	Result RespNote `json:"result"`
}

// UpdateNote updates a note in the server
func UpdateNote(ctx infra.DnoteCtx, uuid, bookUUID, content string, public bool) (UpdateNoteResp, error) {
	encBody, err := crypt.AesGcmEncrypt(ctx.CipherKey, []byte(content))
	if err != nil {
		return UpdateNoteResp{}, errors.Wrap(err, "encrypting the content")
	}

	payload := updateNotePayload{
		BookUUID: &bookUUID,
		Body:     &encBody,
		Public:   &public,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return UpdateNoteResp{}, errors.Wrap(err, "marshaling payload")
	}

	hc := http.Client{}
	endpoint := fmt.Sprintf("/v1/notes/%s", uuid)
	res, err := utils.DoAuthorizedReq(ctx, hc, "PATCH", endpoint, string(b))
	if err != nil {
		return UpdateNoteResp{}, errors.Wrap(err, "patching a note to the server")
	}

	hasErr, message, err := checkRespErr(res)
	if err != nil {
		return UpdateNoteResp{}, errors.Wrap(err, "checking repsonse error")
	}
	if hasErr {
		return UpdateNoteResp{}, errors.New(message)
	}

	var resp UpdateNoteResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return UpdateNoteResp{}, errors.Wrap(err, "decoding payload")
	}

	return resp, nil
}

// DeleteNoteResp is the response from remove note api
type DeleteNoteResp struct {
	Status int      `json:"status"`
	Result RespNote `json:"result"`
}

// DeleteNote removes a note in the server
func DeleteNote(ctx infra.DnoteCtx, uuid string) (DeleteNoteResp, error) {
	hc := http.Client{}
	endpoint := fmt.Sprintf("/v1/notes/%s", uuid)
	res, err := utils.DoAuthorizedReq(ctx, hc, "DELETE", endpoint, "")
	if err != nil {
		return DeleteNoteResp{}, errors.Wrap(err, "patching a note to the server")
	}

	hasErr, message, err := checkRespErr(res)
	if err != nil {
		return DeleteNoteResp{}, errors.Wrap(err, "checking repsonse error")
	}
	if hasErr {
		return DeleteNoteResp{}, errors.New(message)
	}

	var resp DeleteNoteResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return DeleteNoteResp{}, errors.Wrap(err, "decoding payload")
	}

	return resp, nil
}

// GetBooksResp is a response from get books endpoint
type GetBooksResp []struct {
	UUID  string `json:"uuid"`
	Label string `json:"label"`
}

// GetBooks gets books from the server
func GetBooks(ctx infra.DnoteCtx, sessionKey string) (GetBooksResp, error) {
	hc := http.Client{}
	res, err := utils.DoAuthorizedReq(ctx, hc, "GET", "/v1/books", "")
	if err != nil {
		return GetBooksResp{}, errors.Wrap(err, "making http request")
	}

	hasErr, message, err := checkRespErr(res)
	if err != nil {
		return GetBooksResp{}, errors.Wrap(err, "checking repsonse error")
	}
	if hasErr {
		return GetBooksResp{}, errors.New(message)
	}

	var resp GetBooksResp
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return GetBooksResp{}, errors.Wrap(err, "decoding payload")
	}

	return resp, nil
}

// PresigninResponse is a reponse from /v1/presignin endpoint
type PresigninResponse struct {
	Iteration int `json:"iteration"`
}

// GetPresignin gets presignin credentials
func GetPresignin(ctx infra.DnoteCtx, email string) (PresigninResponse, error) {
	res, err := utils.DoReq(ctx, "GET", fmt.Sprintf("/v1/presignin?email=%s", email), "")
	if err != nil {
		return PresigninResponse{}, errors.Wrap(err, "making http request")
	}

	hasErr, message, err := checkRespErr(res)
	if err != nil {
		return PresigninResponse{}, errors.Wrap(err, "checking repsonse error")
	}
	if hasErr {
		return PresigninResponse{}, errors.New(message)
	}

	var resp PresigninResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return PresigninResponse{}, errors.Wrap(err, "decoding payload")
	}

	return resp, nil
}

// SigninPayload is a payload for /v1/signin
type SigninPayload struct {
	Email   string `json:"email"`
	AuthKey string `json:"auth_key"`
}

// SigninResponse is a response from /v1/signin endpoint
type SigninResponse struct {
	Key          string `json:"key"`
	ExpiresAt    int64  `json:"expires_at"`
	CipherKeyEnc string `json:"cipher_key_enc"`
}

// Signin requests a session token
func Signin(ctx infra.DnoteCtx, email, authKey string) (SigninResponse, error) {
	payload := SigninPayload{
		Email:   email,
		AuthKey: authKey,
	}
	b, err := json.Marshal(payload)
	if err != nil {
		return SigninResponse{}, errors.Wrap(err, "marshaling payload")
	}
	res, err := utils.DoReq(ctx, "POST", "/v1/signin", string(b))
	if err != nil {
		return SigninResponse{}, errors.Wrap(err, "making http request")
	}

	if res.StatusCode == http.StatusUnauthorized {
		return SigninResponse{}, ErrInvalidLogin
	}

	hasErr, message, err := checkRespErr(res)
	if err != nil {
		return SigninResponse{}, errors.Wrap(err, "checking repsonse error")
	}
	if hasErr {
		return SigninResponse{}, errors.New(message)
	}

	var resp SigninResponse
	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		return SigninResponse{}, errors.Wrap(err, "decoding payload")
	}

	return resp, nil
}

// Signout deletes a user session on the server side
func Signout(ctx infra.DnoteCtx, sessionKey string) error {
	hc := http.Client{
		// No need to follow redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := utils.DoAuthorizedReq(ctx, hc, "POST", "/v1/signout", "")
	if err != nil {
		return errors.Wrap(err, "making http request")
	}

	hasErr, message, err := checkRespErr(res)
	if err != nil {
		return errors.Wrap(err, "checking repsonse error")
	}
	if hasErr {
		return errors.New(message)
	}

	return nil
}
