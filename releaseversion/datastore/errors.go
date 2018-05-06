package datastore

import "errors"

var NoMatchErr = errors.New("release version not found")
var MultipleMatchErr = errors.New("multiple release versions found")
