package getstream_test

import (
	"encoding/json"
	"testing"
	"time"

	getstream "github.com/GetStream/stream-go"
)

func getFlatFeed(client *getstream.Client) *getstream.FlatFeed {
	f, _ := client.FlatFeed("flat", RandString(8))
	return f
}

func TestExampleFlatFeedAddActivity(t *testing.T) {
	client := PreTestSetup(t)
	feed := getFlatFeed(client)

	activity, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: RandString(8),
		Object:    "flat:eric",
		Actor:     "flat:john",
	})
	if err != nil {
		t.Fatal(err)
	}

	_ = activity
}

func TestFlatFeedAddActivity(t *testing.T) {
	client := PreTestSetup(t)
	feed := getFlatFeed(client)

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
	feed := getFlatFeed(client)

	feedTo := getFlatFeed(client)
	feedToB := getFlatFeed(client)

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
	feed := getFlatFeed(client)

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
	feed := getFlatFeed(client)

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
	feed := getFlatFeed(client)

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
	feed := getFlatFeed(client)

	_, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: RandString(8),
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

}

func TestFlatFeedAddActivities(t *testing.T) {
	client := PreTestSetup(t)
	feed := getFlatFeed(client)

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

func TestFlatFeedFollow(t *testing.T) {
	client := PreTestSetup(t)

	feedA := getFlatFeed(client)
	feedB := getFlatFeed(client)

	err := feedA.FollowFeedWithCopyLimit(feedB, 20)
	if err != nil {
		t.Error(err)
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
		t.Error(err)
	}
	if following[0] != feedB.FeedID() {
		t.Error("Eric's FeedB is not a follower of FeedA")
	}

}

func TestFlatFeedFollowingFollowers(t *testing.T) {
	client := PreTestSetup(t)

	feedA := getFlatFeed(client)
	feedB := getFlatFeed(client)
	feedC := getFlatFeed(client)

	err := feedA.FollowFeedWithCopyLimit(feedB, 20)
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

}

func TestFlatFeedUnFollow(t *testing.T) {
	client := PreTestSetup(t)

	feedA := getFlatFeed(client)
	feedB := getFlatFeed(client)

	err := feedA.FollowFeedWithCopyLimit(feedB, 20)
	if err != nil {
		t.Fail()
	}

	err = feedA.Unfollow(feedB)
	if err != nil {
		t.Fail()
	}

}

func TestFlatFeedUnFollowKeepingHistory(t *testing.T) {
	client := PreTestSetup(t)

	feedA := getFlatFeed(client)
	feedB := getFlatFeed(client)

	err := feedA.FollowFeedWithCopyLimit(feedB, 20)
	if err != nil {
		t.Fail()
	}

	err = feedA.UnfollowKeepingHistory(feedB)
	if err != nil {
		t.Fail()
	}

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
