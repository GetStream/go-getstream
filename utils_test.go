package getstream_test

import (
	"testing"

	getstream "github.com/GetStream/stream-go"
)

func TestFeedSlug(t *testing.T) {
	feedSlug, err := getstream.ValidateFeedSlug("foo")
	if err != nil {
		t.Error(err)
	}
	if feedSlug != "foo" {
		t.Error("feedSlug not 'foo'")
	}

	feedSlug, err = getstream.ValidateFeedSlug("f-o-o")
	if err != nil {
		t.Error(err)
	}
	if feedSlug != "f_o_o" {
		t.Error("feedSlug not 'f_o_o'")
	}
}

func TestFeedID(t *testing.T) {
	feedID, err := getstream.ValidateFeedID("123")
	if err != nil {
		t.Error(err)
	}
	if feedID != "123" {
		t.Error("feedID not '123'")
	}

	feedID, err = getstream.ValidateFeedID("1-2-3")
	if err != nil {
		t.Error(err)
	}
	if feedID != "1_2_3" {
		t.Error("feedID not '1_2_3'")
	}
}

func TestUserID(t *testing.T) {
	userID, err := getstream.ValidateUserID("123")
	if err != nil {
		t.Error(err)
	}
	if userID != "123" {
		t.Error("userID not '123'")
	}

	userID, err = getstream.ValidateUserID("1-2-3")
	if err != nil {
		t.Error(err)
	}
	if userID != "1_2_3" {
		t.Error("userSlug not '1_2_3'")
	}
}
