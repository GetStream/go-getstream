package getstream

import (
	"encoding/json"
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type postFlatFeedOutputActivities struct {
	Activities []*Activity `json:"activities"`
}

// GetFlatFeedInput is used to Get a list of Activities from a FlatFeed
type GetFlatFeedInput struct {
	Limit  int `json:"limit,omitempty"`
	Offset int `json:"offset,omitempty"`

	IDGTE string `json:"id_gte,omitempty"`
	IDGT  string `json:"id_gt,omitempty"`
	IDLTE string `json:"id_lte,omitempty"`
	IDLT  string `json:"id_lt,omitempty"`

	Ranking string `json:"ranking,omitempty"`
}

// GetFlatFeedOutput is the response from a FlatFeed Activities Get Request
type GetFlatFeedOutput struct {
	Duration   string      `json:"duration"`
	Next       string      `json:"next"`
	Activities []*Activity `json:"results"`
}

type getFlatFeedFollowersInput struct {
	Limit int `json:"limit"`
	Skip  int `json:"offset"`
}

type getFlatFeedFollowersOutput struct {
	Duration string                              `json:"duration"`
	Results  []*getFlatFeedFollowersOutputResult `json:"results"`
}

type getFlatFeedFollowersOutputResult struct {
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	FeedID    string `json:"feed_id"`
	TargetID  string `json:"target_id"`
}

type postFlatFeedFollowingInput struct {
	Target            string `json:"target"`
	ActivityCopyLimit int    `json:"activity_copy_limit"`
}

// FlatFeed is a getstream FlatFeed
// Use it to for CRUD on FlatFeed Groups
type FlatFeed struct {
	Client   *Client
	FeedSlug string
	UserID   string
	token    string
}

// Signature is used to sign Requests : "FeedSlugUserID Token"
func (f *FlatFeed) Signature() string {
	if f.Token() == "" {
		return f.FeedIDWithoutColon()
	}
	return f.FeedIDWithoutColon() + " " + f.Token()
}

// FeedID is the combo if the FeedSlug and UserID : "FeedSlug:UserID"
func (f *FlatFeed) FeedID() FeedID {
	return FeedID(f.FeedSlug + ":" + f.UserID)
}

func (f *FlatFeed) FeedIDWithoutColon() string {
	return f.FeedSlug + f.UserID
}

// SignFeed sets the token on a Feed
func (f *FlatFeed) SignFeed(signer *Signer) {
	if f.Client.Signer != nil {
		f.token = signer.GenerateToken(f.FeedIDWithoutColon())
	}
}

// Token returns the token of a Feed
func (f *FlatFeed) Token() string {
	return f.token
}

// GenerateToken returns a new Token for a Feed without setting it to the Feed
func (f *FlatFeed) GenerateToken(signer *Signer) string {
	if f.Client.Signer != nil {
		return signer.GenerateToken(f.FeedSlug + f.UserID)
	}
	return ""
}

// AddActivity is used to add an Activity to a FlatFeed
func (f *FlatFeed) AddActivity(activity *Activity) (*Activity, error) {

	activity.ID = ""

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

// AddActivities is used to add multiple Activities to a FlatFeed
func (f *FlatFeed) AddActivities(activities []*Activity) ([]*Activity, error) {
	for _, activity := range activities {
		activity.ID = ""
	}

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

	output := &postFlatFeedOutputActivities{}
	err = json.Unmarshal(resultBytes, output)
	if err != nil {
		return nil, err
	}

	return output.Activities, err
}

// Activities returns a list of Activities for a FlatFeedGroup
func (f *FlatFeed) Activities(input *GetFlatFeedInput) (*GetFlatFeedOutput, error) {

	var payload []byte
	var err error

	if input != nil {
		payload, err = json.Marshal(input)
		if err != nil {
			return nil, err
		}
	}

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/"

	result, err := f.Client.get(f, endpoint, payload, nil)
	if err != nil {
		return nil, err
	}

	output := &GetFlatFeedOutput{}
	err = json.Unmarshal(result, output)
	if err != nil {
		return nil, err
	}

	return output, err
}

// RemoveActivity removes an Activity from a FlatFeedGroup
func (f *FlatFeed) RemoveActivity(input *Activity) error {

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + input.ID + "/"

	return f.Client.del(f, endpoint, nil, nil)
}

// RemoveActivityByForeignID removes an Activity from a FlatFeedGroup by ForeignID
func (f *FlatFeed) RemoveActivityByForeignID(input *Activity) error {

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

// FollowFeedWithCopyLimit sets a Feed to follow another target Feed
// CopyLimit is the maximum number of Activities to Copy from History
func (f *FlatFeed) FollowFeedWithCopyLimit(target *FlatFeed, copyLimit int) error {

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "following" + "/"

	input := postFlatFeedFollowingInput{
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

// Unfollow is used to Unfollow a target Feed
func (f *FlatFeed) Unfollow(target *FlatFeed) error {
	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "following" + "/" + target.FeedID().Value() + "/"

	return f.Client.del(f, endpoint, nil, nil)
}

// Unfollow is used to Unfollow a target Aggregated Feed
func (f *FlatFeed) UnfollowAggregated(target *AggregatedFeed) error {
	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "following" + "/" + target.FeedID().Value() + "/"

	return f.Client.del(f, endpoint, nil, nil)
}

// Unfollow is used to Unfollow a target Notification Feed
func (f *FlatFeed) UnfollowNotification(target *NotificationFeed) error {
	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "following" + "/" + target.FeedID().Value() + "/"

	return f.Client.del(f, endpoint, nil, nil)
}

// UnfollowKeepingHistory is used to Unfollow a target Feed while keeping the History
// this means that Activities already visibile will remain
func (f *FlatFeed) UnfollowKeepingHistory(target *FlatFeed) error {

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "following" + "/" + target.FeedID().Value() + "/"

	payload, err := json.Marshal(map[string]string{
		"keep_history": "1",
	})
	if err != nil {
		return err
	}

	return f.Client.del(f, endpoint, payload, nil)
}

// FollowersWithLimitAndSkip returns a list of GeneralFeed following the current FlatFeed
func (f *FlatFeed) FollowersWithLimitAndSkip(limit int, skip int) ([]*GeneralFeed, error) {
	var err error

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "followers" + "/"

	payload, err := json.Marshal(&getFlatFeedFollowersInput{
		Limit: limit,
		Skip:  skip,
	})
	if err != nil {
		return nil, err
	}

	resultBytes, err := f.Client.get(f, endpoint, payload, nil)

	output := &getFlatFeedFollowersOutput{}
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
func (f *FlatFeed) FollowingWithLimitAndSkip(limit int, skip int) ([]*GeneralFeed, error) {
	var err error

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "following" + "/"

	payload, err := json.Marshal(&getFlatFeedFollowersInput{
		Limit: limit,
		Skip:  skip,
	})
	if err != nil {
		return nil, err
	}

	resultBytes, err := f.Client.get(f, endpoint, payload, nil)

	output := &getFlatFeedFollowersOutput{}
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
			firstSplit := strings.Split(result.TargetID, ":")

			feed.FeedSlug = firstSplit[0]
			feed.UserID = firstSplit[1]
		}

		outputFeeds = append(outputFeeds, &feed)
	}

	return outputFeeds, err
}

/** FollowFeedsWithCopyLimit sets a Feed to follow one or more other target Feeds
	This method only exists within FlatFeed because only flat feeds can follow other feeds

	Params:
	sourceFeeds, a list of feeds this feed can follow
	copyLimit, optional number of items to copy from history, defaults to 100

 	Returns:
 	error, if any
*/
func (f *FlatFeed) FollowManyFeeds(sourceFeeds []PostFlatFeedFollowingManyInput, copyLimit int) error {

	final_payload, err := json.Marshal(sourceFeeds)
	if err != nil {
		return err
	}

	var params = map[string]string{}
	if copyLimit < 0 {
		copyLimit = 100
	}
	params = map[string]string{
		"activity_copy_limit": strconv.Itoa(copyLimit),
	}

	endpoint := "follow_many/"

	//save_token := ""
	//if f.token != "" {
		//fmt.Println("saving token for later")
		//save_token = f.token
		//f.token = ""
	//}
	_, err = f.Client.post(f, endpoint, final_payload, params)
	//if save_token != "" {
		//fmt.Println("restoring token")
		//f.token = save_token
	//}
	return err
}

type postMultipleActivities struct {
	Activities []*Activity `json:"activities"`
}

func (f *FlatFeed) UpdateActivities(activities []*Activity) error {
	if len(activities) == 0 {
		return errors.New("No activities to update")
	}

	// verify/exclude by foreign id
	var verifiedActivities []*Activity
	for _, activity := range activities {
		if activity.ForeignID != "" {
			verifiedActivities = append(verifiedActivities, activity)
		}
	}

	// verify that there are no more than 100 to update
	if len(verifiedActivities) > 100 {
		return errors.New("Cannot update more than 100 activities at a time")
	}
	if len(verifiedActivities) == 0 {
		return errors.New("No activities to update (no ForeignID values)")
	}

	final_payload, err := json.Marshal(&postMultipleActivities{
		Activities: verifiedActivities,
	})
	if err != nil {
		return err
	}

	endpoint := "activities/"
	params := map[string]string{}

	_, err = f.Client.post(f, endpoint, final_payload, params)
	if err != nil {
		return err
	}

	return nil
}

func (f *FlatFeed) UpdateActivity(activity *Activity) error {
	return f.UpdateActivities([]*Activity{activity})
}
