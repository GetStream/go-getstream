package getstream_test

import (
	"os"

	getstream "github.com/GetStream/stream-go"
)

func PreTestSetup() (*getstream.Client, error) {
	return doTestSetup(&getstream.Config{
		APIKey:     os.Getenv("STREAM_API_KEY"),
		APISecret:  os.Getenv("STREAM_API_SECRET"),
		AppID:      os.Getenv("STREAM_APP_ID"),
		Location:   os.Getenv("STREAM_REGION"),
		TimeoutInt: 1000,
	})
}

func doTestSetup(cfg *getstream.Config) (*getstream.Client, error) {
	return getstream.New(cfg)
}
