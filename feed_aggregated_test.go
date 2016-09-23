package getstream_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/pborman/uuid"
	getstream "github.com/GetStream/stream-go"
)

func TestExampleAggregatedFeed_AddActivity(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feed, err := client.AggregatedFeed("flat", "UserID")
	if err != nil {
		t.Fatal(err)
	}

	activity, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: uuid.New(),
		Object:    getstream.FeedID("flat:eric"),
		Actor:     getstream.FeedID("flat:john"),
	})
	if err != nil {
		t.Fatal(err)
	}

	_ = activity
}

func TestAggregatedFeedAddActivity(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feed, err := client.AggregatedFeed("aggregated", "bob")
	if err != nil {
		t.Fatal(err)
	}

	activity, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: uuid.New(),
		Object:    getstream.FeedID("flat:eric"),
		Actor:     getstream.FeedID("flat:john"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if activity.Verb != "post" && activity.ForeignID != "48d024fe-3752-467a-8489-23febd1dec4e" {
		t.Error(err)
	}

	err = PostTestCleanUp(client, nil, nil, []*getstream.Activity{activity})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAggregatedFeedAddActivityWithTo(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feed, err := client.AggregatedFeed("aggregated", "bob")
	if err != nil {
		t.Fatal(err)
	}

	toFeed, err := client.AggregatedFeed("aggregated", "barry")
	if err != nil {
		t.Fatal(err)
	}

	activity, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: uuid.New(),
		Object:    getstream.FeedID("flat:eric"),
		Actor:     getstream.FeedID("flat:john"),
		To:        []getstream.Feed{toFeed},
	})
	if err != nil {
		t.Fatal(err)
	}

	if activity.Verb != "post" && activity.ForeignID != "48d024fe-3752-467a-8489-23febd1dec4e" {
		t.Fail()
	}

	err = PostTestCleanUp(client, nil, nil, []*getstream.Activity{activity})
	if err != nil {
		t.Fatal(err)
	}
}

func TestAggregatedFeedRemoveActivity(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feed, err := client.AggregatedFeed("aggregated", "bob")
	if err != nil {
		t.Fatal(err)
	}

	activity, err := feed.AddActivity(&getstream.Activity{
		Verb:   "post",
		Object: getstream.FeedID("flat:eric"),
		Actor:  getstream.FeedID("flat:john"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if activity.Verb != "post" {
		t.Fail()
	}

	rmActivity := getstream.Activity{
		ID: activity.ID,
	}

	err = feed.RemoveActivity(&rmActivity)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAggregatedFeedRemoveByForeignIDActivity(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feed, err := client.AggregatedFeed("aggregated", "bob")
	if err != nil {
		t.Fatal(err)
	}

	activity, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: uuid.New(),
		Object:    getstream.FeedID("flat:eric"),
		Actor:     getstream.FeedID("flat:john"),
	})
	if err != nil {
		t.Fatal(err)
	}

	if activity.Verb != "post" && activity.ForeignID != "08f01c47-014f-11e4-aa8f-0cc47a024be0" {
		t.Fail()
	}

	rmActivity := getstream.Activity{
		ForeignID: activity.ForeignID,
	}
	_ = rmActivity

	err = feed.RemoveActivityByForeignID(activity)
	if err != nil {
		t.Fatal(err)
	}

	PostTestCleanUp(client, nil, nil, []*getstream.Activity{activity})
}

func TestAggregatedFeedActivities(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feed, err := client.AggregatedFeed("aggregated", "bob")
	if err != nil {
		t.Fatal(err)
	}

	_, err = feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: uuid.New(),
		Object:    getstream.FeedID("flat:eric"),
		Actor:     getstream.FeedID("flat:john"),
	})
	if err != nil {
		t.Error(err)
	}

	activities, err := feed.Activities(&getstream.GetAggregatedFeedInput{})
	if err != nil {
		t.Fatal(err)
	}

	for _, result := range activities.Results {
		err = PostTestCleanUp(client, nil, nil, result.Activities)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestAggregatedFeedAddActivities(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feed, err := client.AggregatedFeed("aggregated", "bob")
	if err != nil {
		t.Fatal(err)
	}

	activities, err := feed.AddActivities([]*getstream.Activity{
		{
			Verb:      "post",
			ForeignID: uuid.New(),
			Object:    getstream.FeedID("flat:eric"),
			Actor:     getstream.FeedID("flat:john"),
		}, {
			Verb:      "walk",
			ForeignID: uuid.New(),
			Object:    getstream.FeedID("flat:john"),
			Actor:     getstream.FeedID("flat:eric"),
		},
	})
	if err != nil {
		t.Error(err)
	}

	err = PostTestCleanUp(client, nil, nil, activities)
	if err != nil {
		t.Fatal(err)
	}
}

func TestAggregatedFeedFollowUnfollow(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feedA, err := client.AggregatedFeed("aggregated", "bob")
	if err != nil {
		t.Fatal(err)
	}

	feedB, err := client.FlatFeed("flat", "eric")
	if err != nil {
		t.Fatal(err)
	}

	err = feedA.FollowFeedWithCopyLimit(feedB, 20)
	if err != nil {
		t.Fail()
	}

	// get feedB's followers, ensure feedA is there
	followers, err := feedB.FollowersWithLimitAndSkip(5, 0)
	if err != nil {
		t.Error(err)
	}
	if followers[0].UserID != "bob" {
		t.Error("Bob's aggregated feed is not a follower of FeedB")
	}

	// get things that feedA follows, ensure feedB is in there
	following, err := feedA.FollowingWithLimitAndSkip(5, 0)
	if err != nil {
		t.Error(err)
	}
	if following[0].UserID != "eric" {
		t.Error("Eric's FeedB is not a follower of FeedA")
	}

	err = feedA.Unfollow(feedB)
	if err != nil {
		t.Fail()
	}

	PostTestCleanUpFlatFeedFollows(client, []*getstream.FlatFeed{feedB})
}

func TestAggregatedFeedFollowKeepingHistory(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feedA, err := client.AggregatedFeed("aggregated", "bob")
	if err != nil {
		t.Fatal(err)
	}

	feedB, err := client.FlatFeed("flat", "eric")
	if err != nil {
		t.Fatal(err)
	}

	err = feedA.FollowFeedWithCopyLimit(feedB, 20)
	if err != nil {
		t.Fail()
	}

	err = feedA.UnfollowKeepingHistory(feedB)
	if err != nil {
		t.Fail()
	}

	PostTestCleanUpFlatFeedFollows(client, []*getstream.FlatFeed{feedB})
}

func TestAggregatedFeedFollowingFollowers(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feedA, err := client.AggregatedFeed("aggregated", "bob")
	if err != nil {
		t.Fatal(err)
	}

	feedB, err := client.FlatFeed("flat", "eric")
	if err != nil {
		t.Fatal(err)
	}

	feedC, err := client.FlatFeed("flat", "barry")
	if err != nil {
		t.Fatal(err)
	}

	err = feedA.FollowFeedWithCopyLimit(feedB, 20)
	if err != nil {
		t.Fail()
	}

	err = feedA.FollowFeedWithCopyLimit(feedC, 20)
	if err != nil {
		t.Fail()
	}

	_, err = feedA.FollowingWithLimitAndSkip(20, 0)
	if err != nil {
		t.Fail()
	}

	PostTestCleanUpFlatFeedFollows(client, []*getstream.FlatFeed{feedB, feedC})
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
		ForeignID: uuid.New(),
		Actor:     getstream.FeedID("user:eric"),
		Object:    getstream.FeedID("user:bob"),
		Target:    getstream.FeedID("user:john"),
		Origin:    getstream.FeedID("user:barry"),
		Verb:      "post",
		TimeStamp: &now,
		Data:      &raw,
		MetaData: map[string]string{
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
