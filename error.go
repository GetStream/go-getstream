package getstream

import (
	"time"
)

// Credits to https://github.com/hyperworks/go-getstream for the error handling.

// Error is a getstream error
type Error struct {
	Code       int `json:"code"`
	StatusCode int `json:"status_code"`

	Detail          string              `json:"detail"`
	RawDuration     string              `json:"duration"`
	Exception       string              `json:"exception"`
	ExceptionFields map[string][]string `json:"exception_fields"`
}

var _ error = &Error{}

// Duration is the time it took for the request to be handled
func (e *Error) Duration() time.Duration {
	result, err := time.ParseDuration(e.RawDuration)
	if err != nil {
		return time.Duration(0)
	}

	return result
}

func (e *Error) Error() string {
	str := e.Exception
	if e.RawDuration != "" {
		if duration := e.Duration(); duration > 0 {
			str += " (" + duration.String() + ")"
		}
	}

	if e.Detail != "" {
		str += ": " + e.Detail
	}

	return str
}
