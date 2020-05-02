package staffdb

import (
	"errors"
	"fmt"

	"github.com/lootch/httpauth"
	"proxima.alt.za/x/sql"
)

// StaffAuthBackend database and database connection information.
type StaffAuthBackend struct {
	DB        sql.Server
	// prepared statements
	userStmt   *sql.Stmt
	usersStmt  *sql.Stmt
	insertStmt *sql.Stmt
	updateStmt *sql.Stmt
	deleteStmt *sql.Stmt
}

func mksqlerror(msg string) error {
	return errors.New("sqlbackend: " + msg)
}

// NewStaffAuthBackend initializes a new backend by testing the database
// connection and making sure the storage table exists. The table is called
// ff_contacts.
//
// Returns an error if connecting to the database fails, pinging the database
// fails, or creating the table fails.
//
// This uses the databases/sql package to open a connection. Its parameters
// should match the sql.Open function. See
// http://golang.org/pkg/database/sql/#Open for more information.
//
// Be sure to import "database/sql" and your driver of choice. If you're not
// using sql for your own purposes, you'll need to use the underscore to import
// for side effects; see http://golang.org/doc/effective_go.html#blank_import.
func NewStaffAuthBackend(db sql.Server) (b StaffAuthBackend, err error) {
	b.DB = db
	// prepare statements for concurrent use and better preformance
	//
	b.userStmt, err = db.Prepare(
`select
	email, authdata as hash, 'user' as role
from ff_contacts
where username = ?`)
	if err != nil {
		return b, mksqlerror(fmt.Sprintf("userstmt: %v", err))
	}
	b.usersStmt, err = db.Prepare(
`select
	username, email, authdata as hash, 'user' as role
from ff_contacts`)
	if err != nil {
		return b, mksqlerror(fmt.Sprintf("usersstmt: %v", err))
	}
	b.insertStmt, err = db.Prepare(
`insert into ff_contacts
	(username, email, authdata)
values (?, ?, ?)`)
	if err != nil {
		return b, mksqlerror(fmt.Sprintf("insertstmt: %v", err))
	}
	b.updateStmt, err = db.Prepare(
`update goauth
set email = ?, authdata = ?
where Username = ?`)
	if err != nil {
		return b, mksqlerror(fmt.Sprintf("updatestmt: %v", err))
	}
	b.deleteStmt, err = db.Prepare(
`delete from ff_contacts
where username = ?`)
	if err != nil {
		return b, mksqlerror(fmt.Sprintf("deletestmt: %v", err))
	}
	return b, nil
}

// User returns the user with the given username. Error is set to
// ErrMissingUser if user is not found.
func (b StaffAuthBackend) User(username string) (user httpauth.UserData, e error) {
	switch err := b.userStmt.QueryRow(username).
		Scan(&user.Email, &user.Hash, &user.Role); {
	case err == sql.ErrNoRows:
		return user, httpauth.ErrMissingUser
	case err != nil:
		return user, mksqlerror(err.Error())
	}
	user.Username = username
	return user, nil
}

// Users returns a slice of all users.
func (b StaffAuthBackend) Users() (us []httpauth.UserData, e error) {
	var (
		username, email, role string
		hash                  []byte
	)

	rows, err := b.usersStmt.Query()
	if err != nil {
		return us, mksqlerror(err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&username, &email, &hash, &role)
		if err != nil {
			return us, mksqlerror(err.Error())
		}
		us = append(us, httpauth.UserData{username, email, hash, role})
	}
	return us, nil
}

// SaveUser adds a new user, replacing one with the same username.
func (b StaffAuthBackend) SaveUser(user httpauth.UserData) (err error) {
	if _, err = b.User(user.Username); err == nil {
		_, err = b.updateStmt.Exec(user.Email, user.Hash, user.Username)
	} else {
		_, err = b.insertStmt.Exec(user.Username, user.Email, user.Hash)
	}
	return
}

// DeleteUser removes a user, raising ErrDeleteNull if that user was missing.
func (b StaffAuthBackend) DeleteUser(username string) error {
	result, err := b.deleteStmt.Exec(username)
	if err != nil {
		return mksqlerror(err.Error())
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return mksqlerror(err.Error())
	}
	if rows == 0 {
		return httpauth.ErrDeleteNull
	}
	return nil
}

// Close cleans up the backend by terminating the database connection.
func (b StaffAuthBackend) Close() {
	b.userStmt.Close()
	b.usersStmt.Close()
	b.insertStmt.Close()
	b.updateStmt.Close()
	b.deleteStmt.Close()
}
