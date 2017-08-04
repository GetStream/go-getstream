package getstream_test

import (
	"testing"

	"github.com/GetStream/stream-go"
)

func TestGenerateToken(t *testing.T) {

	signer := getstream.Signer{
		Secret: "test_secret",
	}

	token := signer.GenerateToken("some message")
	if token != "8SZVOYgCH6gy-ZjBTq_9vydr7TQ" {
		t.Fail()
		return
	}
}

func TestURLSafe(t *testing.T) {

	signer := getstream.Signer{}

	result := signer.UrlSafe("some+test/string=foo=")
	if result != "some-test_string=foo" {
		t.Fail()
		return
	}
}

func TestFeedScopeToken(t *testing.T) {

	client, err := getstream.New(&getstream.Config{
		APIKey:    "a_key",
		APISecret: "tfq2sdqpj9g446sbv653x3aqmgn33hsn8uzdc9jpskaw8mj6vsnhzswuwptuj9su",
		AppID:     "123456",
		Location:  "us-east"})
	if err != nil {
		t.Fatal(err)
		return
	}

	feed, err := client.FlatFeed("flat", "bob")
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Signer.GenerateJWT(
		getstream.ScopeContextFeed,
		getstream.ScopeActionRead,
		feed.FeedIDWithoutColon())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Signer.GenerateJWT(
		getstream.ScopeContextActivities,
		getstream.ScopeActionWrite,
		feed.FeedIDWithoutColon())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Signer.GenerateJWT(
		getstream.ScopeContextFollower,
		getstream.ScopeActionDelete,
		feed.FeedIDWithoutColon())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Signer.GenerateJWT(
		getstream.ScopeContextAll,
		getstream.ScopeActionAll,
		feed.FeedIDWithoutColon())
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.Signer.GenerateJWT(getstream.ScopeContextFeed, getstream.ScopeActionRead, "")
	if err != nil {
		t.Fatal(err)
	}
}
