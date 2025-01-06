package posgres

import "errors"

var ErrNotFound = errors.New("not found")
var ErrDatabase = errors.New("database error")
var ErrDuplicate = errors.New("already exists")
