package getstream

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlatFeedBasic(t *testing.T) {
	client, err := New(&Config{
		APIKey:    "a key",
		APISecret: "a secret",
		AppID:     "11111",
		Location:  "us-east"})
	if err != nil {
		t.Fatal(err)
	}

	flatFeed := FlatFeed{
		baseFeed{
			Client:   client,
			FeedSlug: "feedGroup",
			UserID:   "feedName",
		},
	}

	if "feedGroupfeedName" != flatFeed.Signature() {
		t.Fatal()
	}

	if "feedGroup:feedName" != string(flatFeed.FeedID()) {
		t.Fatal()
	}

	flatFeed.SignFeed(flatFeed.Client.Signer)
	if "NWH8lcFHfHYEc2xdMs2kOhM-oII" != flatFeed.Token() {
		t.Fatal()
	}

	if "NWH8lcFHfHYEc2xdMs2kOhM-oII" != flatFeed.GenerateToken(flatFeed.Client.Signer) {
		t.Fatal()
	}

	if "feedGroupfeedName NWH8lcFHfHYEc2xdMs2kOhM-oII" != flatFeed.Signature() {
		t.Fatal()
	}
}

func TestNotificationFeedBasic(t *testing.T) {
	client, err := New(&Config{
		APIKey:    "a key",
		APISecret: "a secret",
		AppID:     "11111",
		Location:  "us-east"})
	if err != nil {
		t.Fatal(err)
	}

	notificationFeed := NotificationFeed{
		baseFeed{
			Client:   client,
			FeedSlug: "feedGroup",
			UserID:   "feedName",
		},
	}

	if "feedGroupfeedName" != notificationFeed.Signature() {
		t.Fatal()
	}

	if "feedGroup:feedName" != string(notificationFeed.FeedID()) {
		t.Fatal()
	}

	notificationFeed.SignFeed(notificationFeed.Client.Signer)

	if "NWH8lcFHfHYEc2xdMs2kOhM-oII" != notificationFeed.Token() {
		t.Fatal()
	}

	if "NWH8lcFHfHYEc2xdMs2kOhM-oII" != notificationFeed.GenerateToken(notificationFeed.Client.Signer) {
		t.Fatal()
	}

	if "feedGroupfeedName NWH8lcFHfHYEc2xdMs2kOhM-oII" != notificationFeed.Signature() {
		t.Fatal()
	}
}

func TestUpdateActivities(t *testing.T) {
	client, err := New(&Config{
		APIKey:    os.Getenv("STREAM_API_KEY"),
		APISecret: os.Getenv("STREAM_API_SECRET"),
		AppID:     os.Getenv("STREAM_APP_ID"),
		Location:  os.Getenv("STREAM_REGION"),
	})
	require.Nil(t, err)

	feed, err := client.FlatFeed("flat", "123")
	require.Nil(t, err)

	tt, _ := time.Parse("2006-01-02T15:04:05.999999", "2017-08-31T08:29:21.151279")

	testCases := []struct {
		tt          time.Time
		activities  []*Activity
		shouldError bool
	}{
		{
			activities: []*Activity{
				&Activity{
					ForeignID: "bob:123",
					TimeStamp: &tt,
					Actor:     "bob",
					Verb:      "verb",
					Object:    "object",
				},
			},
			shouldError: false,
		},
		{
			activities: []*Activity{
				&Activity{
					TimeStamp: &tt,
					Actor:     "bob",
					Verb:      "verb",
					Object:    "object",
				},
			},
			shouldError: true,
		},
		{
			activities: []*Activity{
				&Activity{
					ForeignID: "bob:123",
					Actor:     "bob",
					Verb:      "verb",
					Object:    "object",
				},
			},
			shouldError: true,
		},
		{
			activities:  []*Activity{},
			shouldError: true,
		},
		{
			activities: []*Activity{
				&Activity{
					ForeignID: "alice:123",
					TimeStamp: &tt,
					Actor:     "alice",
					Verb:      "verb",
					Object:    "object",
				},
			},
			shouldError: false,
		},
	}

	for _, tc := range testCases {
		err := feed.UpdateActivities(tc.activities...)
		if tc.shouldError {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
		}
	}
}
