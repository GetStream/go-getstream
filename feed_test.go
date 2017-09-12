package getstream_test

import (
	"testing"
	"time"

	getstream "github.com/GetStream/stream-go"
	"github.com/stretchr/testify/require"
)

func TestUpdateActivityToTargets(t *testing.T) {
	client := PreTestSetup(t)
	f1 := getFlatFeed(t, client)
	f2 := getFlatFeed(t, client)
	f3 := getFlatFeed(t, client)

	now := time.Now()

	activity := &getstream.Activity{
		Actor:     "bob",
		Verb:      "like",
		Object:    "cakes",
		ForeignID: "bob:123",
		TimeStamp: &now,
		To:        []getstream.FeedID{f2.FeedID()},
	}
	_, err := f1.AddActivity(activity)
	require.NoError(t, err)

	testCases := []struct {
		activity            *getstream.Activity
		shouldError         bool
		news, adds, removes []getstream.FeedID
	}{
		{
			shouldError: true,
		},
		{
			activity:    &getstream.Activity{},
			shouldError: true,
		},
		{
			activity:    activity,
			news:        []getstream.FeedID{f3.FeedID()},
			shouldError: false,
		},
	}

	for _, tc := range testCases {
		err := f1.UpdateActivityToTargets(tc.activity, tc.news, tc.adds, tc.removes)
		if tc.shouldError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}
