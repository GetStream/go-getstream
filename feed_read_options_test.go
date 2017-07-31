package getstream_test

import (
	"testing"

	getstream "github.com/GetStream/stream-go"
	"github.com/stretchr/testify/assert"
)

func TestNewFeedReadOptions(t *testing.T) {
	opts := getstream.NewFeedReadOptions()
	params := opts.Params()
	assert.Equal(t, params, map[string]string{})
	optsCopy := opts.AddIdGt("1")
	assert.NotEqual(t, opts.Params(), optsCopy.Params())
}

func TestNewFeedReadOptionsLimitOffset(t *testing.T) {
	opts := getstream.NewFeedReadOptions().AddLimit(0).AddOffset(10)
	params := opts.Params()
	assert.Equal(t, params, map[string]string{"limit": "0", "offset": "10"})
}

func TestNewFeedReadOptionsIdGt(t *testing.T) {
	opts := getstream.NewFeedReadOptions().AddIdGt("gt").AddIdGte("gte")
	params := opts.Params()
	assert.Equal(t, params, map[string]string{"id_gt": "gt", "id_gte": "gte"})
}

func TestNewFeedReadOptionsIdLt(t *testing.T) {
	opts := getstream.NewFeedReadOptions().AddIdLt("lt").AddIdLte("lte")
	params := opts.Params()
	assert.Equal(t, params, map[string]string{"id_lt": "lt", "id_lte": "lte"})
}
