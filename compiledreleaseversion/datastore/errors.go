package datastore

import "errors"

var NoMatchErr = errors.New("compiled release version not found")
var MultipleMatchErr = errors.New("multiple compiled release versions found")
