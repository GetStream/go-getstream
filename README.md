# stream-go

[stream-go](https://github.com/GetStream/stream-go) is a (beta) Go client for [Stream](https://getstream.io/).

[![godoc](https://godoc.org/github.com/GetStream/stream-go?status.svg)](https://godoc.org/github.com/GetStream/stream-go)

[![codecov](https://codecov.io/gh/GetStream/stream-go/branch/master/graph/badge.svg)](https://codecov.io/gh/GetStream/stream-go)

### Full documentation

Documentation for this Go client are available at the [Stream website](https://getstream.io/docs/?language=go).

### Example Usage

Creating a client:
```go
import (
	"fmt"
	"github.com/GetStream/stream-go"
)

// we recommend getting your API credentials using os.Getenv()
client, err := getstream.New(&getstream.Config{
    APIKey:      os.Getenv("STREAM_API_KEY"),
    APISecret:   os.Getenv("STREAM_API_SECRET"),
    AppID:       os.Getenv("STREAM_APP_ID"),
    Location:    os.Getenv("STREAM_REGION"),
})
if err != nil {
    return err
}

// but you can define the variables in code as well, of course
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
    APIKey:      APIKey,
    APISecret:   APISecret,
    AppID:       AppID,
    Location:    Location,
    TimeoutInt:  TimeoutInt,
})

```

Creating a Feed object for a user:

```go
// this code assumes you've created a flat feed named "flat-feed-name" for your app
// and similarly-named feeds for aggregated feeds and notification feeds
// we also recommend using UUID values for users

bobFlatFeed, err := client.FlatFeed("flat-feed-name", "bob-uuid")
if err != nil {
    return err
}

bobAggregatedFeed, err := client.AggregatedFeed("aggregated-feed-name", "bob-uuid")
if err != nil {
    return err
}

bobNotificationFeed, err := client.NotificationFeed("notification-feed-name", "bob-uuid")
if err != nil {
    return err
}
```

Creating an activity on Bob's flat feed:
```go

activity, err := bobFeed.AddActivity(&Activity{
    Verb:      "post",
    ForeignID: "42",
    Object:    "flat:eric",
    Actor:     "flat:john",
})
if err != nil {
    return err
}
```

The library is gradually introducing JWT support. You can generate a client token
for a feed using the following example:

```go
// create a client using your API key and secret
client, err := getstream.New(&getstream.Config{
    APIKey:    os.Getenv("STREAM_API_KEY"),
    APISecret: os.Getenv("STREAM_API_SECRET"),
    AppID:     os.Getenv("STREAM_APP_ID"),
    Location:  os.Getenv("STREAM_REGION"),
})

// create a feed
feed, err := client.FlatFeed("flat-feed-name", "bob-uuid")
if err != nil {
    return err
}

// create a JWT token for the feed
token, err := client.Signer.GenerateFeedScopeToken(
    getstream.ScopeContextFeed,
    getstream.ScopeActionRead,
    bobFeed)
if err != nil {
    fmt.Println(err)
}

// create a new client using the token
// note in the struct below that we're not setting "APISecret"
// but setting "Token" instead:
bobFlatFeedJWTClient, err := getstream.NewWithToken(&getstream.Config{
    APIKey:    os.Getenv("STREAM_API_KEY"),
    Token:     token, // not setting APISecret
    AppID:     os.Getenv("STREAM_APP_ID"),
    Location:  os.Getenv("STREAM_REGION"),
})
if err != nil {
  return err
}
```

JWT support is not yet fully tested on the library, but we'd love to
hear any feedback you have as you try it out.

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
