package getstream

import (
	"encoding/json"
	"fmt"
	"time"
)

// Credits to https://github.com/hyperworks/go-getstream for the error handling.

// Error is a getstream error
type Error struct {
	Code            int                 `json:"code"`
	StatusCode      int                 `json:"status_code"`
	Duration        time.Duration       `json:"duration"`
	Detail          string              `json:"detail"`
	Exception       string              `json:"exception"`
	ExceptionFields map[string][]string `json:"exception_fields"`
}

var _ error = &Error{}

func (e *Error) Error() string {
	str := fmt.Sprintf("%s (%s)", e.Exception, e.Duration)
	if e.Detail != "" {
		str += ": " + e.Detail
	}
	return str
}

func (e *Error) UnmarshalJSON(b []byte) error {
	type alias Error
	aux := &struct {
		Duration string `json:"duration"`
		*alias
	}{
		alias: (*alias)(e),
	}
	err := json.Unmarshal(b, &aux)
	if err != nil {
		return err
	}
	dur, err := time.ParseDuration(aux.Duration)
	if err != nil {
		return err
	}
	e.Duration = dur
	return nil
}
