package getstream

import (
	"testing"
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
