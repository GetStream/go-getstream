package getstream_test

import (
	"fmt"
	"testing"

	"math/rand"
	"time"

	getstream "github.com/GetStream/stream-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var src = rand.NewSource(time.Now().UnixNano())

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

func RandString(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}

func TestFeedSlug(t *testing.T) {
	feedSlug, err := getstream.ValidateFeedSlug("foo")
	if err != nil {
		t.Error(err)
	}
	if feedSlug != "foo" {
		t.Error("feedSlug not 'foo'")
	}

	feedSlug, err = getstream.ValidateFeedSlug("f-o-o")
	if err != nil {
		t.Error(err)
	}
	if feedSlug != "f_o_o" {
		t.Error("feedSlug not 'f_o_o'")
	}
}

func TestFeedID(t *testing.T) {
	feedID, err := getstream.ValidateFeedID("123")
	if err != nil {
		t.Error(err)
	}
	if feedID != "123" {
		t.Error("feedID not '123'")
	}

	feedID, err = getstream.ValidateFeedID("1-2-3")
	if err != nil {
		t.Error(err)
	}
	if feedID != "1_2_3" {
		t.Error("feedID not '1_2_3'")
	}
}

func TestUserID(t *testing.T) {
	userID, err := getstream.ValidateUserID("123")
	if err != nil {
		t.Error(err)
	}
	if userID != "123" {
		t.Error("userID not '123'")
	}

	userID, err = getstream.ValidateUserID("1-2-3")
	if err != nil {
		t.Error(err)
	}
	if userID != "1_2_3" {
		t.Error("userSlug not '1_2_3'")
	}
}

func prepareUnfollowKeepingHistory(t *testing.T, feedA, feedB getstream.Feed, fn func()) {
	activities := make([]*getstream.Activity, 10)
	for i := range activities {
		activities[i] = &getstream.Activity{Actor: "test", Verb: "like", Object: fmt.Sprintf("obj-%d", i)}
	}
	_, err := feedB.AddActivities(activities)
	require.NoError(t, err)

	err = feedA.FollowFeedWithCopyLimit(feedB, 20)
	require.NoError(t, err)

	fn()
	time.Sleep(3 * time.Second)

	err = feedA.UnfollowKeepingHistory(feedB)
	assert.NoError(t, err)
}
