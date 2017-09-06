package getstream_test

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	getstream "github.com/GetStream/stream-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	flatFeedName   = "flat"
	rankedFeedName = "ranked"
)

func getFlatFeedByName(t *testing.T, client *getstream.Client, name string) *getstream.FlatFeed {
	feed, err := client.FlatFeed(name, RandString(8))
	require.NoError(t, err)
	return feed
}

func getFlatFeed(t *testing.T, client *getstream.Client) *getstream.FlatFeed {
	return getFlatFeedByName(t, client, flatFeedName)
}

func getRankedFeed(t *testing.T, client *getstream.Client) *getstream.FlatFeed {
	return getFlatFeedByName(t, client, rankedFeedName)
}

func TestExampleFlatFeedAddActivity(t *testing.T) {
	client := PreTestSetup(t)
	feed := getFlatFeed(t, client)

	_, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: RandString(8),
		Object:    "flat:eric",
		Actor:     "flat:john",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestFlatFeedAddActivity(t *testing.T) {
	client := PreTestSetup(t)
	feed := getFlatFeed(t, client)

	fid := RandString(8)
	activity, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: fid,
		Object:    "flat:eric",
		Actor:     "flat:john",
	})
	if err != nil {
		t.Fatal(err)
	}

	if activity.Verb != "post" && activity.ForeignID != fid {
		t.Fail()
	}

}

func TestFlatFeedAddActivityWithTo(t *testing.T) {
	client := PreTestSetup(t)
	feed := getFlatFeed(t, client)
	feedTo := getFlatFeed(t, client)
	feedToB := getFlatFeed(t, client)

	fid := RandString(8)
	_, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: fid,
		Object:    "flat:eric",
		Actor:     "flat:john",
		To:        []getstream.FeedID{feedTo.FeedID(), feedToB.FeedID()},
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestFlatFeedUUID(t *testing.T) {
	client := PreTestSetup(t)
	feed := getFlatFeed(t, client)

	activity, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: RandString(8),
		Object:    "flat:eric",
		Actor:     "flat:john",
	})

	if err != nil {
		t.Error(err)
	}

	err = feed.RemoveActivityByForeignID(activity.ForeignID)
	if err != nil {
		t.Error(err)
	}
}

func TestFlatFeedRemoveActivity(t *testing.T) {
	client := PreTestSetup(t)
	feed := getFlatFeed(t, client)

	activity, err := feed.AddActivity(&getstream.Activity{
		Verb:   "post",
		Object: "flat:eric",
		Actor:  "flat:john",
	})

	if err != nil {
		t.Error(err)
	}

	err = feed.RemoveActivity(activity.ID)
	if err != nil {
		t.Error(err)
	}
}

func TestFlatFeedRemoveByForeignIDActivity(t *testing.T) {
	client := PreTestSetup(t)
	feed := getFlatFeed(t, client)

	activity, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: RandString(8),
		Object:    "flat:eric",
		Actor:     "flat:john",
	})
	if err != nil {
		t.Error(err)
	}

	if activity.Verb != "post" && activity.ForeignID != "08f01c47-014f-11e4-aa8f-0cc47a024be0" {
		t.Fail()
	}

	err = feed.RemoveActivityByForeignID(activity.ForeignID)
	if err != nil {
		t.Fatal(err)
	}

}

func TestFlatFeedGetActivities(t *testing.T) {
	client := PreTestSetup(t)
	feed := getFlatFeed(t, client)

	_, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: RandString(8),
		Object:    "flat:eric",
		Actor:     "flat:john",
	})
	if err != nil {
		t.Fatal(err)
	}

	activities, err := feed.Activities(getstream.NewFeedReadOptions())
	if err != nil {
		t.Fatal(err)
	}

	if activities.Activities[0].Actor != "flat:john" {
		t.Fatal("Activity read from stream did not match")
	}

}

func TestFlatFeedAddActivities(t *testing.T) {
	client := PreTestSetup(t)
	feed := getFlatFeed(t, client)

	_, err := feed.AddActivities([]*getstream.Activity{
		{
			Verb:      "post",
			ForeignID: RandString(8),
			Object:    "flat:eric",
			Actor:     "flat:john",
		}, {
			Verb:      "walk",
			ForeignID: RandString(8),
			Object:    "flat:john",
			Actor:     "flat:eric",
		},
	})
	if err != nil {
		t.Fatal(err)
	}

}

func TestFlatFeedFollow(t *testing.T) {
	client := PreTestSetup(t)

	feedA := getFlatFeed(t, client)
	feedB := getFlatFeed(t, client)

	if err := feedA.FollowFeedWithCopyLimit(feedB, 20); err != nil {
		t.Fatal(err)
	}

	// get feedB's followers, ensure feedA is there
	followers, err := feedB.GetFollowers(5, 0)
	if err != nil {
		t.Error(err)
	}
	if followers[0] != feedA.FeedID() {
		t.Error("Bob's FeedA is not a follower of FeedB")
	}

	// get things that feedA follows, ensure feedB is in there
	following, err := feedA.GetFollowings(5, 0)
	if err != nil {
		t.Fatal(err)
	}
	if following[0] != feedB.FeedID() {
		t.Error("Eric's FeedB is not a follower of FeedA")
	}
}

func TestFlatFeedFollowingFollowers(t *testing.T) {
	client := PreTestSetup(t)

	feedA := getFlatFeed(t, client)
	feedB := getFlatFeed(t, client)
	feedC := getFlatFeed(t, client)

	err := feedA.FollowFeedWithCopyLimit(feedB, 20)
	if err != nil {
		t.Fatal(err)
	}

	err = feedA.FollowFeedWithCopyLimit(feedC, 20)
	if err != nil {
		t.Fatal(err)
	}

	_, err = feedA.GetFollowings(20, 0)
	if err != nil {
		t.Fatal(err)
	}

	_, err = feedB.GetFollowers(20, 0)
	if err != nil {
		t.Fatal(err)
	}

}

func TestFlatFeedUnFollow(t *testing.T) {
	client := PreTestSetup(t)

	feedA := getFlatFeed(t, client)
	feedB := getFlatFeed(t, client)

	err := feedA.FollowFeedWithCopyLimit(feedB, 20)
	if err != nil {
		t.Fatal(err)
	}

	err = feedA.Unfollow(feedB)
	if err != nil {
		t.Fatal(err)
	}

}

func TestFlatFeedUnFollowKeepingHistory(t *testing.T) {
	client := PreTestSetup(t)
	feedA := getFlatFeed(t, client)
	feedB := getFlatFeed(t, client)

	prepareUnfollowKeepingHistory(t, feedA, feedB, func() {
		_, err := feedA.Activities(getstream.NewFeedReadOptions())
		require.NoError(t, err)
	})

	out, err := feedA.Activities(getstream.NewFeedReadOptions())
	require.NoError(t, err)

	assert.Len(t, out.Activities, 10)
}

func TestFlatActivityMetaData(t *testing.T) {
	now := time.Now()

	data := struct {
		Foo  string
		Fooz string
	}{
		Foo:  "foo",
		Fooz: "fooz",
	}

	dataB, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	raw := json.RawMessage(dataB)

	activity := getstream.Activity{
		ForeignID: RandString(8),
		Actor:     "user:eric",
		Object:    "user:bob",
		Target:    "user:john",
		Verb:      "post",
		TimeStamp: &now,
		Data:      &raw,
		MetaData: map[string]interface{}{
			"meta": "data",
		},
	}

	b, err := json.Marshal(&activity)
	if err != nil {
		t.Fatal(err)
	}

	b2, err := json.Marshal(activity)
	if err != nil {
		t.Fatal(err)
	}

	resultActivity := getstream.Activity{}
	err = json.Unmarshal(b, &resultActivity)
	if err != nil {
		t.Error(err)
	}

	resultActivity2 := getstream.Activity{}
	err = json.Unmarshal(b2, &resultActivity2)
	if err != nil {
		t.Error(err)
	}

	if resultActivity.ForeignID != activity.ForeignID {
		t.Error(activity.ForeignID, resultActivity.ForeignID)
	}
	if resultActivity.Actor != activity.Actor {
		t.Error(activity.Actor, resultActivity.Actor)
	}
	if resultActivity.Verb != activity.Verb {
		t.Error(activity.Verb, resultActivity.Verb)
	}
	if resultActivity.Object != activity.Object {
		t.Error(activity.Object, resultActivity.Object)
	}
	if resultActivity.Target != activity.Target {
		t.Error(activity.Target, resultActivity.Target)
	}
	if resultActivity.TimeStamp.Format("2006-01-02T15:04:05.999999") != activity.TimeStamp.Format("2006-01-02T15:04:05.999999") {
		t.Error(activity.TimeStamp, resultActivity.TimeStamp)
	}
	if resultActivity.MetaData["meta"] != activity.MetaData["meta"] {
		t.Error(activity.MetaData, resultActivity.MetaData)
	}
	if string(*resultActivity.Data) != string(*activity.Data) {
		t.Error(string(*activity.Data), string(*resultActivity.Data))
	}
}

func TestRankedFlatFeedScore(t *testing.T) {
	client := PreTestSetup(t)
	feed := getRankedFeed(t, client)

	activity := &getstream.Activity{
		Actor:  "flat:john",
		Verb:   "post",
		Object: "flat:eric",
	}
	_, err := feed.AddActivity(activity)
	if err != nil {
		t.Fatal(err)
	}

	testCases := []struct {
		name          string
		opts          getstream.FeedReadOptions
		shouldError   bool
		expectedScore bool
	}{
		{
			name:          "when ranking is not requested",
			shouldError:   false,
			opts:          getstream.NewFeedReadOptions(),
			expectedScore: false,
		},
		{
			name:        "when ranking is not configured",
			shouldError: true,
			opts:        getstream.NewFeedReadOptions().AddRanking("unknown"),
		},
		{
			name:          "when ranking is configured",
			shouldError:   false,
			opts:          getstream.NewFeedReadOptions().AddRanking("popularity"),
			expectedScore: true,
		},
	}

	for _, tc := range testCases {
		log.Println(tc.name) // this should be t.Run, but it's not supported on older versions of Go
		out, err := feed.Activities(tc.opts)
		if tc.shouldError && err == nil {
			t.Fatal("expected error, got nil")
		}
		if !tc.shouldError && err != nil {
			t.Fatalf("expected no errors, got %s", err)
			if len(out.Activities) != 1 {
				t.Fatalf("expected to have returned 1 activity, returned %d", len(out.Activities))
			}
			if tc.expectedScore && out.Activities[0].Score == nil {
				t.Fatalf("expected score, got nil")
			}
			if !tc.expectedScore && out.Activities[0].Score != nil {
				t.Fatalf("expected nil score, got %f", *out.Activities[0].Score)
			}
		}
	}
}

func TestFlatFeedUpdateActivities(t *testing.T) {
	client := PreTestSetup(t)
	feed := getFlatFeed(t, client)

	tt, _ := time.Parse("2006-01-02T15:04:05.999999", "2017-08-31T08:29:21.151279")
	err := feed.UpdateActivities(
		&getstream.Activity{
			ForeignID: "bob:123",
			TimeStamp: &tt,
			Actor:     "bob",
			Verb:      "verb",
			Object:    "object",
		},
	)

	assert.NoError(t, err)
}
