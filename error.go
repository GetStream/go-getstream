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
		Duration        string                 `json:"duration"`
		ExceptionFields map[string]interface{} `json:"exception_fields"`
		*alias
	}{alias: (*alias)(e)}

	if err := json.Unmarshal(b, &aux); err != nil {
		return err
	}

	dur, err := time.ParseDuration(aux.Duration)
	if err != nil {
		return err
	}
	e.Duration = dur
	e.ExceptionFields = makeExceptionFields(aux.ExceptionFields)
	return nil
}

func makeExceptionFields(data map[string]interface{}) map[string][]string {
	if len(data) == 0 {
		return nil
	}

	exceptionFields := make(map[string][]string)
	for k, v := range data {
		slice, ok := v.([]interface{})
		if !ok {
			continue
		}
		exceptionFields[k] = make([]string, len(slice))
		i := 0
		for _, elem := range slice {
			if s, ok := elem.(string); ok {
				exceptionFields[k][i] = s
				i++
			}
		}
		exceptionFields[k] = exceptionFields[k][:i]
	}
	return exceptionFields
}
