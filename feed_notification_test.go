package getstream_test

import (
	"encoding/json"
	"testing"
	"time"

	getstream "github.com/GetStream/stream-go"
)

func getNotificationFeed(client *getstream.Client) *getstream.NotificationFeed {
	f, _ := client.NotificationFeed("notification", RandString(8))
	return f
}

func TestExampleNotificationFeed_AddActivity(t *testing.T) {
	client := PreTestSetup(t)
	feed := getNotificationFeed(client)

	_, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: RandString(8),
		Object:    "flat:eric",
		Actor:     "flat:john",
	})
	if err != nil {
		t.Error(err)
	}
}

func TestNotificationFeedAddActivity(t *testing.T) {
	client := PreTestSetup(t)
	feed := getNotificationFeed(client)

	activity, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: RandString(8),
		Object:    "flat:eric",
		Actor:     "flat:john",
	})
	if err != nil {
		t.Error(err)
	}

	if activity.Verb != "post" && activity.ForeignID != "48d024fe-3752-467a-8489-23febd1dec4e" {
		t.Fail()
	}

}

func TestNotificationFeedAddActivityWithTo(t *testing.T) {
	client := PreTestSetup(t)
	feed := getNotificationFeed(client)
	feedTo := getNotificationFeed(client)

	_, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: RandString(8),
		Object:    "flat:eric",
		Actor:     "flat:john",
		To:        []getstream.FeedID{feedTo.FeedID()},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestNotificationFeedRemoveActivity(t *testing.T) {
	client := PreTestSetup(t)
	feed := getNotificationFeed(client)

	activity, err := feed.AddActivity(&getstream.Activity{
		Verb:   "post",
		Object: "flat:eric",
		Actor:  "flat:john",
	})
	if err != nil {
		t.Error(err)
	}

	if activity.Verb != "post" {
		t.Fail()
	}

	err = feed.RemoveActivity(activity.ID)
	if err != nil {
		t.Fatal(err)
	}
}

func TestNotificationFeedRemoveByForeignIDActivity(t *testing.T) {
	client := PreTestSetup(t)
	feed := getNotificationFeed(client)

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
		t.Fatal(err)
	}

}

func TestNotificationFeedActivities(t *testing.T) {
	client := PreTestSetup(t)
	feed := getNotificationFeed(client)

	_, err := feed.AddActivity(&getstream.Activity{
		Verb:      "post",
		ForeignID: RandString(8),
		Object:    "flat:eric",
		Actor:     "flat:john",
	})
	if err != nil {
		t.Error(err)
	}

	_, err = feed.Activities(getstream.NewNotificationFeedReadOptions())
	if err != nil {
		t.Error(err)
	}

}

func TestNotificationFeedAddActivities(t *testing.T) {
	client := PreTestSetup(t)
	feed := getNotificationFeed(client)

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

func TestNotificationFeedFollow(t *testing.T) {
	client := PreTestSetup(t)

	feedA := getNotificationFeed(client)
	feedB := getFlatFeed(t, client)

	if err := feedA.FollowFeedWithCopyLimit(feedB, 20); err != nil {
		t.Fatal(err)
	}

	if err := feedA.Unfollow(feedB); err != nil {
		t.Fatal(err)
	}

}

func TestNotificationFeedFollowKeepingHistory(t *testing.T) {
	client := PreTestSetup(t)

	feedA := getNotificationFeed(client)
	feedB := getFlatFeed(t, client)

	err := feedA.FollowFeedWithCopyLimit(feedB, 20)
	if err != nil {
		t.Fatal(err)
	}

	err = feedA.UnfollowKeepingHistory(feedB)
	if err != nil {
		t.Fatal(err)
	}

}

func TestMarkAsSeen(t *testing.T) {
	client := PreTestSetup(t)

	feed := getNotificationFeed(client)

	feed.AddActivities([]*getstream.Activity{
		{
			Actor:  "flat:larry",
			Object: "notification:larry",
			Verb:   "post",
		},
	})

	output, _ := feed.Activities(getstream.NewNotificationFeedReadOptions())
	if output.Unseen == 0 {
		t.Fail()
	}

	feed.Activities(getstream.NewNotificationFeedReadOptions().MarkAllSeen())

	output, _ = feed.Activities(getstream.NewNotificationFeedReadOptions())
	if output.Unseen != 0 {
		t.Fail()
	}

}

//func TestMarkAsRead(t *testing.T) {
//	client := PreTestSetup(t)
//
//	feed := getNotificationFeed(client)
//
//	feed.AddActivities([]*getstream.Activity{
//		{
//			Actor:  "flat:larry",
//			Object: "notification:larry",
//			Verb:   "post",
//		},
//	})
//
//	output, _ := feed.Activities(getstream.NewNotificationFeedReadOptions())
//	if output.Unread == 0 {
//		t.Fail()
//	}
//
//	for _, result := range output.Results {
//		err := feed.MarkActivitiesAsRead(result.Activities)
//		if err != nil {
//			t.Fatal(err)
//		}
//	}
//
//	output, _ = feed.Activities(getstream.NewNotificationFeedReadOptions())
//	if output.Unread != 0 {
//		t.Fail()
//	}
//}

func TestNotificationActivityMetaData(t *testing.T) {

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
		t.Fatal(err)
	}

	resultActivity2 := getstream.Activity{}
	err = json.Unmarshal(b2, &resultActivity2)
	if err != nil {
		t.Fatal(err)
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
