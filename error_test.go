package getstream_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	getstream "github.com/GetStream/stream-go"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {

	errorResponse := "{\"code\": 5, \"detail\": \"some detail\", \"duration\": \"36ms\", \"exception\": \"an exception\", \"status_code\": 400}"

	var getStreamError getstream.Error
	err := json.Unmarshal([]byte(errorResponse), &getStreamError)
	if err != nil {
		t.Fatal(err)
	}

	testError := getstream.Error{
		Code:       5,
		Detail:     "some detail",
		Duration:   36 * time.Millisecond,
		Exception:  "an exception",
		StatusCode: 400,
	}

	assert.Equal(t, getStreamError, testError)
	assert.Equal(t, getStreamError.Duration, time.Millisecond*36)
	assert.Equal(t, getStreamError.Error(), "an exception (36ms): some detail")
}

func TestErrorBadDuration(t *testing.T) {
	var e getstream.Error

	testCases := []struct {
		duration string
	}{
		{duration: ""},
		{duration: "asd"},
		{duration: "123asd"},
	}

	for _, tc := range testCases {
		payload := fmt.Sprintf(`{"code":42, "duration":"%s"}`, tc.duration)
		err := json.Unmarshal([]byte(payload), &e)
		assert.NotNil(t, err)
	}
}
