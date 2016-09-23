package getstream

// GeneralFeed is a container for Feeds returned from request
// The specific Type will be unknown so no Actions are associated with a GeneralFeed
type GeneralFeed struct {
	Client   *Client
	FeedSlug string
	UserID   string
	token    string
}

// Signature is used to sign Requests : "FeedSlugUserID Token"
func (f *GeneralFeed) Signature() string {
	if f.Token() == "" {
		return f.FeedIDWithoutColon()
	}
	return f.FeedIDWithoutColon() + " " + f.Token()
}

// FeedID is the combo if the FeedSlug and UserID : "FeedSlug:UserID"
func (f *GeneralFeed) FeedID() FeedID {
	return FeedID(f.FeedSlug + ":" + f.UserID)
}

func (f *GeneralFeed) FeedIDWithoutColon() string {
	return f.FeedSlug + f.UserID
}

// SignFeed sets the token on a Feed
func (f *GeneralFeed) SignFeed(signer *Signer) {
	if f.Client.Signer != nil {
		f.token = signer.GenerateToken(f.FeedIDWithoutColon())
	}
}

// Token returns the token of a Feed
func (f *GeneralFeed) Token() string {
	return f.token
}

// GenerateToken returns a new Token for a Feed without setting it to the Feed
func (f *GeneralFeed) GenerateToken(signer *Signer) string {
	if f.Client.Signer != nil {
		return signer.GenerateToken(f.FeedSlug + f.UserID)
	}
	return ""
}

// Unfollow is used to Unfollow a target Feed
func (f *GeneralFeed) Unfollow(client *Client, target *FlatFeed) error {
	f.Client = client
	f.SignFeed(f.Client.Signer)

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "following" + "/" + target.FeedID().Value() + "/"

	return f.Client.del(f, endpoint, nil, nil)
}

// UnfollowAggregated is used to Unfollow a target Aggregated Feed
func (f *GeneralFeed) UnfollowAggregated(client *Client, target *AggregatedFeed) error {
	f.Client = client
	f.SignFeed(f.Client.Signer)

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "following" + "/" + target.FeedID().Value() + "/"

	return f.Client.del(f, endpoint, nil, nil)
}

// UnfollowNotification is used to Unfollow a target Notification Feed
func (f *GeneralFeed) UnfollowNotification(client *Client, target *NotificationFeed) error {
	f.Client = client
	f.SignFeed(f.Client.Signer)

	endpoint := "feed/" + f.FeedSlug + "/" + f.UserID + "/" + "following" + "/" + target.FeedID().Value() + "/"

	return f.Client.del(f, endpoint, nil, nil)
}
