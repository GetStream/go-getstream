package getstream

// FeedID is a typealias of string to create some value safety
type FeedID string

// Value returns a String Representation of FeedID
func (f FeedID) Value() string {
	return string(f)
}

// Feed is the interface bundling all Feed Types
// It exposes methods needed for all Types
type Feed interface {
	Signature() string
	FeedID() FeedID
	FeedIDWithoutColon() string
	Token() string
	SignFeed(signer *Signer)
	GenerateToken(signer *Signer) string
}
