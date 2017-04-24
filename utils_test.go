package getstream_test

import (
	"testing"

	getstream "github.com/GetStream/stream-go"
)

func TestFeedSlug(t *testing.T) {
	var tests = []struct {
		//given
		description string
		slug        string

		//expected
		validatedSlug string
		errorMsg      string
	}{
		{
			"Slug with word chars",
			"foo",

			"foo",
			"",
		},
		{
			"Slug with word chars and dashes",
			"f-o-o",

			"f-o-o",
			"",
		},
		{
			"Slug with invalid chars",
			"foo@",

			"",
			"invalid feedSlug",
		},
	}

	for _, test := range tests {
		feedSlug, err := getstream.ValidateFeedSlug(test.slug)
		if feedSlug != test.validatedSlug {
			t.Errorf("Expected slug %v, got %v", test.validatedSlug, feedSlug)
		}

		if err != nil && err.Error() != test.errorMsg {
			t.Errorf("Error %v does not match expected %v", test.errorMsg, err.Error())
		}
	}
}

func TestFeedID(t *testing.T) {
	var tests = []struct {
		//given
		description string
		feedID      string

		//expected
		validatedFeedID string
		errorMsg        string
	}{
		{
			"Feed ID with word chars",
			"123",

			"123",
			"",
		},
		{
			"Feed ID word chars and dashes",
			"1-2-3",

			"1-2-3",
			"",
		},
		{
			"Feed ID with invalid chars",
			"123@",

			"",
			"invalid feedID",
		},
	}

	for _, test := range tests {
		feedID, err := getstream.ValidateFeedID(test.feedID)
		if feedID != test.validatedFeedID {
			t.Errorf("Expected feedID %v, got %v", test.validatedFeedID, feedID)
		}

		if err != nil && err.Error() != test.errorMsg {
			t.Errorf("Error %v does not match expected %v", test.errorMsg, err.Error())
		}
	}
}

func TestUserID(t *testing.T) {
	var tests = []struct {
		//given
		description string
		userID      string

		//expected
		validatedUserID string
		errorMsg        string
	}{
		{
			"User ID with word chars",
			"123",

			"123",
			"",
		},
		{
			"User ID word chars and dashes",
			"1-2-3",

			"1-2-3",
			"",
		},
		{
			"User ID with invalid chars",
			"123@",

			"",
			"invalid userID",
		},
	}

	for _, test := range tests {
		userID, err := getstream.ValidateUserID(test.userID)
		if userID != test.validatedUserID {
			t.Errorf("Expected userID %v, got %v", test.validatedUserID, userID)
		}

		if err != nil && err.Error() != test.errorMsg {
			t.Errorf("Error %v does not match expected %v", test.errorMsg, err.Error())
		}
	}
}
