package datastore

import "errors"

var UnsupportedOperationErr = errors.New("unsupported operation")
var NoMatchErr = errors.New("compiled release version not found")
var MultipleMatchErr = errors.New("multiple compiled release versions found")
