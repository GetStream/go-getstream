# stream-go

[![godoc](https://godoc.org/github.com/GetStream/stream-go?status.svg)](https://godoc.org/github.com/stream-go)
[![codecov](https://codecov.io/gh/GetStream/stream-go/branch/master/graph/badge.svg)](https://codecov.io/gh/GetStream/stream-go)

Beta client library for [GetStream.io](//getstream.io).

This library could not exist without the efforts of several open-source community members, 
including the awesome folks at [MrHenry](//github.com/mrhenry) and 
[HyperWorks](//github.com/hyperworks)(HyperWorks). Thank you so much for contributing
to our community!

The code provided by the MrHenry team is used with written permission; we are working 
with them to get a final license in place. Stream will be modifying the codebase 
together with MrHenry over time, so we especially want to point out how great they 
have been working with us to release this library.

### Beta

We are releasing this as our v0.9.0 beta since there may be bugs, and inevitable cleanup 
will happen along the way. Please do not hesitate to [contact us](mailto:support@getstream.io)
if you see something strange happening. We'd be happy to consider any and all pull 
requests from the community as well!

### Example Usage

Creating a client:
```go
import (
	"fmt"
	"github.com/GetStream/stream-go"
)

// we recommend getting your API credentials using os.Getenv()
APIKey string = "your-api-key"
APISecret string = "your-api-secret"

// your application ID, found on your GetStream.io dashboard
AppID string = "16013"

// Location is optional; leaving it blank will default the
// hostname to "api.getstream.io"
// but we do have geographic-specific choices:
// "us-east", "us-west" and "eu-west"
Location string = "us-east"

// TimeoutInt is an optional integer parameter to define
// the number of seconds before your connection will hang-up
// during a request; you can set this to any non-negative
// and non-zero whole number, and will default to 3 
TimeoutInt: 3

client, err := getstream.New(&getstream.Config{
    APIKey:    os.Getenv("STREAM_APIKEY"),
    APISecret: os.Getenv("STREAM_APISECRET"),
    AppID:     os.Getenv("STREAM_APPID"),
    Location:  os.Getenv("STREAM_LOCATION"),
})
if err != nil {
    return err
}
```

Creating a Feed object for a user:

```go
// this code assumes you've created a flat feed named "flat-feed-name" for your app
// we also recommend using UUID values for users

bobFeed, err := client.FlatFeed("flat-feed-name", "bob-uuid")
if err != nil {
    return err
}
```

Creating an activity on Bob's flat feed:
```go
import "github.com/pborman/uuid"

activity, err := bobFeed.AddActivity(&Activity{
    Verb:      "post",
    ForeignID: uuid.New(),
    Object:    FeedID("flat:eric"),
    Actor:     FeedID("flat:john"),
})
if err != nil {
    return err
}
```

The library is gradually introducing JWT support. You can generate a client token
for a feed using the following example:

```go
client, err := getstream.New(&getstream.Config{
    APIKey:    os.Getenv("STREAM_APIKEY"),
    APISecret: os.Getenv("STREAM_APISECRET"),
    AppID:     os.Getenv("STREAM_APPID"),
    Location:  os.Getenv("STREAM_LOCATION"),
})

feed, err := client.FlatFeed("flat-feed-name", "bob-uuid")
if err != nil {
    return err
}

token, err := client.Signer.GenerateFeedScopeToken(
    getstream.ScopeContextFeed, 
    getstream.ScopeActionRead, 
    bobFeed)
if err != nil {
    fmt.Println(err)
}

// note in the struct below that we're not setting APISecret
// but setting "Token" instead:
clientSideClient, err := getstream.NewWithToken(&getstream.Config{
    APIKey:    os.Getenv("STREAM_APIKEY"),
    Token:     token, // not setting APISecret
    AppID:     os.Getenv("STREAM_APPID"),
    Location:  os.Getenv("STREAM_LOCATION"),
})
if err != nil {
  return err
}
```

JWT support is not yet fully tested on the library, but we'd love to 
hear any feedback you have if you try it out.

### API Support

Flat Feed

- [x] Add one or more Activities (AddActivity, AddActivities)
- [x] Remove Activity (RemoveActivity, RemoveActivityByForeignID)
- [x] Get a list of Activities on the Feed (Activities)
- [x] Follow another Feed (FollowFeedWithCopyLimit)
- [x] UnFollow another Feed (Unfollow, UnfollowAggregated, UnfollowNotification, UnfollowKeepingHistory)
- [x] Get Followers of this Feed (FollowersWithLimitAndSkip)
- [x] Get list of Feeds this Feed is Following (FollowingWithLimitAndSkip)
- [x] Follow Many Feeds (FollowManyFeeds)
- [x] Update one or more Activities (UpdateActivity, UpdateActivities)

Aggregated Feed

- [x] Add one or more Activities (AddActivity, AddActivities)
- [x] Remove Activity (RemoveActivity, RemoveActivityByForeignID)
- [x] Get a list of Activities on the Feed (Activities)
- [x] Follow another Feed (FollowFeedWithCopyLimit)
- [x] UnFollow another Feed (Unfollow, UnfollowKeepingHistory)
- [x] Get Followers of this Feed (FollowersWithLimitAndSkip)
- [x] Get list of Feeds this Feed is Following (FollowingWithLimitAndSkip)
- [ ] Update one or more Activities

Notification Feed

- [x] Add one or more Activities (AddActivity, AddActivities)
- [x] Remove Activity (RemoveActivity, RemoveActivityByForeignID)
- [x] Get a list of Activities on the Feed (Activities)
- [x] Follow another Feed (FollowFeedWithCopyLimit)
- [x] UnFollow another Feed (Unfollow, UnfollowKeepingHistory)
- [x] Get list of Feeds this Feed is Following (FollowingWithLimitAndSkip)
- [x] Mark Read (MarkActivitiesAsRead)
- [x] Mark Seen (MarkActivitiesAsSeenWithLimit)
- [x] Get Followers of this Feed (FollowersWithLimitAndSkip)
- [ ] Update one or more Activities

### Activity Payload Structure

Payload building Follows our API standards for all request payloads
- `data` : Statically typed payloads as `json.RawMessage`
- `metadata` : Top-level key/value pairs

You can/should use `data` to send Go structures through the library. This 
will give you the benefit of Go's static type system. If you are unable
to determine a type (or compatible type) for the contents of an Activity,
you can use `metadata` which is a `map[string]string`; encoding this to
JSON will move these values to the top-level, so any keys you define in
your `metadata` which conflict with our standard top-level keys will be
overwritten.

The benefit of this `metadata` structure is that these key/value pairs 
will be exposed to Stream's internals such as ranking.

### Design Choices

Many design choices in the library were inherited from the team at MrHenry,
with some choices to refactor some of the test code as its own getstream_test
package. This choice meant exposing some attributes that perhaps should
be left private, and we expect this re-refactoring will take place over
time.

The MrHenry team noted this about the Flat / Aggregated / Notification 
Feed types:
- they have separate structures and methods to prevent the impact of 
future Stream changes
- if two types of feeds grow farther apart, incorporated future changes
in this client should not breaking everything

### Credits

Have we mentioned the team at [MrHenry](//github.com/mrhenry) yet??

##### Credits from MrHenry that we wanted to pass along as well:

This pkg started out as a fork from [HyperWorks](//github.com/hyperworks/go-getstream)
and still borrows snippets (token & errors) from the original. We 
decided to make this a separate repo entirely because drastic changes 
were made to the interface and internal workings.

We received great support from Stream while creating this pkg for which
we are very grateful, and we also want to thank them for creating 
Stream in the first place.

