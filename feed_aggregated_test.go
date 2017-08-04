package getstream_test

import (
	"encoding/json"
	"testing"
	"time"

	getstream "github.com/GetStream/stream-go"
)

func getAggregatedFeed(client *getstream.Client) *getstream.AggregatedFeed {
	f, _ := client.AggregatedFeed("aggregated", RandString(8))
	return f
}

func TestExampleAggregatedFeed_AddActivity(t *testing.T) {
	client := PreTestSetup(t)
	feed := getAggregatedFeed(client)

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

func TestAggregatedFeedAddActivity(t *testing.T) {
	client := PreTestSetup(t)
	feed := getAggregatedFeed(client)

	fid := RandString(8)
	_, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: fid,
		Object:    "flat:eric",
		Actor:     "flat:john",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestAggregatedFeedAddActivityWithTo(t *testing.T) {
	client := PreTestSetup(t)
	feed := getAggregatedFeed(client)
	toFeed := getAggregatedFeed(client)

	_, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: RandString(8),
		Object:    "flat:eric",
		Actor:     "flat:john",
		To:        []getstream.FeedID{toFeed.FeedID()},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestAggregatedFeedRemoveActivity(t *testing.T) {
	client := PreTestSetup(t)
	feed := getAggregatedFeed(client)

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
		t.Fatal(err)
	}
}

func TestAggregatedFeedRemoveByForeignIDActivity(t *testing.T) {
	client := PreTestSetup(t)
	feed := getAggregatedFeed(client)

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

	rmActivity := getstream.Activity{
		ForeignID: activity.ForeignID,
	}
	_ = rmActivity

	err = feed.RemoveActivityByForeignID(activity.ID)
	if err != nil {
		t.Fatal(err)
	}

}

func TestAggregatedFeedActivities(t *testing.T) {
	client := PreTestSetup(t)
	feed := getAggregatedFeed(client)

	_, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: RandString(8),
		Object:    "flat:eric",
		Actor:     "flat:john",
	})
	if err != nil {
		t.Error(err)
	}

	_, err = feed.Activities(getstream.NewFeedReadOptions())
	if err != nil {
		t.Fatal(err)
	}

}

func TestAggregatedFeedAddActivities(t *testing.T) {
	client := PreTestSetup(t)
	feed := getAggregatedFeed(client)

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
		t.Error(err)
	}

}

func TestAggregatedFeedFollowUnfollow(t *testing.T) {
	client := PreTestSetup(t)
	feedA := getAggregatedFeed(client)
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
		t.Error("feedA aggregated feed is not a follower of FeedB")
	}

	// get things that feedA follows, ensure feedB is in there
	following, err := feedA.GetFollowings(5, 0)
	if err != nil {
		t.Error(err)
	}
	if following[0] != feedB.FeedID() {
		t.Error("Eric's FeedB is not a follower of FeedA")
	}

	if err := feedA.Unfollow(feedB); err != nil {
		t.Fatal(err)
	}

}

func TestAggregatedFeedFollowKeepingHistory(t *testing.T) {
	client := PreTestSetup(t)

	feedA := getAggregatedFeed(client)
	feedB := getFlatFeed(t, client)

	if err := feedA.FollowFeedWithCopyLimit(feedB, 20); err != nil {
		t.Fatal(err)
	}

	if err := feedA.UnfollowKeepingHistory(feedB); err != nil {
		t.Fatal(err)
	}

}

func TestAggregatedActivityMetaData(t *testing.T) {
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
		Origin:    "user:barry",
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
	if resultActivity.Origin != activity.Origin {
		t.Error(activity.Origin, resultActivity.Origin)
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
