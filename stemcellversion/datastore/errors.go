package datastore

import "errors"

var NoMatchErr = errors.New("stemcell version not found")
var MultipleMatchErr = errors.New("multiple stemcell versions found")
