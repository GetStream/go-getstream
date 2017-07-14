package getstream

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
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

	RemoveActivity(input *Activity) error
	RemoveActivityByForeignID(input *Activity) error

	Unfollow(target *FlatFeed) error
	UnfollowKeepingHistory(target *FlatFeed) error

	FollowFeedWithCopyLimit(target Feed, copyLimit int) error

	//TODO: change this to return a list of []Feed
	FollowersWithLimitAndSkip(limit int, skip int) ([]*GeneralFeed, error)
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

// AddActivity is used to add an Activity
func (f *baseFeed) AddActivity(activity *Activity) (*Activity, error) {

	payload, err := json.Marshal(activity)
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

	payload, err := json.Marshal(map[string][]*Activity{
		"activities": activities,
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

// RemoveActivity removes an Activity
func (f *baseFeed) RemoveActivity(input *Activity) error {

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + input.ID + "/"

	return f.Client.del(f, endpoint, nil, nil)
}

// RemoveActivityByForeignID removes an Activity by ForeignID
func (f *baseFeed) RemoveActivityByForeignID(input *Activity) error {

	if input.ForeignID == "" {
		return errors.New("no ForeignID")
	}

	r, err := regexp.Compile("^[a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12}$")
	if err != nil {
		return err
	}
	if !r.MatchString(input.ForeignID) {
		return errors.New("invalid ForeignID")
	}

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + input.ForeignID + "/"

	return f.Client.del(f, endpoint, nil, map[string]string{
		"foreign_id": "1",
	})
}

// Unfollow is used to Unfollow a target Feed
func (f *baseFeed) Unfollow(target *FlatFeed) error {
	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "following" + "/" + target.FeedID().Value() + "/"
	return f.Client.del(f, endpoint, nil, nil)
}

// UnfollowKeepingHistory is used to Unfollow a target Feed while keeping the History
// this means that Activities already visibile will remain
func (f *baseFeed) UnfollowKeepingHistory(target *FlatFeed) error {
	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "following" + "/" + target.FeedID().Value() + "/"

	payload, err := json.Marshal(map[string]string{
		"keep_history": "1",
	})
	if err != nil {
		return err
	}

	return f.Client.del(f, endpoint, payload, nil)
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

// FollowersWithLimitAndSkip returns a list of GeneralFeed that are following the feed
// TODO: rename skip into offset
func (f *baseFeed) FollowersWithLimitAndSkip(limit int, skip int) ([]*GeneralFeed, error) {
	var err error

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "followers" + "/"
	params := map[string]string{
		"limit":  strconv.Itoa(limit),
		"offset": strconv.Itoa(skip),
	}

	if err != nil {
		return nil, err
	}

	resultBytes, err := f.Client.get(f, endpoint, nil, params)

	output := &feedFollows{}
	err = json.Unmarshal(resultBytes, output)
	if err != nil {
		return nil, err
	}

	var outputFeeds []*GeneralFeed
	for _, result := range output.Results {

		feed := GeneralFeed{}

		var match bool
		match, err = regexp.MatchString(`^.*?:.*?$`, result.FeedID)
		if err != nil {
			continue
		}

		if match {
			firstSplit := strings.Split(result.FeedID, ":")

			feed.FeedSlug = firstSplit[0]
			feed.UserID = firstSplit[1]
		}

		outputFeeds = append(outputFeeds, &feed)
	}

	return outputFeeds, err
}

// FollowingWithLimitAndSkip returns a list of GeneralFeed followed by the current FlatFeed
// TODO: need to support filters
func (f *baseFeed) FollowingWithLimitAndSkip(limit int, skip int) ([]*GeneralFeed, error) {
	var err error

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "following" + "/"

	params := map[string]string{
		"limit":  strconv.Itoa(limit),
		"offset": strconv.Itoa(skip),
	}

	resultBytes, err := f.Client.get(f, endpoint, nil, params)

	output := feedFollows{}
	err = json.Unmarshal(resultBytes, &output)
	if err != nil {
		return nil, err
	}

	var outputFeeds []*GeneralFeed
	for _, result := range output.Results {

		feed := GeneralFeed{}

		var match bool
		match, err = regexp.MatchString(`^.*?:.*?$`, result.FeedID)
		if err != nil {
			continue
		}

		if match {
			firstSplit := strings.Split(result.TargetID, ":")

			feed.FeedSlug = firstSplit[0]
			feed.UserID = firstSplit[1]
		}

		outputFeeds = append(outputFeeds, &feed)
	}

	return outputFeeds, err
}

// FollowFeedWithCopyLimit sets a Feed to follow another target Feed
// CopyLimit is the maximum number of Activities to Copy from History
func (f *baseFeed) FollowFeedWithCopyLimit(target Feed, copyLimit int) error {
	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "following" + "/"

	input := struct {
		Target            string `json:"target"`
		ActivityCopyLimit int    `json:"activity_copy_limit"`
	}{
		Target:            target.FeedID().Value(),
		ActivityCopyLimit: copyLimit,
	}

	payload, err := json.Marshal(input)
	if err != nil {
		return err
	}

	_, err = f.Client.post(f, endpoint, payload, nil)
	return err
}

type typeFeedReadOptions struct {
	Limit  int
	Offset int

	IDGTE string
	IDGT  string
	IDLTE string
	IDLT  string
}

func (i *typeFeedReadOptions) Params() (params map[string]string) {
	params = make(map[string]string)

	if i.Limit != 0 {
		params["limit"] = fmt.Sprintf("%d", i.Limit)
	}
	if i.Offset != 0 {
		params["offset"] = fmt.Sprintf("%d", i.Offset)
	}
	if i.IDGTE != "" {
		params["id_gte"] = i.IDGTE
	}
	if i.IDGT != "" {
		params["id_gt"] = i.IDGT
	}
	if i.IDLTE != "" {
		params["id_lte"] = i.IDLTE
	}
	if i.IDLT != "" {
		params["id_lt"] = i.IDLT
	}
	return params
}
