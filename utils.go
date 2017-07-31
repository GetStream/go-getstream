package getstream

import (
	"errors"
	"regexp"
)

func ValidateFeedSlug(feedSlug string) (string, error) {
	validFeedSlug := regexp.MustCompile(`^\w+$`)

	if !validFeedSlug.MatchString(feedSlug) {
		return "", errors.New("invalid feedSlug")
	}

	return feedSlug, nil
}

func ValidateFeedID(feedID string) (string, error) {
	validFeedID := regexp.MustCompile(`^\w+$`)

	if !validFeedID.MatchString(feedID) {
		return "", errors.New("invalid feedID")
	}

	return feedID, nil
}

func ValidateUserID(userID string) (string, error) {
	validUserID := regexp.MustCompile(`^[\w-]+$`)

	if !validUserID.MatchString(userID) {
		return "", errors.New("invalid userID")
	}

	return userID, nil
}
