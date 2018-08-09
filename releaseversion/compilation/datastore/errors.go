package datastore

import "errors"

var UnsupportedOperationErr = errors.New("unsupported operation")
var NoMatchErr = errors.New("no match found")
var MultipleMatchErr = errors.New("multiple matches found")
