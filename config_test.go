package getstream_test

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	getstream "github.com/GetStream/stream-go"
)

func TestWerckerConfig(t *testing.T) {

	testAPIKey := os.Getenv("key")
	testAPISecret := os.Getenv("secret")
	testAppID := os.Getenv("app_id")
	testRegion := os.Getenv("region")

	if testAPIKey == "" || testAPISecret == "" || testAppID == "" || testRegion == "" {
		t.Fatal()
	}
}

func TestConfig_SetAPIKey(t *testing.T) {
	cfg := getstream.Config{APIKey: "123"}
	if cfg.APIKey != "123" {
		t.Error(fmt.Sprintf("building cfg didn't set APIKey to '234', got %s", cfg.APIKey))
	}

	chk := cfg.SetAPIKey("234")
	if chk != "234" {
		t.Error(fmt.Sprintf("SetAPIKey didn't return '234', got %s", chk))
	}
	if cfg.APIKey != "234" {
		t.Error(fmt.Sprintf("cfg.APIKey isn't 234, got %s", cfg.APIKey))
	}
}

func TestConfig_SetAPISecret(t *testing.T) {
	cfg := getstream.Config{APISecret: "123"}
	if cfg.APISecret != "123" {
		t.Error(fmt.Sprintf("building cfg didn't set APISecret to '234', got %s", cfg.APISecret))
	}

	chk := cfg.SetAPISecret("234")
	if chk != "234" {
		t.Error(fmt.Sprintf("SetAPISecret didn't return '234', got %s", chk))
	}
	if cfg.APISecret != "234" {
		t.Error(fmt.Sprintf("cfg.APISecret isn't 234, got %s", cfg.APISecret))
	}
}

func TestConfig_SetAppID(t *testing.T) {
	cfg := getstream.Config{AppID: "123"}
	if cfg.AppID != "123" {
		t.Error(fmt.Sprintf("building cfg didn't set AppID to '234', got %s", cfg.AppID))
	}

	chk := cfg.SetAppID("234")
	if chk != "234" {
		t.Error(fmt.Sprintf("SetAppID didn't return '234', got %s", chk))
	}
	if cfg.AppID != "234" {
		t.Error(fmt.Sprintf("cfg.AppID isn't 234, got %s", cfg.AppID))
	}
}

func TestConfig_SetLocation(t *testing.T) {
	cfg := getstream.Config{Location: "123"}
	if cfg.Location != "123" {
		t.Error(errors.New(fmt.Sprintf("building cfg didn't set Location to '234', got %s", cfg.Location)))
	}

	chk := cfg.SetLocation("234")
	if chk != "234" {
		t.Error(errors.New(fmt.Sprintf("SetLocation didn't return '234', got %s", chk)))
	}
	if cfg.Location != "234" {
		t.Error(errors.New(fmt.Sprintf("cfg.Location isn't 234, got %s", cfg.Location)))
	}
}

func TestConfig_SetTimeout(t *testing.T) {
	cfg := getstream.Config{TimeoutInt: 123, TimeoutDuration: time.Duration(456)}
	if cfg.TimeoutInt != 123 {
		t.Error(fmt.Sprintf("building cfg didn't set TimeoutInt to '123', got %d", cfg.TimeoutInt))
	}
	if cfg.TimeoutDuration != 456 {
		t.Error(fmt.Sprintf("building cfg didn't set timeoutDuration to time.Duration(456), got %d", cfg.TimeoutDuration))
	}

	chk := cfg.SetTimeout(234)
	if chk != time.Duration(234*time.Second) {
		t.Error(fmt.Sprintf("SetTimeoutInt didn't return time.Duration(234), got %d", chk))
	}
	if cfg.TimeoutInt != 234 {
		t.Error(fmt.Sprintf("cfg.TimeoutInt isn't int(234), got %s", cfg.TimeoutInt))
	}
	if cfg.TimeoutDuration != time.Duration(234*time.Second) {
		t.Error(fmt.Sprintf("setting cfg.TimeoutInt didn't set timeoutDuraction to time.Duraction(234), got %d", cfg.TimeoutInt))
	}
}

func TestConfig_SetVersion(t *testing.T) {
	cfg := getstream.Config{Version: "123"}
	if cfg.Version != "123" {
		t.Error(fmt.Sprintf("building cfg didn't set Version to '234', got %s", cfg.Version))
	}

	chk := cfg.SetVersion("234")
	if chk != "234" {
		t.Error(fmt.Sprintf("SetVersion didn't return '234', got %s", chk))
	}
	if cfg.Version != "234" {
		t.Error(fmt.Sprintf("cfg.Version isn't 234, got %s", cfg.Version))
	}
}

func TestConfig_SetToken(t *testing.T) {
	cfg := getstream.Config{Token: "123"}
	if cfg.Token != "123" {
		t.Error(fmt.Sprintf("building cfg didn't set Token to '234', got %s", cfg.Token))
	}

	chk := cfg.SetToken("234")
	if chk != "234" {
		t.Error(fmt.Sprintf("SetToken didn't return '234', got %s", chk))
	}
	if cfg.Token != "234" {
		t.Error(fmt.Sprintf("cfg.Token isn't '234', got '%s'", cfg.Token))
	}
}

func TestConfig_SetBaseURL(t *testing.T) {
	url, _ := url.Parse("http://api.getstream.io")
	cfg := getstream.Config{BaseURL: url}
	if cfg.BaseURL != url {
		t.Error(fmt.Sprintf("building cfg didn't set BaseURL to 'http://use-east-api.getstream.io', got %s", cfg.BaseURL))
	}

	url, _ = url.Parse("http://use-east-api.getstream.io")
	chk := cfg.SetBaseURL(url)
	if chk != url {
		t.Error(fmt.Sprintf("SetBaseURL didn't return 'http://use-east-api.getstream.io', got %s", chk))
	}
	if cfg.BaseURL != url {
		t.Error(fmt.Sprintf("cfg.BaseURL isn't http://use-east-api.getstream.io, got %s", cfg.BaseURL))
	}
}
