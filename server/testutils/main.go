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

// Package testutils provides utilities used in tests
package testutils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/dnote/dnote/server/database"
	"github.com/stripe/stripe-go"

	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

func checkEqual(a interface{}, b interface{}, message string) (bool, string) {
	if a == b {
		return true, ""
	}

	var m string
	if len(message) == 0 {
		m = fmt.Sprintf("%v != %v", a, b)
	} else {
		m = message
	}
	errorMessage := fmt.Sprintf("%s. Actual: %+v. Expected: %+v.", m, a, b)

	return false, errorMessage
}

// AssertEqual errors a test if the actual does not match the expected
func AssertEqual(t *testing.T, a interface{}, b interface{}, message string) {
	ok, m := checkEqual(a, b, message)
	if !ok {
		t.Error(m)
	}
}

// AssertEqualf fails a test if the actual does not match the expected
func AssertEqualf(t *testing.T, a interface{}, b interface{}, message string) {
	ok, m := checkEqual(a, b, message)
	if !ok {
		t.Fatal(m)
	}
}

// AssertNotEqual fails a test if the actual matches the expected
func AssertNotEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if a != b {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Errorf("%s. Expected %+v to not equal %+v.", message, a, b)
}

// AssertDeepEqual fails a test if the actual does not deeply equal the expected
func AssertDeepEqual(t *testing.T, a interface{}, b interface{}, message string) {
	if reflect.DeepEqual(a, b) {
		return
	}

	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Errorf("%s.\nActual:   %+v.\nExpected: %+v.", message, a, b)
}

// AssertEqualJSON asserts that two JSON strings are equal
func AssertEqualJSON(t *testing.T, a, b, message string) {
	var o1 interface{}
	var o2 interface{}

	err := json.Unmarshal([]byte(a), &o1)
	if err != nil {
		panic(fmt.Errorf("Error mashalling string 1 :: %s", err.Error()))
	}
	err = json.Unmarshal([]byte(b), &o2)
	if err != nil {
		panic(fmt.Errorf("Error mashalling string 2 :: %s", err.Error()))
	}

	if reflect.DeepEqual(o1, o2) {
		return
	}

	if len(message) == 0 {
		message = fmt.Sprintf("%v != %v", a, b)
	}
	t.Errorf("%s.\nActual:   %+v.\nExpected: %+v.", message, a, b)
}

// ReadJSON reads JSON fixture to the struct at the destination address
func ReadJSON(path string, destination interface{}) {
	var dat []byte
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(errors.Wrap(err, "reading file"))
	}
	if err := json.Unmarshal(dat, destination); err != nil {
		panic(errors.Wrap(err, "unmarshalling json"))
	}
}

// ReadFile reads file and returns the byte
func ReadFile(path string) []byte {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(errors.Wrap(err, "reading file"))
	}

	return dat
}

// InitTestDB establishes connection pool with the test database specified by
// the environment variable configuration and initalizes a new schema
func InitTestDB() {
	database.InitDB()
	database.InitSchema()
}

// SetupUserData creates and returns a new user for testing purposes
func SetupUserData() database.User {
	db := database.DBConn

	user := database.User{
		APIKey: "test-api-key",
		Name:   "user-name",
		Cloud:  true,
	}

	if err := db.Save(&user).Error; err != nil {
		panic(errors.Wrap(err, "Failed to prepare user"))
	}

	return user
}

// SetupAccountData creates and returns a new account for the user
func SetupAccountData(user database.User, email string) database.Account {
	db := database.DBConn

	// email: alice@example.com
	// password: pass1234
	// masterKey: WbUvagj9O6o1Z+4+7COjo7Uqm4MD2QE9EWFXne8+U+8=
	// authKey: /XCYisXJ6/o+vf6NUEtmrdYzJYPz+T9oAUCtMpOjhzc=
	account := database.Account{
		UserID:             user.ID,
		Salt:               "Et0joOigYjdgHBKMN/ijxg==",
		AuthKeyHash:        "SeN3PMz4H/7q9lINB+VPKpygexAuK68wO8pDAgQ4OOQ=",
		CipherKeyEnc:       "f7aFFCh7YS1WlHEOxAmDfs8rUQQoX5tr8AB7ZJQaTYCEM8NhAZCbQTsjFgKOf5iPQhhkm8eDAgPNTuhO",
		ClientKDFIteration: 100000,
		ServerKDFIteration: 100000,
	}
	if email != "" {
		account.Email = database.ToNullString(email)
	}

	if err := db.Save(&account).Error; err != nil {
		panic(errors.Wrap(err, "Failed to prepare account"))
	}

	return account
}

// SetupSession creates and returns a new user session
func SetupSession(t *testing.T, user database.User) database.Session {
	db := database.DBConn

	session := database.Session{
		Key:       "Vvgm3eBXfXGEFWERI7faiRJ3DAzJw+7DdT9J1LEyNfI=",
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 24),
	}
	if err := db.Save(&session).Error; err != nil {
		t.Fatal(errors.Wrap(err, "Failed to prepare user"))
	}

	return session
}

// SetupEmailPreferenceData creates and returns a new email frequency for a user
func SetupEmailPreferenceData(user database.User, digestWeekly bool) database.EmailPreference {
	db := database.DBConn

	frequency := database.EmailPreference{
		UserID:       user.ID,
		DigestWeekly: digestWeekly,
	}

	if err := db.Save(&frequency).Error; err != nil {
		panic(errors.Wrap(err, "Failed to prepare email frequency"))
	}

	return frequency
}

// ClearData deletes all records from the database
func ClearData() {
	db := database.DBConn

	if err := db.Delete(&database.Book{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear books"))
	}
	if err := db.Delete(&database.Note{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear notes"))
	}
	if err := db.Delete(&database.Notification{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear notifications"))
	}
	if err := db.Delete(&database.User{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear users"))
	}
	if err := db.Delete(&database.Account{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear accounts"))
	}
	if err := db.Delete(&database.Token{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear reset_tokens"))
	}
	if err := db.Delete(&database.EmailPreference{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear reset_tokens"))
	}
	if err := db.Delete(&database.Session{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear sessions"))
	}
	if err := db.Delete(&database.Digest{}).Error; err != nil {
		panic(errors.Wrap(err, "Failed to clear digests"))
	}
}

// HTTPDo makes an HTTP request and returns a response
func HTTPDo(t *testing.T, req *http.Request) *http.Response {
	hc := http.Client{
		// Do not follow redirects.
		// e.g. /logout redirects to a page but we'd like to test the redirect
		// itself, not what happens after the redirect
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	res, err := hc.Do(req)
	if err != nil {
		t.Fatal(errors.Wrap(err, "performing http request"))
	}

	return res
}

// HTTPAuthDo makes an HTTP request with an appropriate authorization header for a user
func HTTPAuthDo(t *testing.T, req *http.Request, user database.User) *http.Response {
	db := database.DBConn

	session := database.Session{
		Key:       "Vvgm3eBXfXGEFWERI7faiRJ3DAzJw+7DdT9J1LEyNfI=",
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(time.Hour * 10 * 24),
	}
	if err := db.Save(&session).Error; err != nil {
		t.Fatal(errors.Wrap(err, "Failed to prepare user"))
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", session.Key))

	return HTTPDo(t, req)

}

// MakeReq makes an HTTP request and returns a response
func MakeReq(server *httptest.Server, method, url, data string) *http.Request {
	endpoint := fmt.Sprintf("%s%s", server.URL, url)

	req, err := http.NewRequest(method, endpoint, strings.NewReader(data))
	if err != nil {
		panic(errors.Wrap(err, "constructing http request"))
	}

	return req
}

// AssertStatusCode asserts that the reponse's status code is equal to the
// expected
func AssertStatusCode(t *testing.T, res *http.Response, expected int, message string) {
	if res.StatusCode != expected {
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			panic(errors.Wrap(err, "reading body"))
		}

		t.Errorf("status code mismatch. %s: got %v want %v. Message was: '%s'", message, res.StatusCode, expected, string(body))
	}
}

// MustExec fails the test if the given database query has error
func MustExec(t *testing.T, db *gorm.DB, message string) {
	if err := db.Error; err != nil {
		t.Fatalf("%s: %s", message, err.Error())
	}
}

// MustMarshalJSON marshalls the given interface into JSON.
// If there is any error, it fails the test.
func MustMarshalJSON(t *testing.T, v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal("marshalling data")
	}

	return b
}

// GetCookieByName returns a cookie with the given name
func GetCookieByName(cookies []*http.Cookie, name string) *http.Cookie {
	var ret *http.Cookie

	for i := 0; i < len(cookies); i++ {
		if cookies[i].Name == name {
			ret = cookies[i]
			break
		}
	}

	return ret
}

// CreateMockStripeBackend returns a mock stripe backend implementation that uses
// the given test server
func CreateMockStripeBackend(ts *httptest.Server) *stripe.BackendImplementation {
	c := ts.Client()
	bi := stripe.BackendImplementation{
		Type:       stripe.APIBackend,
		URL:        ts.URL + "/v1",
		HTTPClient: c,
	}

	return &bi
}
