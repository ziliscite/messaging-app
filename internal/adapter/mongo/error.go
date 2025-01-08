package mongo

import "errors"

var ErrNotFound = errors.New("document not found")
var ErrDecode = errors.New("failed to decode document")
