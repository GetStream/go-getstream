package getstream_test

import (
	"os"

	getstream "github.com/GetStream/stream-go"
	"testing"
)

func PreTestSetup(t *testing.T) *getstream.Client {
	client, err := doTestSetup(&getstream.Config{
		APIKey:     os.Getenv("STREAM_API_KEY"),
		APISecret:  os.Getenv("STREAM_API_SECRET"),
		AppID:      os.Getenv("STREAM_APP_ID"),
		Location:   os.Getenv("STREAM_REGION"),
		TimeoutInt: 1000,
	})
	if err != nil {
		t.FailNow()
	}
	return client
}

func doTestSetup(cfg *getstream.Config) (*getstream.Client, error) {
	return getstream.New(cfg)
}
