package getstream

import (
	"errors"
	"regexp"
	"strings"
)

func ValidateFeedSlug(feedSlug string) (string, error) {
	r, err := regexp.Compile(`^\w+$`)
	if err != nil {
		return "", err
	}

	feedSlug = strings.Replace(feedSlug, "-", "_", -1)

	if !r.MatchString(feedSlug) {
		return "", errors.New("invalid feedSlug")
	}

	return feedSlug, nil
}

func ValidateFeedID(feedID string) (string, error) {
	r, err := regexp.Compile(`^\w+$`)
	if err != nil {
		return "", err
	}

	feedID = strings.Replace(feedID, "-", "_", -1)

	if !r.MatchString(feedID) {
		return "", errors.New("invalid feedID")
	}

	return feedID, nil
}

func ValidateUserID(userID string) (string, error) {
	r, err := regexp.Compile(`^\w+$`)
	if err != nil {
		return "", err
	}

	userID = strings.Replace(userID, "-", "_", -1)

	if !r.MatchString(userID) {
		return "", errors.New("invalid userID")
	}

	return userID, nil
}
