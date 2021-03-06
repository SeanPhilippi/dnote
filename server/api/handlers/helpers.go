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
	crand "crypto/rand"
	"encoding/base64"
	"strings"

	"github.com/dnote/dnote/server/database"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

const (
	demoUserEmail = "demo@dnote.io"
)

func generateRandomToken(bits int) (string, error) {
	b := make([]byte, bits)

	_, err := crand.Read(b)
	if err != nil {
		return "", errors.Wrap(err, "generating random bytes")
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

func generateResetToken() (string, error) {
	ret, err := generateRandomToken(16)
	if err != nil {
		return "", errors.Wrap(err, "generating random token")
	}

	return ret, nil
}

func generateVerificationCode() (string, error) {
	ret, err := generateRandomToken(16)
	if err != nil {
		return "", errors.Wrap(err, "generating random token")
	}

	return ret, nil
}

func paginate(conn *gorm.DB, page int) *gorm.DB {
	limit := 30

	// Paginate
	if page > 0 {
		offset := limit * (page - 1)
		conn = conn.Offset(offset)
	}

	conn = conn.Limit(limit)

	return conn
}

func getBookIDs(books []database.Book) []int {
	ret := []int{}

	for _, book := range books {
		ret = append(ret, book.ID)
	}

	return ret
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return errors.New("Password should be longer than 8 characters")
	}

	return nil
}

func getClientType(origin string) string {
	if strings.HasPrefix(origin, "moz-extension://") {
		return "firefox-extension"
	}

	if strings.HasPrefix(origin, "chrome-extension://") {
		return "chrome-extension"
	}

	return "web"
}
