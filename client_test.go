package getstream_test

import (
	"testing"
	"time"

	getstream "github.com/GetStream/stream-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient(t *testing.T) {
	_, err := getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  "us-east",
	})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClientMissingAPIKey(t *testing.T) {
	_, err := getstream.New(&getstream.Config{
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  "us-east",
	})
	if err == nil {
		t.Fatal(err)
	}
	if err.Error() != "Required API Key was not set" {
		t.Fatal("Expected to get an error about missing APIKey")
	}
}

func TestClientMissingAPISecretAndToken(t *testing.T) {
	_, err := getstream.New(&getstream.Config{
		APIKey:   "my_secret",
		AppID:    "111111",
		Location: "us-east",
	})
	if err == nil {
		t.Fatal(err)
	}
	if err.Error() != "API Secret or Token was not set, one or the other is required" {
		t.Fatal("Expected to get an error about missing APISecret or Token")
	}
}

func TestClientLocalhost(t *testing.T) {
	client, err := getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  "localhost",
	})
	if err != nil {
		t.Fatal(err)
	}
	if client.BaseURL.String() != "http://localhost:8000/api/v1.0/" {
		t.Fatal("Location=localhost should be represented in non-SSL URL on port 8000, got", client.BaseURL.String())
	}
}

func TestClientToken(t *testing.T) {
	client, err := getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  "localhost",
	})
	if err != nil {
		t.Fatal(err)
	}
	if client.BaseURL.String() != "http://localhost:8000/api/v1.0/" {
		t.Fatal("Location=localhost should be represented in non-SSL URL on port 8000, got", client.BaseURL.String())
	}
}

func TestClient_FlatFeed(t *testing.T) {

	client, err := getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  "us-east"})
	if err != nil {
		t.Fatal(err)
	}

	feed, err := client.FlatFeed("flat", "UserID")
	if err != nil {
		t.Fatal(err)
	}

	_ = feed
}

func TestClient_FlatFeedBadSlug(t *testing.T) {

	client, err := getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  "us-east"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.FlatFeed("", "UserID")
	if err == nil {
		t.Fatal("Expected a slug validation error")
	}
	if err.Error() != "invalid feedSlug" {
		t.Fatal("Expected error about bad FlatFeed slug mismatch, got:", err.Error())
	}
}

func TestClient_NotificationFeed(t *testing.T) {
	client, err := getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  "us-east"})
	if err != nil {
		t.Fatal(err)
	}

	feed, err := client.NotificationFeed("flat", "UserID")
	if err != nil {
		t.Fatal(err)
	}

	_ = feed
}

func TestClient_NotificationFeedBadSlug(t *testing.T) {
	client, err := getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  "us-east"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.NotificationFeed("", "UserID")
	if err == nil {
		t.Fatal("Expected a slug validation error")
	}
	if err.Error() != "invalid feedSlug" {
		t.Fatal("Expected error about bad NotificationFeed slug mismatch, got:", err.Error())
	}
}

func TestClient_AggregatedFeed(t *testing.T) {
	client, err := getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  "us-east"})
	if err != nil {
		t.Fatal(err)
	}

	feed, err := client.AggregatedFeed("aggregated", "UserID")
	if err != nil {
		t.Fatal(err)
	}

	_ = feed
}

func TestClient_AggregatedFeedBadSlug(t *testing.T) {
	client, err := getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  "us-east"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.AggregatedFeed("", "UserID")
	if err == nil {
		t.Fatal("Expected a slug validation error")
	}
	if err.Error() != "invalid feedSlug" {
		t.Fatal("Expected error about bad AggregatedFeed slug mismatch, got:", err.Error())
	}
}

func TestClient_AggregatedFeedBadUserID(t *testing.T) {
	client, err := getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  "us-east"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.AggregatedFeed("aggregated", "")
	if err == nil {
		t.Fatal("Expected a userId validation error")
	}
	if err.Error() != "invalid userID" {
		t.Fatal("Expected error about bad AggregatedFeed userId mismatch, got:", err.Error())
	}
}

func TestFlatFeedInputValidation(t *testing.T) {
	client, err := getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  "us-east"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.FlatFeed("user", RandString(8))
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.FlatFeed("user", "tester@mail.com")
	if err == nil {
		t.Fatal(err)
	}
}

func TestNotificationFeedInputValidation(t *testing.T) {
	client, err := getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  "us-east"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.NotificationFeed("user", RandString(8))
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.NotificationFeed("user", "tester@mail.com")
	if err == nil {
		t.Fatal(err)
	}
}

func TestClientInit(t *testing.T) {
	_, err := getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  "!#@#$%ˆ&*((*=/*-+[]',.><"})
	if err == nil {
		t.Fatal(err)
	}

	_, err = getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  ""})
	if err != nil {
		t.Fatal(err)
	}

	_, err = getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  "us-east"})
	if err != nil {
		t.Fatal(err)
	}
}

func TestClientBaseURL(t *testing.T) {
	client, err := getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  "us-east"})
	if err != nil {
		t.Fatal(err)
	}

	if string(client.BaseURL.String()) != "https://us-east-api.getstream.io/api/v1.0/" {
		t.Fatal()
	}
}

func TestClientAbsoluteURL(t *testing.T) {
	client, err := getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  "us-east"})
	if err != nil {
		t.Fatal(err)
	}

	url, err := client.AbsoluteURL("user")
	if err != nil {
		t.Fatal(err)
	}

	if url.String() != "https://us-east-api.getstream.io/api/v1.0/user?api_key=my_key&location=us-east" {
		t.Fatal(err)
	}

	client, err = getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  ""})
	if err != nil {
		t.Fatal(err)
	}

	url, err = client.AbsoluteURL("flat")
	if err != nil {
		t.Fatal(err)
	}

	if url.String() != "https://api.getstream.io/api/v1.0/flat?api_key=my_key&location=unspecified" {
		t.Fatal()
	}

	client, err = getstream.New(&getstream.Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  ""})
	if err != nil {
		t.Fatal(err)
	}

	url, err = client.AbsoluteURL("!#@#$%ˆ&*((*=/*-+[]',.><")
	if err == nil {
		t.Fatal(err)
	}
}

func TestAddActivityToMany(t *testing.T) {
	client := PreTestSetup(t)

	feeds := []string{}

	bobFeed, err := client.FlatFeed("flat", "bob")
	if err != nil {
		t.Fatal(err)
	}
	feeds = append(feeds, string(bobFeed.FeedID()))

	sallyFeed, err := client.FlatFeed("flat", "sally")
	if err != nil {
		t.Fatal(err)
	}
	feeds = append(feeds, string(sallyFeed.FeedID()))

	activity := &getstream.Activity{
		Verb:      "post",
		ForeignID: RandString(8),
		Object:    "flat:eric",
		Actor:     "flat:john",
	}

	err = client.AddActivityToMany(*activity, feeds)
	if err != nil {
		t.Fatal(err)
	}

	// cleanup
	err = sallyFeed.RemoveActivityByForeignID(activity.ForeignID)
	if err != nil {
		t.Error(err)
	}
	err = bobFeed.RemoveActivityByForeignID(activity.ForeignID)
	if err != nil {
		t.Error(err)
	}
}

func prepareTestFollowMany(t *testing.T, client *getstream.Client) (*getstream.FlatFeed, []*getstream.PostFlatFeedFollowingManyInput, []*getstream.FlatFeed, *getstream.Activity) {
	flat, err := client.FlatFeed("flat", RandString(10))
	require.NoError(t, err)

	activity := &getstream.Activity{Actor: "bob", Verb: "like", Object: "cake"}
	activity, err = flat.AddActivity(activity)
	require.NoError(t, err)

	relations := make([]*getstream.PostFlatFeedFollowingManyInput, 20)
	feeds := make([]*getstream.FlatFeed, 20)
	for i := range relations {
		feeds[i], err = client.FlatFeed("flat", RandString(10))
		require.NoError(t, err)
		relations[i] = client.PrepFollowFlatFeed(flat, feeds[i])
	}
	return flat, relations, feeds, activity
}

func TestFollowMany(t *testing.T) {
	client := PreTestSetup(t)
	flat, relations, feeds, activity := prepareTestFollowMany(t, client)

	err := client.FollowMany(relations)
	assert.NoError(t, err)
	time.Sleep(2 * time.Second)

	followers, err := flat.GetFollowers(100, 0)
	require.NoError(t, err)
	assert.Len(t, followers, len(relations))
	time.Sleep(2 * time.Second)

	out, err := feeds[0].Activities(getstream.NewFeedReadOptions())
	require.NoError(t, err)
	assert.Len(t, out.Activities, 1)
	assert.Equal(t, out.Activities[0].ID, activity.ID)
}

func TestFollowManyCopyLimit(t *testing.T) {
	client := PreTestSetup(t)
	flat, relations, feeds, _ := prepareTestFollowMany(t, client)

	err := client.FollowManyCopyLimit(relations, 0)
	assert.NoError(t, err)
	time.Sleep(2 * time.Second)

	followers, err := flat.GetFollowers(100, 0)
	require.NoError(t, err)
	assert.Len(t, followers, len(relations))

	out, err := feeds[0].Activities(getstream.NewFeedReadOptions())
	require.NoError(t, err)
	assert.Len(t, out.Activities, 0)
}

func Test_ClientLocation(t *testing.T) {
	testCases := []struct {
		location string
		expected string
		version  string
	}{
		{
			location: "",
			expected: "https://api.getstream.io/api/v1.0/",
		},
		{
			location: "localhost",
			expected: "http://localhost:8000/api/v1.0/",
		},
		{
			location: "example",
			expected: "https://api-example.getstream.io/api/v1.0/",
		},
		{
			location: "example",
			version:  "v42",
			expected: "https://api-example.getstream.io/api/v42/",
		},
		{
			location: "example",
			expected: "https://api-example.getstream.io/api/v1.0/",
		},
		{
			location: "http://example.com",
			expected: "http://example.com/api/v1.0/",
		},
		{
			location: "https://example.com:1234/",
			expected: "https://example.com:1234/api/v1.0/",
		},
	}

	for _, tc := range testCases {
		cfg := &getstream.Config{
			APIKey:    "key",
			APISecret: "secret",
			Location:  tc.location,
		}
		if tc.version != "" {
			cfg.Version = tc.version
		}
		client, err := getstream.New(cfg)
		require.NoError(t, err)
		assert.Equal(t, tc.expected, client.BaseURL.String())
	}
}
