package getstream

import "regexp"

type ContextMatcher struct {
	re      *regexp.Regexp
	auth    AuthenticationMethod
	context ScopeContext
}

var contextMatchers = []ContextMatcher{
	{
		re:      regexp.MustCompile("^follow_many/$"),
		context: ScopeContextNoContext,
		auth:    ApplicationAuthentication,
	},
	{
		re:      regexp.MustCompile("^feed/add_to_many/$"),
		context: ScopeContextNoContext,
		auth:    ApplicationAuthentication,
	},
	{
		re:      regexp.MustCompile("^activities/$"),
		context: ScopeContextActivities,
		auth:    FeedAuthentication,
	},
	{
		re:      regexp.MustCompile("^feed/.+/.+/followers/$"),
		context: ScopeContextFollower,
		auth:    FeedAuthentication,
	},
	{
		re:      regexp.MustCompile("^feed/.+/.+/follows/$"),
		context: ScopeContextFollower,
		auth:    FeedAuthentication,
	},
	{
		re:      regexp.MustCompile("^feed/.+/.+/following/.+/$"),
		context: ScopeContextFollower,
		auth:    FeedAuthentication,
	},
	{
		re:      regexp.MustCompile("^feed/add_to_many/$"),
		context: ScopeContextNoContext,
		auth:    FeedAuthentication,
	},
	{
		re:      regexp.MustCompile("^follow_many/$"),
		context: ScopeContextNoContext,
		auth:    FeedAuthentication,
	},
	{
		re:      regexp.MustCompile("^feed/"),
		context: ScopeContextFeed,
		auth:    FeedAuthentication,
	},
}
