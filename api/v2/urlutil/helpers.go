package urlutil

import (
	"errors"
	"fmt"
	"net/http"
)

func simpleQueryLookup(r *http.Request, param string) (string, error) {
	paramValue, ok := r.URL.Query()[param]
	if !ok {
		return "", fmt.Errorf("parameter '%s': %v", param, ParamMissingError)
	} else if len(paramValue) != 1 {
		return "", fmt.Errorf("parameter '%s': %v", param, fmt.Errorf("expected 1 value, but found %d", len(paramValue)))
	} else if len(paramValue[0]) == 0 {
		return "", fmt.Errorf("parameter '%s': %v", param, errors.New("expected non-empty value"))
	}

	return paramValue[0], nil
}
