package datastore

import "errors"

var NoMatchErr = errors.New("analysis not found")
var MultipleMatchErr = errors.New("multiple analyses found")
