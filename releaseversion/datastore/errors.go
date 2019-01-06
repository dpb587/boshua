package datastore

import (
	"errors"
	"fmt"
)

var UnsupportedOperationErr = errors.New("unsupported operation")

type UnexpectedCountComparisonErr struct {
	expectedOperator string
	expectedCount    int
	actualCount      int
}

var _ error = UnexpectedCountComparisonErr{}

func (e UnexpectedCountComparisonErr) Error() string {
	return fmt.Sprintf("expected %s%d results, but found %d", e.expectedOperator, e.expectedCount, e.actualCount)
}

func NewUnexpectedMinCountError(expected, actual int) UnexpectedCountComparisonErr {
	return UnexpectedCountComparisonErr{
		expectedOperator: ">=",
		expectedCount:    expected,
		actualCount:      actual,
	}
}

func NewUnexpectedMaxCountError(expected, actual int) UnexpectedCountComparisonErr {
	return UnexpectedCountComparisonErr{
		expectedOperator: "<=",
		expectedCount:    expected,
		actualCount:      actual,
	}
}
