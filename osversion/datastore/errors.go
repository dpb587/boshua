package datastore

import "errors"

var NoMatchErr = errors.New("os version not found")
var MultipleMatchErr = errors.New("multiple os versions found")
