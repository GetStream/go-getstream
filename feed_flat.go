package getstream

import (
	"encoding/json"
)

type FlatFeed struct {
	baseFeed
}

// GetFlatFeedInput is used to Get a list of Activities from a FlatFeed
type GetFlatFeedInput struct {
	typeFeedReadOptions

	Ranking string
}

// GetFlatFeedOutput is the response from a FlatFeed Activities Get Request
type GetFlatFeedOutput struct {
	Duration   string      `json:"duration"`
	Next       string      `json:"next"`
	Activities []*Activity `json:"results"`
}

// Activities returns a list of Activities for a FlatFeedGroup
func (f *FlatFeed) Activities(input *GetFlatFeedInput) (*GetFlatFeedOutput, error) {
	var err error

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/"

	result, err := f.Client.get(f, endpoint, nil, input.Params())

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
