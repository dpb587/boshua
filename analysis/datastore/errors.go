package datastore

import "errors"

var UnsupportedOperationErr = errors.New("unsupported operation")
var NoMatchErr = errors.New("analysis not found")
var MultipleMatchErr = errors.New("multiple analyses found")
