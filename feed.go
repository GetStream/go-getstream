package getstream

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type FeedID string

// Value returns a String Representation of FeedID
func (f FeedID) Value() string {
	return string(f)
}

// Feed is the interface bundling all Feed Types
// It exposes methods needed for all Types
type Feed interface {
	Signature() string
	FeedID() FeedID
	FeedIDWithoutColon() string
	Token() string
	SignFeed(signer *Signer)
	GenerateToken(signer *Signer) string

	AddActivity(activity *Activity) (*Activity, error)
	AddActivities(activities []*Activity) ([]*Activity, error)

	RemoveActivity(activityId string) error
	RemoveActivityByForeignID(foreignId string) error

	UpdateActivityToTargets(activity *Activity, news, adds, removes []string) error

	Follow(target Feed) error
	FollowFeedWithCopyLimit(target Feed, copyLimit int) error
	Unfollow(target Feed) error
	UnfollowKeepingHistory(target Feed) error

	GetFollowers(limit int, offset int) ([]FeedID, error)
	GetFollowings(limit int, offset int) ([]FeedID, error)
	GetFollowingsFiltered(limit int, offset int, filter []string) ([]FeedID, error)
}

// A collection of common code between all feeds
type baseFeed struct {
	FeedSlug string
	UserID   string
	Client   *Client
	token    string
}

func (f *baseFeed) FeedIDWithoutColon() string {
	return f.FeedSlug + f.UserID
}

// Signature is used to sign Requests : "FeedSlugUserID Token"
func (f *baseFeed) Signature() string {
	if f.Token() == "" {
		return f.FeedIDWithoutColon()
	}
	return f.FeedIDWithoutColon() + " " + f.Token()
}

// FeedID is the combo if the FeedSlug and UserID : "FeedSlug:UserID"
func (f *baseFeed) FeedID() FeedID {
	return FeedID(f.FeedSlug + ":" + f.UserID)
}

// SignFeed sets the token on a Feed
func (f *baseFeed) SignFeed(signer *Signer) {
	if f.Client.Signer != nil {
		f.token = signer.GenerateToken(f.FeedIDWithoutColon())
	}
}

// Token returns the token of a Feed
func (f *baseFeed) Token() string {
	return f.token
}

// GenerateToken returns a new Token for a Feed without setting it to the Feed
func (f *baseFeed) GenerateToken(signer *Signer) string {
	if f.Client.Signer != nil {
		return signer.GenerateToken(f.FeedSlug + f.UserID)
	}
	return ""
}

func (f *baseFeed) signToField(activity *Activity) {
	for i := range activity.To {
		bits := strings.Split(string(activity.To[i]), ":")
		if len(bits) != 2 {
			continue
		}
		signature := f.Client.Signer.SignFeed(strings.Join(bits, ""))
		signedTo := fmt.Sprintf("%s %s", activity.To[i], signature)
		activity.signedTo = append(activity.signedTo, signedTo)
	}
}

// AddActivity is used to add an Activity
func (f *baseFeed) AddActivity(activity *Activity) (*Activity, error) {
	activityCopy := &Activity{}

	*activityCopy = *activity
	f.signToField(activityCopy)
	payload, err := json.Marshal(activityCopy)

	if err != nil {
		return nil, err
	}

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/"

	resultBytes, err := f.Client.post(f, endpoint, payload, nil)
	if err != nil {
		return nil, err
	}

	output := &Activity{}
	err = json.Unmarshal(resultBytes, output)
	if err != nil {
		return nil, err
	}

	return output, err
}

// AddActivities is used to add multiple Activities
func (f *baseFeed) AddActivities(activities []*Activity) ([]*Activity, error) {
	var (
		activitiesCopy []*Activity
	)

	for i := range activities {
		activityCopy := &Activity{}
		*activityCopy = *activities[i]
		f.signToField(activityCopy)
		activitiesCopy = append(activitiesCopy, activityCopy)
	}

	payload, err := json.Marshal(map[string][]*Activity{
		"activities": activitiesCopy,
	})
	if err != nil {
		return nil, err
	}

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/"

	resultBytes, err := f.Client.post(f, endpoint, payload, nil)
	if err != nil {
		return nil, err
	}

	output := &struct {
		Activities []*Activity `json:"activities"`
	}{}

	err = json.Unmarshal(resultBytes, output)
	if err != nil {
		return nil, err
	}

	return output.Activities, err
}

// RemoveActivity removes an Activity by its ID
func (f *baseFeed) RemoveActivity(activityId string) error {
	if activityId == "" {
		return errors.New("activityId must be a non empty string")
	}
	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + activityId + "/"
	return f.Client.del(f, endpoint, nil, nil)
}

// RemoveActivityByForeignID performs a delete by ForeignID
func (f *baseFeed) RemoveActivityByForeignID(foreignId string) error {
	if foreignId == "" {
		return errors.New("foreignId must be a non empty string")
	}
	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + foreignId + "/"
	return f.Client.del(f, endpoint, nil, map[string]string{
		"foreign_id": "1",
	})
}

func (f *baseFeed) UpdateActivities(activities ...*Activity) error {
	payload, err := json.Marshal(UpdateActivitiesRequest{Activities: activities})
	if err != nil {
		return fmt.Errorf("cannot marshal payload: %s", err)
	}

	endpoint := "activities/"
	_, err = f.Client.post(nil, endpoint, payload, nil)
	if err != nil {
		return fmt.Errorf("cannot update activities: %s", err)
	}
	return nil
}

// Unfollow is used to Unfollow a target Feed
func (f *baseFeed) Unfollow(target Feed) error {
	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "following" + "/" + target.FeedID().Value() + "/"
	return f.Client.del(f, endpoint, nil, nil)
}

// UnfollowKeepingHistory is used to Unfollow a target Feed while keeping the History
// this means that Activities already visibile will remain
func (f *baseFeed) UnfollowKeepingHistory(target Feed) error {
	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "following" + "/" + target.FeedID().Value() + "/"
	params := map[string]string{
		"keep_history": "1",
	}
	return f.Client.del(f, endpoint, nil, params)
}

type feedFollows struct {
	Duration string `json:"duration"`
	Results  []*struct {
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
		FeedID    string `json:"feed_id"`
		TargetID  string `json:"target_id"`
	} `json:"results"`
}

// GetFollowers returns a list of GeneralFeed that are following the feed
func (f *baseFeed) GetFollowers(limit int, offset int) ([]FeedID, error) {
	var (
		err         error
		outputFeeds []FeedID
	)

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "followers" + "/"
	params := map[string]string{
		"limit":  strconv.Itoa(limit),
		"offset": strconv.Itoa(offset),
	}

	if err != nil {
		return nil, err
	}

	resultBytes, err := f.Client.get(f, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	output := &feedFollows{}
	err = json.Unmarshal(resultBytes, output)
	if err != nil {
		return nil, err
	}

	for _, result := range output.Results {
		outputFeeds = append(outputFeeds, FeedID(result.FeedID))
	}

	return outputFeeds, err
}

// GetFollowings returns a list of FeedID followed by the current FlatFeed
func (f *baseFeed) GetFollowings(limit int, offset int) ([]FeedID, error) {
	return f.GetFollowingsFiltered(limit, offset, nil)
}

// GetFollowingsFiltered returns a list of FeedID followed by the current FlatFeed, filtered with the provided feed IDs.
func (f *baseFeed) GetFollowingsFiltered(limit int, offset int, filter []string) ([]FeedID, error) {
	var (
		err         error
		outputFeeds []FeedID
	)

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "follows" + "/"

	params := map[string]string{
		"limit":  strconv.Itoa(limit),
		"offset": strconv.Itoa(offset),
	}
	if filter != nil {
		params["filter"] = strings.Join(filter, ",")
	}

	resultBytes, err := f.Client.get(f, endpoint, nil, params)
	if err != nil {
		return nil, err
	}

	output := feedFollows{}
	err = json.Unmarshal(resultBytes, &output)
	if err != nil {
		return nil, err
	}

	for _, result := range output.Results {
		outputFeeds = append(outputFeeds, FeedID(result.TargetID))
	}

	return outputFeeds, err
}

func (f *baseFeed) follow(target Feed, copyLimit *int) error {
	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "follows" + "/"

	input := struct {
		Target            string `json:"target"`
		ActivityCopyLimit *int   `json:"activity_copy_limit,omitempty"`
	}{
		Target: target.FeedID().Value(),
	}

	if copyLimit != nil {
		input.ActivityCopyLimit = copyLimit
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return err
	}

	_, err = f.Client.post(f, endpoint, payload, nil)
	return err
}

// Follow sets a Feed to follow another target Feed
func (f *baseFeed) Follow(target Feed) error {
	return f.follow(target, nil)
}

// FollowFeedWithCopyLimit sets a Feed to follow another target Feed
// CopyLimit is the maximum number of Activities to Copy from History
func (f *baseFeed) FollowFeedWithCopyLimit(target Feed, copyLimit int) error {
	return f.follow(target, &copyLimit)
}

// UpdateActivityToTargets updates the to targets for the provided Activity using
// the string slices arguments.
func (f *baseFeed) UpdateActivityToTargets(activity *Activity, news, adds, removes []string) error {
	if activity == nil {
		return fmt.Errorf("activity cannot be nil")
	}
	if activity.TimeStamp == nil {
		return fmt.Errorf("activity timestamp cannot be nil")
	}
	data := struct {
		ForeignID string   `json:"foreign_id,omitempty"`
		Time      string   `json:"time,omitempty"`
		New       []string `json:"new_targets,omitempty"`
		Adds      []string `json:"added_targets,omitempty"`
		Removes   []string `json:"removed_targets,omitempty"`
	}{
		ForeignID: activity.ForeignID,
		Time:      activity.TimeStamp.Format(timeLayout),
		New:       news,
		Adds:      adds,
		Removes:   removes,
	}
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("cannot marshal data: %s", err)
	}

	endpoint := "feed_targets/" + f.FeedSlug + "/" + f.UserID + "/" + "activity_to_targets" + "/"
	_, err = f.Client.post(f, endpoint, payload, nil)
	return err
}
