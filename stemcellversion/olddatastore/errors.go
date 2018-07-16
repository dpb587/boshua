package datastore

import "errors"

var NoMatchErr = errors.New("stemcell version not found")
var MultipleMatchErr = errors.New("multiple stemcell versions found")
var UnsupportedOperationErr = errors.New("unsupported operation")
