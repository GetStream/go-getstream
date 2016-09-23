package getstream

import (
	"testing"
	"net/url"
)

func TestClientRequestBadPath(t *testing.T) {
	client, err := New(&Config{
		APIKey:    "my_key",
		APISecret: "my_secret",
		AppID:     "111111",
		Location:  "us-east",
	})
	if err != nil {
		t.Fatal(err)
	}

	_, err = client.request(nil, "get", ":hfi", []byte{}, map[string]string{})
	if err.Error() != "parse :hfi: missing protocol scheme" {
		t.Fatal("Expected error about bad URL path mismatch, got:", err.Error())
	}
}

func TestClientSetStandardParams(t *testing.T) {
	baseURL, _ := url.Parse("https://google.com/")

	client := &Client{
		BaseURL: baseURL,
		Config: &Config{
			APIKey: "my_api_key",
		},
	}
	tmpUrl, err := url.Parse("/")
	if err != nil {
		t.Fatal(err)
	}
	apiUrl := client.BaseURL.ResolveReference(tmpUrl)
	query := apiUrl.Query()
	query = client.setStandardParams(query)

	if _, ok := query["api_key"]; !ok {
		t.Error("API key not set in query params:", query)
	}
	if query["api_key"][0] != "my_api_key" {
		t.Error("API key didn't set as a URL param as expected, got:", query["api_key"][0])
	}

	if _, ok := query["location"]; !ok {
		t.Error("Location not set in query params:", query)
	}
	if query["location"][0] != "unspecified" {
		t.Error("Location didn't set as a URL param as expected, got:", query["location"][0])
	}
}

func TestClientSetRequestParams(t *testing.T) {
	baseURL, _ := url.Parse("https://google.com/")

	client := &Client{
		BaseURL: baseURL,
		Config: &Config{
			APIKey: "my_api_key",
		},
	}
	tmpUrl, err := url.Parse("/")
	if err != nil {
		t.Fatal(err)
	}
	apiUrl := client.BaseURL.ResolveReference(tmpUrl)
	query := apiUrl.Query()

	// passing no params should give us nothing back
	values := client.setRequestParams(query, map[string]string{})
	if len(values) != 0 {
		t.Fatal("Expected empty urlValues, got:", values)
	}

	values = client.setRequestParams(query, map[string]string{
		"foo":"bar",
	})
	if values == nil {
		t.Fatal("Expected urlValues, got nil")
	}

	if _, ok := query["foo"]; !ok {
		t.Error("foo key not set in query params:", query)
	}
	if query["foo"][0] != "bar" {
		t.Error("foo key didn't set as a URL param as expected, got:", query["foo"][0])
	}
}
