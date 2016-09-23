package getstream

import (
	"net/url"
	"time"
)

// Config structure for configuration data
type Config struct {
	APIKey          string
	APISecret       string
	AppID           string
	Location        string
	TimeoutInt      int64
	TimeoutDuration time.Duration
	Version         string
	Token           string
	BaseURL         *url.URL
}

// SetAPIKey sets the API key for your GetStream.io account
// You can find/generate this value on GetStream.io when visiting your application's page
func (c *Config) SetAPIKey(apiKey string) string {
	c.APIKey = apiKey
	return c.APIKey
}

// SetAPISecret sets the API secret for your GetStream.io account
// You can find/generate this value on GetStream.io when visiting your application's page
func (c *Config) SetAPISecret(apiSecret string) string {
	c.APISecret = apiSecret
	return c.APISecret
}

// SetAppID sets the app id (aka site id) for your GetStream.io application
// current site_id/app_id values are integers, but parsing as a string will allow for greater
// flexibility in the future.
func (c *Config) SetAppID(appID string) string {
	c.AppID = appID
	return c.AppID
}

// SetLocation sets the region where your application runs
// Acceptable values are: "us-east"
func (c *Config) SetLocation(location string) string {
	c.Location = location
	return c.Location
}

// SetTimeout sets the timeout, in seconds, for request timers
// timeout param needs to be a time.Duration, ultimately, but we'll convert it here from an int64.
func (c *Config) SetTimeout(timeout int64) time.Duration {
	c.TimeoutInt = timeout
	c.TimeoutDuration = time.Duration(timeout) * time.Second
	return c.TimeoutDuration
}

// SetVersion sets the version of the GetStream API to use
// Currently, the only acceptable value is "v1.0"
func (c *Config) SetVersion(version string) string {
	c.Version = version
	return c.Version
}

// SetBaseURL sets the core url of the GetStream API to use
func (c *Config) SetToken(token string) string {
	c.Token = token
	return c.Token
}

// SetBaseURL sets the core url of the GetStream API to use
func (c *Config) SetBaseURL(baseURL *url.URL) *url.URL {
	c.BaseURL = baseURL
	return c.BaseURL
}
