# Go Session Authentication
[![Build Status](https://travis-ci.org/apexskier/httpauth.svg?branch=master)](https://travis-ci.org/apexskier/httpauth)
[![GoDoc](https://godoc.org/github.com/apexskier/httpauth?status.png)](https://godoc.org/github.com/apexskier/httpauth)

This package uses the [Gorilla web toolkit](http://www.gorillatoolkit.org/)'s
sessions and package to implement a user authentication and authorization
system for Go web servers.

Multiple user data storage backends are available, and new ones can be
implemented relatively easily.

- [File based](https://godoc.org/github.com/apexskier/goauth#NewGobFileAuthBackend) ([gob](http://golang.org/pkg/encoding/gob/))
- [Various SQL Databases](https://godoc.org/github.com/apexskier/httpauth#NewSqlAuthBackend)

Access can be restricted by a users' role.

Uses [bcrypt](http://codahale.com/how-to-safely-store-a-password/) for password
hashing.

Run `go run server.go` from the examples directory and visit `localhost:8080`
for an example. You can login with the username and password "admin".

**Note**

This is the first time I've worked with implementing the details of cookie
storage, authentication or any sort of real security. There are no guarantees
that this will work as expected, but I'd love feedback. If you have any issues
or suggestions, please [let me
know](https://github.com/Wombats/goauth/issues/new).

### TODO

- User roles
- SMTP email validation (key based)
