package core

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/dnote/actions"
	"github.com/dnote/cli/infra"
	"github.com/dnote/cli/log"
	"github.com/dnote/cli/utils"
	"github.com/pkg/errors"
)

// RewindAction rewinds the given action according to its type
func RewindAction(ctx infra.DnoteCtx, tx *sql.Tx, action actions.Action) error {
	var err error

	switch action.Type {
	case actions.ActionAddBook:
		err = rewindAddBook(ctx, tx, action)
		//	case actions.ActionRemoveBook:
		//		err = rewindRemoveBook(ctx, tx, action)
		//	case actions.ActionAddNote:
		//		err = rewindAddNote(ctx, tx, action)
		//	case actions.ActionRemoveNote:
		//		err = rewindRemoveNote(ctx, tx, action)
		//	case actions.ActionEditNote:
		//		err = rewindEditNote(ctx, tx, action)
	default:
		return errors.Wrapf(err, "Unsupported action type: %s", action.Type)
	}

	if err != nil {
		return errors.Wrapf(err, "rewinding an action %s", action.Type)
	}

	return nil
}

// rewindAddBook rewinds the given add_book action
func rewindAddBook(ctx infra.DnoteCtx, tx *sql.Tx, action actions.Action) error {
	if action.Schema != 1 {
		return errors.Errorf("unsupported schema %d", action.Schema)
	}

	var data actions.AddBookDataV1
	if err := json.Unmarshal(action.Data, &data); err != nil {
		return errors.Wrap(err, "parsing action data")
	}

	if _, err := tx.Exec("DELETE FROM books WHERE label = ?", data.BookName); err != nil {
		return errors.Wrap(err, "removing book")
	}

	return nil
}

// rewindRemoveBook rewinds the given remove_book action
func rewindRemoveBook(ctx infra.DnoteCtx, tx *sql.Tx, action actions.Action) error {
	if action.Schema != 1 {
		return errors.Errorf("unsupported schema %d", action.Schema)
	}

	var data actions.RemoveBookDataV1
	if err := json.Unmarshal(action.Data, &data); err != nil {
		return errors.Wrap(err, "parsing action data")
	}

	return nil
}

// PlayAction transitions the local dnote state by consuming the action returned
// from the server
func PlayAction(ctx infra.DnoteCtx, tx *sql.Tx, action actions.Action) error {
	log.Debug("playing %s. uuid: %s, schema: %d, timestamp: %d\n", action.Type, action.UUID, action.Schema, action.Timestamp)

	var err error

	switch action.Type {
	case actions.ActionAddNote:
		err = playAddNote(ctx, tx, action)
	case actions.ActionRemoveNote:
		err = playRemoveNote(ctx, tx, action)
	case actions.ActionEditNote:
		err = playEditNote(ctx, tx, action)
	case actions.ActionAddBook:
		err = playAddBook(ctx, tx, action)
	case actions.ActionRemoveBook:
		err = playRemoveBook(ctx, tx, action)
	default:
		return errors.Errorf("Unsupported action %s", action.Type)
	}

	if err != nil {
		return errors.Wrapf(err, "reducing %s", action.Type)
	}

	return nil
}

func getBookUUIDWithTx(tx *sql.Tx, bookLabel string) (string, error) {
	var ret string
	err := tx.QueryRow("SELECT uuid FROM books WHERE label = ?", bookLabel).Scan(&ret)
	if err == sql.ErrNoRows {
		return ret, errors.Errorf("book '%s' not found", bookLabel)
	} else if err != nil {
		return ret, errors.Wrap(err, "querying the book")
	}

	return ret, nil
}

// AddNote adds a note and logs the action
func AddNote(tx *sql.Tx, bookUUID, bookLabel, content string, ts int64) error {
	uuid, err := performAddNote(tx, bookUUID, bookLabel, content, ts)
	if err != nil {
		return errors.Wrap(err, "adding note")
	}

	err = LogActionAddNote(tx, uuid, bookLabel, content, ts)
	if err != nil {
		return errors.Wrap(err, "logging action")
	}

	return nil
}

func performAddNote(tx *sql.Tx, bookUUID, bookLabel, content string, ts int64) (string, error) {
	uuid := utils.GenerateUUID()

	_, err := tx.Exec(`INSERT INTO notes (uuid, book_uuid, content, added_on, public)
		VALUES (?, ?, ?, ?, ?);`, uuid, bookUUID, content, ts, false)
	if err != nil {
		return uuid, errors.Wrap(err, "creating the note")
	}

	return uuid, nil
}

func playAddNote(ctx infra.DnoteCtx, tx *sql.Tx, action actions.Action) error {
	if action.Schema != 2 {
		return errors.Errorf("data schema '%d' not supported", action.Schema)
	}

	var data actions.AddNoteDataV2
	if err := json.Unmarshal(action.Data, &data); err != nil {
		return errors.Wrap(err, "parsing the action data")
	}

	log.Debug("data: %+v\n", data)

	bookUUID, err := getBookUUIDWithTx(tx, data.BookName)
	if err != nil {
		return errors.Wrap(err, "getting book uuid")
	}

	var noteCount int
	if err := tx.
		QueryRow("SELECT count(uuid) FROM notes WHERE uuid = ? AND book_uuid = ?", data.NoteUUID, bookUUID).
		Scan(&noteCount); err != nil {
		return errors.Wrap(err, "counting note")
	}

	if noteCount > 0 {
		// if a duplicate exists, it is because the same action has been previously synced to the server
		// but the client did not bring the bookmark up-to-date at the time because it had error reducing
		// the returned actions.
		// noop so that the client can update bookmark
		return nil
	}

	_, err = tx.Exec(`INSERT INTO notes
	(uuid, book_uuid, content, added_on, public)
	VALUES (?, ?, ?, ?, ?)`, data.NoteUUID, bookUUID, data.Content, action.Timestamp, data.Public)
	if err != nil {
		return errors.Wrap(err, "inserting a note")
	}

	return nil
}

func playRemoveNote(ctx infra.DnoteCtx, tx *sql.Tx, action actions.Action) error {
	if action.Schema != 2 {
		return errors.Errorf("data schema '%d' not supported", action.Schema)
	}

	var data actions.RemoveNoteDataV2
	if err := json.Unmarshal(action.Data, &data); err != nil {
		return errors.Wrap(err, "parsing the action data")
	}

	log.Debug("data: %+v\n", data)

	_, err := tx.Exec("DELETE FROM notes WHERE uuid = ?", data.NoteUUID)
	if err != nil {
		return errors.Wrap(err, "removing a note")
	}

	return nil
}

func buildEditNoteQuery(ctx infra.DnoteCtx, tx *sql.Tx, noteUUID string, ts int64, data actions.EditNoteDataV3) (string, []interface{}, error) {
	setTmpl := "edited_on = ?"
	queryArgs := []interface{}{ts}

	if data.Content != nil {
		setTmpl = fmt.Sprintf("%s, content = ?", setTmpl)
		queryArgs = append(queryArgs, *data.Content)
	}
	if data.Public != nil {
		setTmpl = fmt.Sprintf("%s, public = ?", setTmpl)
		queryArgs = append(queryArgs, *data.Public)
	}
	if data.BookName != nil {
		setTmpl = fmt.Sprintf("%s, book_uuid = ?", setTmpl)

		bookUUID, err := getBookUUIDWithTx(tx, *data.BookName)
		if err != nil {
			return setTmpl, queryArgs, errors.Wrap(err, "getting book uuid")
		}

		queryArgs = append(queryArgs, bookUUID)
	}

	queryTmpl := fmt.Sprintf("UPDATE notes SET %s WHERE uuid = ?", setTmpl)
	queryArgs = append(queryArgs, noteUUID)

	return queryTmpl, queryArgs, nil
}

func playEditNote(ctx infra.DnoteCtx, tx *sql.Tx, action actions.Action) error {
	if action.Schema != 3 {
		return errors.Errorf("data schema '%d' not supported", action.Schema)
	}

	var data actions.EditNoteDataV3
	err := json.Unmarshal(action.Data, &data)
	if err != nil {
		return errors.Wrap(err, "parsing the action data")
	}

	log.Debug("data: %+v\n", data)

	queryTmpl, queryArgs, err := buildEditNoteQuery(ctx, tx, data.NoteUUID, action.Timestamp, data)
	if err != nil {
		return errors.Wrap(err, "building edit note query")
	}
	_, err = tx.Exec(queryTmpl, queryArgs...)
	if err != nil {
		return errors.Wrap(err, "updating a note")
	}

	return nil
}

// HandleAddBook adds a new book and logs the action
func HandleAddBook(tx *sql.Tx, bookLabel string) (string, error) {
	uuid, err := performAddBook(tx, bookLabel)
	if err != nil {
		return uuid, errors.Wrap(err, "performing add_book")
	}

	if err := LogActionAddBook(tx, bookLabel); err != nil {
		tx.Rollback()
		return uuid, errors.Wrap(err, "logging action")
	}

	return uuid, nil
}

// performAddBook encapsulates the actual low-level details of adding a book
func performAddBook(tx *sql.Tx, bookLabel string) (string, error) {
	uuid := utils.GenerateUUID()

	if _, err := tx.Exec("INSERT INTO books (uuid, label) VALUES (?, ?)", uuid, bookLabel); err != nil {
		tx.Rollback()
		return uuid, errors.Wrap(err, "creating the book")
	}

	return uuid, nil
}

// playAddBook checks if the given add_book action needs to be played and plays it
func playAddBook(ctx infra.DnoteCtx, tx *sql.Tx, action actions.Action) error {
	if action.Schema != 1 {
		return errors.Errorf("data schema '%d' not supported", action.Schema)
	}

	var data actions.AddBookDataV1
	err := json.Unmarshal(action.Data, &data)
	if err != nil {
		return errors.Wrap(err, "parsing the action data")
	}

	log.Debug("data: %+v\n", data)

	var bookCount int
	err = tx.QueryRow("SELECT count(uuid) FROM books WHERE label = ?", data.BookName).Scan(&bookCount)
	if err != nil {
		return errors.Wrap(err, "counting books")
	}
	if bookCount > 0 {
		// If book already exists, bookmark was not updated for any error
		// noop
		return nil
	}

	if _, err = performAddBook(tx, data.BookName); err != nil {
		return errors.Wrap(err, "inserting a book")
	}

	return nil
}

func playRemoveBook(ctx infra.DnoteCtx, tx *sql.Tx, action actions.Action) error {
	if action.Schema != 1 {
		return errors.Errorf("data schema '%d' not supported", action.Schema)
	}

	var data actions.RemoveBookDataV1
	if err := json.Unmarshal(action.Data, &data); err != nil {
		return errors.Wrap(err, "parsing the action data")
	}

	log.Debug("data: %+v\n", data)

	var bookCount int
	if err := tx.
		QueryRow("SELECT count(uuid) FROM books WHERE label = ?", data.BookName).
		Scan(&bookCount); err != nil {
		return errors.Wrap(err, "counting note")
	}

	if bookCount == 0 {
		// If book does not exist, another client added and removed the book, making the add_book action
		// obsolete. noop.
		return nil
	}

	bookUUID, err := getBookUUIDWithTx(tx, data.BookName)
	if err != nil {
		return errors.Wrap(err, "getting book uuid")
	}

	_, err = tx.Exec("DELETE FROM notes WHERE book_uuid = ?", bookUUID)
	if err != nil {
		return errors.Wrap(err, "removing notes")
	}

	_, err = tx.Exec("DELETE FROM books WHERE uuid = ?", bookUUID)
	if err != nil {
		return errors.Wrap(err, "removing a book")
	}

	return nil
}
