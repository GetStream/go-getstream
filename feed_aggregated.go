package getstream

import (
	"encoding/json"
)

type AggregatedFeed struct {
	baseFeed
}

// GetAggregatedFeedInput is used to Get a list of Activities from a AggregatedFeed
type GetAggregatedFeedInput struct {
	typeFeedReadOptions
}

// GetAggregatedFeedOutput is the response from a AggregatedFeed Activities Get Request
type GetAggregatedFeedOutput struct {
	Duration string
	Next     string
	Results  []*struct {
		Activities    []*Activity
		ActivityCount int
		ActorCount    int
		CreatedAt     string
		Group         string
		ID            string
		UpdatedAt     string
		Verb          string
	}
}

type getAggregatedFeedOutput struct {
	Duration string                           `json:"duration"`
	Next     string                           `json:"next"`
	Results  []*getAggregatedFeedOutputResult `json:"results"`
}

func (a getAggregatedFeedOutput) output() *GetAggregatedFeedOutput {

	output := GetAggregatedFeedOutput{
		Duration: a.Duration,
		Next:     a.Next,
	}

	var results []*struct {
		Activities    []*Activity
		ActivityCount int
		ActorCount    int
		CreatedAt     string
		Group         string
		ID            string
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
			UpdatedAt     string
			Verb          string
		}{
			ActivityCount: result.ActivityCount,
			ActorCount:    result.ActorCount,
			CreatedAt:     result.CreatedAt,
			Group:         result.Group,
			ID:            result.ID,
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

type getAggregatedFeedOutputResult struct {
	Activities    []*Activity `json:"activities"`
	ActivityCount int         `json:"activity_count"`
	ActorCount    int         `json:"actor_count"`
	CreatedAt     string      `json:"created_at"`
	Group         string      `json:"group"`
	ID            string      `json:"id"`
	UpdatedAt     string      `json:"updated_at"`
	Verb          string      `json:"verb"`
}

// Activities returns a list of Activities for a NotificationFeedGroup
func (f *AggregatedFeed) Activities(input GetAggregatedFeedInput) (*GetAggregatedFeedOutput, error) {
	var (
		payload []byte
		err     error
	)

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/"
	result, err := f.Client.get(f, endpoint, payload, input.Params())

	if err != nil {
		return nil, err
	}

	output := &getAggregatedFeedOutput{}
	err = json.Unmarshal(result, output)
	if err != nil {
		return nil, err
	}

	return output.output(), err
}
