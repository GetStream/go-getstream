package getstream

import (
	"encoding/json"
	"strconv"
	"strings"
)

type NotificationFeed struct {
	baseFeed
}

// GetNotificationFeedOutput is the response from a NotificationFeed Activities Get Request
type GetNotificationFeedOutput struct {
	Duration string
	Next     string
	Results  []*struct {
		Activities    []*Activity
		ActivityCount int
		ActorCount    int
		CreatedAt     string
		Group         string
		ID            string
		IsRead        bool
		IsSeen        bool
		UpdatedAt     string
		Verb          string
	}
	Unread int
	Unseen int
}

type getNotificationFeedOutput struct {
	Duration string                             `json:"duration"`
	Next     string                             `json:"next"`
	Results  []*getNotificationFeedOutputResult `json:"results"`
	Unread   int                                `json:"unread"`
	Unseen   int                                `json:"unseen"`
}

func (a getNotificationFeedOutput) output() *GetNotificationFeedOutput {

	output := GetNotificationFeedOutput{
		Duration: a.Duration,
		Next:     a.Next,
		Unread:   a.Unread,
		Unseen:   a.Unseen,
	}

	var results []*struct {
		Activities    []*Activity
		ActivityCount int
		ActorCount    int
		CreatedAt     string
		Group         string
		ID            string
		IsRead        bool
		IsSeen        bool
		UpdatedAt     string
		Verb          string
	}

	for _, result := range a.Results {

		outputResult := struct {
			Activities    []*Activity
			ActivityCount int
			ActorCount    int
			CreatedAt     string
			Group         string
			ID            string
			IsRead        bool
			IsSeen        bool
			UpdatedAt     string
			Verb          string
		}{
			ActivityCount: result.ActivityCount,
			ActorCount:    result.ActorCount,
			CreatedAt:     result.CreatedAt,
			Group:         result.Group,
			ID:            result.ID,
			IsRead:        result.IsRead,
			IsSeen:        result.IsSeen,
			UpdatedAt:     result.UpdatedAt,
			Verb:          result.Verb,
		}

		for _, activity := range result.Activities {
			outputResult.Activities = append(outputResult.Activities, activity)
		}

		results = append(results, &outputResult)
	}
	output.Results = results
	return &output
}

type getNotificationFeedOutputResult struct {
	Activities    []*Activity `json:"activities"`
	ActivityCount int         `json:"activity_count"`
	ActorCount    int         `json:"actor_count"`
	CreatedAt     string      `json:"created_at"`
	Group         string      `json:"group"`
	ID            string      `json:"id"`
	IsRead        bool        `json:"is_read"`
	IsSeen        bool        `json:"is_seen"`
	UpdatedAt     string      `json:"updated_at"`
	Verb          string      `json:"verb"`
}

// MarkActivitiesAsRead marks activities as read for this feed
func (f *NotificationFeed) MarkActivitiesAsRead(activities []*Activity) error {

	var ids []string
	for _, activity := range activities {
		ids = append(ids, activity.ID)
	}

	idStr := strings.Join(ids, ",")

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/"

	_, err := f.Client.get(f, endpoint, nil, map[string]string{
		"mark_read": idStr,
	})

	return err
}

// MarkActivitiesAsSeenWithLimit marks activities as seen for this feed
func (f *NotificationFeed) MarkActivitiesAsSeenWithLimit(limit int) error {

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/"

	_, err := f.Client.get(f, endpoint, nil, map[string]string{
		"mark_seen": "true",
		"limit":     strconv.Itoa(limit),
	})

	return err
}

// Activities returns a list of Activities for a NotificationFeedGroup
func (f *NotificationFeed) Activities(input FeedReadOptions) (*GetNotificationFeedOutput, error) {
	var (
		payload []byte
		err     error
	)

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/"
	result, err := f.Client.get(f, endpoint, payload, input.Params())

	if err != nil {
		return nil, err
	}

	response := getNotificationFeedOutput{}
	err = json.Unmarshal(result, &response)
	if err != nil {
		return nil, err
	}

	return response.output(), err
}
