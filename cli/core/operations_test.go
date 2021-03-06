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

package core

import (
	"fmt"
	"testing"

	"github.com/dnote/dnote/cli/testutils"
	"github.com/pkg/errors"
)

func TestInsertSystem(t *testing.T) {
	testCases := []struct {
		key string
		val string
	}{
		{
			key: "foo",
			val: "1558089284",
		},
		{
			key: "baz",
			val: "quz",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("insert %s %s", tc.key, tc.val), func(t *testing.T) {
			// Setup
			ctx := testutils.InitEnv(t, "../tmp", "../testutils/fixtures/schema.sql", true)
			defer testutils.TeardownEnv(ctx)

			// execute
			db := ctx.DB

			tx, err := db.Begin()
			if err != nil {
				t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
			}

			if err := InsertSystem(tx, tc.key, tc.val); err != nil {
				tx.Rollback()
				t.Fatalf(errors.Wrap(err, "executing for test case").Error())
			}

			tx.Commit()

			// test
			var key, val string
			testutils.MustScan(t, "getting the saved record",
				db.QueryRow("SELECT key, value FROM system WHERE key = ?", tc.key), &key, &val)

			testutils.AssertEqual(t, key, tc.key, "key mismatch for test case")
			testutils.AssertEqual(t, val, tc.val, "val mismatch for test case")
		})
	}
}

func TestUpsertSystem(t *testing.T) {
	testCases := []struct {
		key        string
		val        string
		countDelta int
	}{
		{
			key:        "foo",
			val:        "1558089284",
			countDelta: 1,
		},
		{
			key:        "baz",
			val:        "quz2",
			countDelta: 0,
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("insert %s %s", tc.key, tc.val), func(t *testing.T) {
			// Setup
			ctx := testutils.InitEnv(t, "../tmp", "../testutils/fixtures/schema.sql", true)
			defer testutils.TeardownEnv(ctx)

			db := ctx.DB
			testutils.MustExec(t, "inserting a system configuration", db, "INSERT INTO system (key, value) VALUES (?, ?)", "baz", "quz")

			var initialSystemCount int
			testutils.MustScan(t, "counting records", db.QueryRow("SELECT count(*) FROM system"), &initialSystemCount)

			// execute
			tx, err := db.Begin()
			if err != nil {
				t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
			}

			if err := UpsertSystem(tx, tc.key, tc.val); err != nil {
				tx.Rollback()
				t.Fatalf(errors.Wrap(err, "executing for test case").Error())
			}

			tx.Commit()

			// test
			var key, val string
			testutils.MustScan(t, "getting the saved record",
				db.QueryRow("SELECT key, value FROM system WHERE key = ?", tc.key), &key, &val)
			var systemCount int
			testutils.MustScan(t, "counting records",
				db.QueryRow("SELECT count(*) FROM system"), &systemCount)

			testutils.AssertEqual(t, key, tc.key, "key mismatch")
			testutils.AssertEqual(t, val, tc.val, "val mismatch")
			testutils.AssertEqual(t, systemCount, initialSystemCount+tc.countDelta, "count mismatch")
		})
	}
}

func TestGetSystem(t *testing.T) {
	t.Run(fmt.Sprintf("get string value"), func(t *testing.T) {
		// Setup
		ctx := testutils.InitEnv(t, "../tmp", "../testutils/fixtures/schema.sql", true)
		defer testutils.TeardownEnv(ctx)

		// execute
		db := ctx.DB
		testutils.MustExec(t, "inserting a system configuration", db, "INSERT INTO system (key, value) VALUES (?, ?)", "foo", "bar")

		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
		}
		var dest string
		if err := GetSystem(tx, "foo", &dest); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing for test case").Error())
		}
		tx.Commit()

		// test
		testutils.AssertEqual(t, dest, "bar", "dest mismatch")
	})

	t.Run(fmt.Sprintf("get int64 value"), func(t *testing.T) {
		// Setup
		ctx := testutils.InitEnv(t, "../tmp", "../testutils/fixtures/schema.sql", true)
		defer testutils.TeardownEnv(ctx)

		// execute
		db := ctx.DB
		testutils.MustExec(t, "inserting a system configuration", db, "INSERT INTO system (key, value) VALUES (?, ?)", "foo", 1234)

		tx, err := db.Begin()
		if err != nil {
			t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
		}
		var dest int64
		if err := GetSystem(tx, "foo", &dest); err != nil {
			tx.Rollback()
			t.Fatalf(errors.Wrap(err, "executing for test case").Error())
		}
		tx.Commit()

		// test
		testutils.AssertEqual(t, dest, int64(1234), "dest mismatch")
	})
}

func TestUpdateSystem(t *testing.T) {
	testCases := []struct {
		key        string
		val        string
		countDelta int
	}{
		{
			key: "foo",
			val: "1558089284",
		},
		{
			key: "foo",
			val: "bar",
		},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("update %s %s", tc.key, tc.val), func(t *testing.T) {
			// Setup
			ctx := testutils.InitEnv(t, "../tmp", "../testutils/fixtures/schema.sql", true)
			defer testutils.TeardownEnv(ctx)

			db := ctx.DB
			testutils.MustExec(t, "inserting a system configuration", db, "INSERT INTO system (key, value) VALUES (?, ?)", "foo", "fuz")
			testutils.MustExec(t, "inserting a system configuration", db, "INSERT INTO system (key, value) VALUES (?, ?)", "baz", "quz")

			var initialSystemCount int
			testutils.MustScan(t, "counting records", db.QueryRow("SELECT count(*) FROM system"), &initialSystemCount)

			// execute
			tx, err := db.Begin()
			if err != nil {
				t.Fatalf(errors.Wrap(err, "beginning a transaction").Error())
			}

			if err := UpdateSystem(tx, tc.key, tc.val); err != nil {
				tx.Rollback()
				t.Fatalf(errors.Wrap(err, "executing for test case").Error())
			}

			tx.Commit()

			// test
			var key, val string
			testutils.MustScan(t, "getting the saved record",
				db.QueryRow("SELECT key, value FROM system WHERE key = ?", tc.key), &key, &val)
			var systemCount int
			testutils.MustScan(t, "counting records",
				db.QueryRow("SELECT count(*) FROM system"), &systemCount)

			testutils.AssertEqual(t, key, tc.key, "key mismatch")
			testutils.AssertEqual(t, val, tc.val, "val mismatch")
			testutils.AssertEqual(t, systemCount, initialSystemCount, "count mismatch")
		})
	}
}
