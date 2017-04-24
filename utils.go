package getstream

import (
	"errors"
	"regexp"
)

var validWord = regexp.MustCompile(`^[\w|-]+$`)

func ValidateFeedSlug(feedSlug string) (string, error) {
	if !validWord.MatchString(feedSlug) {
		return "", errors.New("invalid feedSlug")
	}

	return feedSlug, nil
}

func ValidateFeedID(feedID string) (string, error) {
	if !validWord.MatchString(feedID) {
		return "", errors.New("invalid feedID")
	}

	return feedID, nil
}

func ValidateUserID(userID string) (string, error) {
	if !validWord.MatchString(userID) {
		return "", errors.New("invalid userID")
	}

	return userID, nil
}
