package getstream

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Client is used to connect to getstream.io
type Client struct {
	HTTP    *http.Client
	BaseURL *url.URL // https://api.getstream.io/api/
	Config  *Config
	Signer  *Signer
}

/**
 * New returns a GetStream client.
 *
 * Params:
 *   cfg, pointer to a Config structure which takes the API credentials, Location, etc
 * Returns:
 *   Client struct
 */
func New(cfg *Config) (*Client, error) {
	var (
		timeout int64
	)

	if cfg.APIKey == "" {
		return nil, errors.New("Required API Key was not set")
	}

	if cfg.APISecret == "" && cfg.Token == "" {
		return nil, errors.New("API Secret or Token was not set, one or the other is required")
	}

	if cfg.TimeoutInt <= 0 {
		timeout = 3
	} else {
		timeout = cfg.TimeoutInt
	}
	cfg.SetTimeout(timeout)

	if cfg.Version == "" {
		cfg.Version = "v1.0"
	}

	location := "api"
	port := ""
	secure := "s"
	mainDomain := ".getstream.io"

	if cfg.Location == "localhost" {
		location = "localhost"
		port = ":8000"
		secure = ""
		mainDomain = ""
	} else if cfg.Location != "" {
		location = cfg.Location + "-api"
		if cfg.Location == "qa" {
			secure = ""
		}
	}

	baseURL, err := url.Parse("http" + secure + "://" + location + mainDomain + port + "/api/" + cfg.Version + "/")
	if err != nil {
		return nil, err
	}
	cfg.SetBaseURL(baseURL)

	var secret string
	if cfg.Token != "" {
		// build the Signature mechanism based on a Token value passed to the client setup
		cfg.SetAPISecret("")
		secret = cfg.Token
	} else {
		// build the Signature based on the API Secret
		cfg.SetToken("")
		secret = cfg.APISecret
	}

	signer := &Signer{
		Key:    cfg.APIKey,
		Secret: secret,
	}

	client := &Client{
		HTTP: &http.Client{
			Transport: GETSTREAM_TRANSPORT,
			Timeout:   cfg.TimeoutDuration,
		},
		BaseURL: baseURL,
		Config:  cfg,
		Signer:  signer,
	}

	return client, nil
}

// FlatFeed returns a getstream feed
// Slug is the FlatFeed Group name
// id is the Specific FlatFeed inside a FlatFeed Group
// to get the feed for Bob you would pass something like "user" as slug and "bob" as the id
func (c *Client) FlatFeed(feedSlug string, userID string) (*FlatFeed, error) {
	var err error

	feedSlug, err = ValidateFeedSlug(feedSlug)
	if err != nil {
		return nil, err
	}
	userID, err = ValidateUserID(userID)
	if err != nil {
		return nil, err
	}

	feed := &FlatFeed{
		baseFeed{
			Client:   c,
			FeedSlug: feedSlug,
			UserID:   userID,
		},
	}

	feed.SignFeed(c.Signer)
	return feed, nil
}

// NotificationFeed returns a getstream feed
// Slug is the NotificationFeed Group name
// id is the Specific NotificationFeed inside a NotificationFeed Group
// to get the feed for Bob you would pass something like "user" as slug and "bob" as the id
func (c *Client) NotificationFeed(feedSlug string, userID string) (*NotificationFeed, error) {
	var err error

	feedSlug, err = ValidateFeedSlug(feedSlug)
	if err != nil {
		return nil, err
	}
	userID, err = ValidateUserID(userID)
	if err != nil {
		return nil, err
	}

	feed := &NotificationFeed{
		baseFeed{
			Client:   c,
			FeedSlug: feedSlug,
			UserID:   userID,
		},
	}

	feed.SignFeed(c.Signer)
	return feed, nil
}

// AggregatedFeed returns a getstream feed
// Slug is the AggregatedFeed Group name
// id is the Specific AggregatedFeed inside a AggregatedFeed Group
// to get the feed for Bob you would pass something like "user" as slug and "bob" as the id
func (c *Client) AggregatedFeed(feedSlug string, userID string) (*AggregatedFeed, error) {
	var err error

	feedSlug, err = ValidateFeedSlug(feedSlug)
	if err != nil {
		return nil, err
	}
	userID, err = ValidateUserID(userID)
	if err != nil {
		return nil, err
	}

	feed := &AggregatedFeed{
		baseFeed{
			Client:   c,
			FeedSlug: feedSlug,
			UserID:   userID,
		},
	}

	feed.SignFeed(c.Signer)
	return feed, nil
}

// absoluteUrl create a url.URL instance and sets query params (bad!!!)
func (c *Client) AbsoluteURL(path string) (*url.URL, error) {
	result, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	result = c.BaseURL.ResolveReference(result)

	qs := result.Query()
	qs.Set("api_key", c.Config.APIKey)
	if c.Config.Location == "" || c.Config.Location == "localhost" {
		qs.Set("location", "unspecified")
	} else {
		qs.Set("location", c.Config.Location)
	}
	result.RawQuery = qs.Encode()

	return result, nil
}

// get request helper
func (c *Client) get(f Feed, path string, payload []byte, params map[string]string) ([]byte, error) {
	// we force an empty body payload because GET requests cannot have a body with our API
	return c.request(f, "GET", path, []byte{}, params)
}

// post request helper
func (c *Client) post(f Feed, path string, payload []byte, params map[string]string) ([]byte, error) {
	return c.request(f, "POST", path, payload, params)
}

// delete request helper
func (c *Client) del(f Feed, path string, payload []byte, params map[string]string) error {
	_, err := c.request(f, "DELETE", path, payload, params)
	return err
}

func (c *Client) getMatcher(path string) (ContextMatcher, error) {
	for _, matcher := range contextMatchers {
		if matcher.re.MatchString(path) {
			return matcher, nil
		}
	}
	return ContextMatcher{}, fmt.Errorf("invalid request path")
}

// request helper
func (c *Client) request(feed Feed, method string, path string, payload []byte, params map[string]string) ([]byte, error) {
	var requestBody io.Reader

	apiUrl, err := url.Parse(path)
	if err != nil {
		return nil, err
	}

	apiUrl = c.BaseURL.ResolveReference(apiUrl)

	query := apiUrl.Query()
	query = c.setStandardParams(query)
	query = c.setRequestParams(query, params)
	apiUrl.RawQuery = query.Encode()

	if payload == nil {
		requestBody = nil
	} else {
		requestBody = bytes.NewBuffer(payload)
	}

	req, err := http.NewRequest(method, apiUrl.String(), requestBody)
	if err != nil {
		return nil, err
	}

	c.setBaseHeaders(req)

	matcher, err := c.getMatcher(path)
	if err != nil {
		return nil, err
	}
	c.setAuthSigAndHeaders(req, feed, path, matcher)

	// perform the http request
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	// handle the response
	switch {
	case resp.StatusCode/100 == 2: // SUCCESS
		return body, nil
	default:
		var respErr Error
		err = json.Unmarshal(body, &respErr)
		if err != nil {
			return nil, err
		}
		return nil, &respErr
	}
}

func (c *Client) setStandardParams(query url.Values) url.Values {
	query.Set("api_key", c.Config.APIKey)
	if c.Config.Location == "" || c.Config.Location == "localhost" {
		query.Set("location", "unspecified")
	} else {
		query.Set("location", c.Config.Location)
	}

	return query
}

func (c *Client) setRequestParams(query url.Values, params map[string]string) url.Values {
	for key, value := range params {
		query.Set(key, value)
	}
	return query
}

/* setBaseHeaders - set common headers for every request
 * params:
 *    request, pointer to http.Request
 */
func (c *Client) setBaseHeaders(request *http.Request) {
	request.Header.Set("X-Stream-Client", "stream-go-client-"+VERSION)
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("X-Api-Key", c.Config.APIKey)

	t := time.Now()
	request.Header.Set("Date", t.Format("Mon, 2 Jan 2006 15:04:05 MST"))
}

func (c *Client) setAuthSigAndHeaders(request *http.Request, feed Feed, path string, matcher ContextMatcher) error {
	authenticator, ok := authenticators[matcher.auth]
	if !ok {
		return fmt.Errorf("missing authentication method")
	}
	return authenticator.Authenticate(c.Signer, request, matcher.context, feed)
}

type PostFlatFeedFollowingManyInput struct {
	Source string `json:"source"`
	Target string `json:"target"`
}

/** PrepFollowFlatFeed - prepares JSON needed for one feed to follow another

Params:
targetFeed, FlatFeed which wants to follow another
sourceFeed, FlatFeed which is to be followed

Returns:
[]byte, array of bytes of JSON suitable for API consumption
*/
func (c *Client) PrepFollowFlatFeed(targetFeed *FlatFeed, sourceFeed *FlatFeed) *PostFlatFeedFollowingManyInput {
	return &PostFlatFeedFollowingManyInput{
		Source: sourceFeed.FeedSlug + ":" + sourceFeed.UserID,
		Target: targetFeed.FeedSlug + ":" + targetFeed.UserID,
	}
}
func (c *Client) PrepFollowAggregatedFeed(targetFeed *FlatFeed, sourceFeed *AggregatedFeed) *PostFlatFeedFollowingManyInput {
	return &PostFlatFeedFollowingManyInput{
		Source: sourceFeed.FeedSlug + ":" + sourceFeed.UserID,
		Target: targetFeed.FeedSlug + ":" + targetFeed.UserID,
	}
}
func (c *Client) PrepFollowNotificationFeed(targetFeed *FlatFeed, sourceFeed *NotificationFeed) *PostFlatFeedFollowingManyInput {
	return &PostFlatFeedFollowingManyInput{
		Source: sourceFeed.FeedSlug + ":" + sourceFeed.UserID,
		Target: targetFeed.FeedSlug + ":" + targetFeed.UserID,
	}
}

type PostActivityToManyInput struct {
	Activity Activity `json:"activity"`
	FeedIDs  []string `json:"feeds"`
}

func (c *Client) AddActivityToMany(activity Activity, feeds []string) error {
	payload := &PostActivityToManyInput{
		Activity: activity,
		FeedIDs:  feeds,
	}

	final_payload, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	endpoint := "feed/add_to_many/"
	params := map[string]string{}
	_, err = c.post(nil, endpoint, final_payload, params)
	return err
}
