package getstream_test

import (
	"encoding/json"
	"testing"
	"time"

	getstream "github.com/GetStream/stream-go"
	"github.com/pborman/uuid"
)

func TestExampleFlatFeedAddActivity(t *testing.T) {

	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feed, err := client.FlatFeed("flat", "UserID")
	if err != nil {
		t.Fatal(err)
	}

	activity, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: uuid.New(),
		Object:    "flat:eric",
		Actor:     "flat:john",
	})
	if err != nil {
		t.Fatal(err)
	}

	_ = activity
}

func TestFlatFeedAddActivity(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feed, err := client.FlatFeed("flat", "bob")
	if err != nil {
		t.Fatal(err)
	}

	activity, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: uuid.New(),
		Object:    "flat:eric",
		Actor:     "flat:john",
	})
	if err != nil {
		t.Fatal(err)
	}

	if activity.Verb != "post" && activity.ForeignID != "48d024fe-3752-467a-8489-23febd1dec4e" {
		t.Fail()
	}

	err = PostTestCleanUp(client, []*getstream.Activity{activity}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFlatFeedAddActivityWithTo(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feed, err := client.FlatFeed("flat", "bob")
	if err != nil {
		t.Fatal(err)
	}

	feedTo, err := client.FlatFeed("flat", "barry")
	if err != nil {
		t.Fatal(err)
	}

	feedToB, err := client.FlatFeed("flat", "larry")
	if err != nil {
		t.Fatal(err)
	}

	activity, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: uuid.New(),
		Object:    "flat:eric",
		Actor:     "flat:john",
		To:        []getstream.FeedID{feedTo.FeedID(), feedToB.FeedID()},
	})
	if err != nil {
		t.Fatal(err)
	}

	if activity.Verb != "post" && activity.ForeignID != "48d024fe-3752-467a-8489-23febd1dec4e" {
		t.Fail()
	}

	err = PostTestCleanUp(client, []*getstream.Activity{activity}, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFlatFeedUUID(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feed, err := client.FlatFeed("flat", "bob")
	if err != nil {
		t.Fatal(err)
	}

	var activities []*getstream.Activity

	for i := 0; i < 10; i++ {
		activity, err := feed.AddActivity(&getstream.Activity{
			Verb:      "post",
			ForeignID: uuid.New(),
			Object:    "flat:eric",
			Actor:     "flat:john",
		})
		if err != nil {
			t.Log("fail add activity with UUID :", err)
			continue
		}

		err = feed.RemoveActivityByForeignID(activity.ForeignID)
		if err != nil {
			t.Log("fail remove activity with UUID :", err)
		}

		activities = append(activities, activity)
	}

	PostTestCleanUp(client, activities, nil, nil)
}

func TestFlatFeedRemoveActivity(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feed, err := client.FlatFeed("flat", "bob")
	if err != nil {
		t.Fatal(err)
	}

	activity, err := feed.AddActivity(&getstream.Activity{
		Verb:   "post",
		Object: "flat:eric",
		Actor:  "flat:john",
	})
	if err != nil {
		t.Fatal(err)
	}

	if activity.Verb != "post" {
		t.Fail()
	}

	err = feed.RemoveActivity(activity.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFlatFeedRemoveByForeignIDActivity(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feed, err := client.FlatFeed("flat", "bob")
	if err != nil {
		t.Fatal(err)
	}

	activity, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: uuid.New(),
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

	PostTestCleanUp(client, []*getstream.Activity{activity}, nil, nil)
}

func TestFlatFeedGetActivities(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feed, err := client.FlatFeed("flat", "bob")
	if err != nil {
		t.Fatal(err)
	}

	_, err = feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: uuid.New(),
		Object:    "flat:eric",
		Actor:     "flat:john",
	})
	if err != nil {
		t.Error(err)
	}

	activities, err := feed.Activities(&getstream.GetFlatFeedInput{})
	if err != nil {
		t.Error(err)
	}

	if activities.Activities[0].Actor != "flat:john" {
		t.Error("Activity read from stream did not match")
	}

	err = PostTestCleanUp(client, activities.Activities, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFlatFeedAddActivities(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feed, err := client.FlatFeed("flat", "bob")
	if err != nil {
		t.Fatal(err)
	}

	activities, err := feed.AddActivities([]*getstream.Activity{
		{
			Verb:      "post",
			ForeignID: uuid.New(),
			Object:    "flat:eric",
			Actor:     "flat:john",
		}, {
			Verb:      "walk",
			ForeignID: uuid.New(),
			Object:    "flat:john",
			Actor:     "flat:eric",
		},
	})
	if err != nil {
		t.Error(err)
	}

	err = PostTestCleanUp(client, activities, nil, nil)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFlatFeedFollow(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feedA, err := client.FlatFeed("flat", "bob")
	if err != nil {
		t.Fatal(err)
	}

	feedB, err := client.FlatFeed("flat", "eric")
	if err != nil {
		t.Fatal(err)
	}

	err = feedA.FollowFeedWithCopyLimit(feedB, 20)
	if err != nil {
		t.Error(err)
	}

	// get feedB's followers, ensure feedA is there
	followers, err := feedB.GetFollowers(5, 0)
	if err != nil {
		t.Error(err)
	}
	if followers[0].UserID != "bob" {
		t.Error("Bob's FeedA is not a follower of FeedB")
	}

	// get things that feedA follows, ensure feedB is in there
	following, err := feedA.GetFollowings(5, 0)
	if err != nil {
		t.Error(err)
	}
	if following[0].UserID != "eric" {
		t.Error("Eric's FeedB is not a follower of FeedA")
	}

	PostTestCleanUpFlatFeedFollows(client, []*getstream.FlatFeed{feedA, feedB})
}

func TestFlatFeedFollowingFollowers(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feedA, err := client.FlatFeed("flat", "bob")
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

	_, err = feedA.GetFollowings(20, 0)
	if err != nil {
		t.Fail()
	}

	_, err = feedB.GetFollowers(20, 0)
	if err != nil {
		t.Fail()
	}

	PostTestCleanUpFlatFeedFollows(client, []*getstream.FlatFeed{feedA, feedB, feedC})
}

func TestFlatFeedUnFollow(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feedA, err := client.FlatFeed("flat", "bob")
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

	err = feedA.Unfollow(feedB)
	if err != nil {
		t.Fail()
	}

	PostTestCleanUpFlatFeedFollows(client, []*getstream.FlatFeed{feedA, feedB})
}

func TestFlatFeedUnFollowKeepingHistory(t *testing.T) {
	client, err := PreTestSetup()
	if err != nil {
		t.Fatal(err)
	}

	feedA, err := client.FlatFeed("flat", "bob")
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

	PostTestCleanUpFlatFeedFollows(client, []*getstream.FlatFeed{feedA, feedB})
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
		ForeignID: uuid.New(),
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
