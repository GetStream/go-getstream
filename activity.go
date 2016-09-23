package getstream

import (
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"time"
)

// Activity is a getstream Activity
// Use it to post activities to Feeds
// It is also the response from Fetch and List Requests
type Activity struct {
	ID        string
	Actor     FeedID
	Verb      string
	Object    FeedID
	Target    FeedID
	Origin    FeedID
	TimeStamp *time.Time

	ForeignID string
	Data      *json.RawMessage
	MetaData  map[string]string

	To []Feed
}

// MarshalJSON is the custom marshal function for Activities
// It will be used by json.Marshal()
func (a Activity) MarshalJSON() ([]byte, error) {

	payload := make(map[string]interface{})

	for key, value := range a.MetaData {
		payload[key] = value
	}

	payload["actor"] = a.Actor.Value()
	payload["verb"] = a.Verb
	payload["object"] = a.Object.Value()
	payload["origin"] = a.Origin.Value()

	if a.ID != "" {
		payload["id"] = a.ID
	}
	if a.Target != "" {
		payload["target"] = a.Target.Value()
	}

	if a.Data != nil {
		payload["data"] = a.Data
	}

	if a.ForeignID != "" {
		r, err := regexp.Compile("^[a-z0-9]{8}-[a-z0-9]{4}-[1-5][a-z0-9]{3}-[a-z0-9]{4}-[a-z0-9]{12}$")
		if err != nil {
			return nil, err
		}
		if !r.MatchString(a.ForeignID) {
			return nil, errors.New("invalid ForeignID")
		}
		payload["foreign_id"] = a.ForeignID
	}

	if a.TimeStamp == nil {
		payload["time"] = time.Now().Format("2006-01-02T15:04:05.999999")
	} else {
		payload["time"] = a.TimeStamp.Format("2006-01-02T15:04:05.999999")
	}

	var tos []string
	for _, feed := range a.To {
		to := feed.FeedID().Value()
		if feed.Token() != "" {
			to += " " + feed.Token()
		}
		tos = append(tos, to)
	}

	if len(tos) > 0 {
		payload["to"] = tos
	}

	return json.Marshal(payload)

}

// UnmarshalJSON is the custom unmarshal function for Activities
// It will be used by json.Unmarshal()
func (a *Activity) UnmarshalJSON(b []byte) (err error) {

	rawPayload := make(map[string]*json.RawMessage)
	metadata := make(map[string]string)

	err = json.Unmarshal(b, &rawPayload)
	if err != nil {
		return err
	}

	for key, value := range rawPayload {
		lowerKey := strings.ToLower(key)

		if value == nil {
			continue
		}

		if lowerKey == "id" {
			var strValue string
			json.Unmarshal(*value, &strValue)
			a.ID = strValue
		} else if lowerKey == "actor" {
			var strValue string
			json.Unmarshal(*value, &strValue)
			a.Actor = FeedID(strValue)
		} else if lowerKey == "verb" {
			var strValue string
			json.Unmarshal(*value, &strValue)
			a.Verb = strValue
		} else if lowerKey == "foreign_id" {
			var strValue string
			json.Unmarshal(*value, &strValue)
			a.ForeignID = strValue
		} else if lowerKey == "object" {
			var strValue string
			json.Unmarshal(*value, &strValue)
			a.Object = FeedID(strValue)
		} else if lowerKey == "origin" {
			var strValue string
			json.Unmarshal(*value, &strValue)
			a.Origin = FeedID(strValue)
		} else if lowerKey == "target" {
			var strValue string
			json.Unmarshal(*value, &strValue)
			a.Target = FeedID(strValue)
		} else if lowerKey == "time" {
			var strValue string
			err := json.Unmarshal(*value, &strValue)
			if err != nil {
				continue
			}
			timeStamp, err := time.Parse("2006-01-02T15:04:05.999999", strValue)
			if err != nil {
				continue
			}
			a.TimeStamp = &timeStamp
		} else if lowerKey == "data" {
			a.Data = value
		} else if lowerKey == "to" {

			var to1D []string
			var to2D [][]string

			err := json.Unmarshal(*value, &to1D)
			if err != nil {
				err = nil
				err = json.Unmarshal(*value, &to2D)
				if err != nil {
					continue
				}

				for _, to := range to2D {
					if len(to) == 2 {
						feedStr := to[0] + " " + to[1]
						to1D = append(to1D, feedStr)
					} else if len(to) == 1 {
						to1D = append(to1D, to[0])
					}
				}
			}

			for _, to := range to1D {
				feed := GeneralFeed{}

				match, err := regexp.MatchString(`^\w+:\w+ .*?$`, to)
				if err != nil {
					continue
				}

				if match {
					firstSplit := strings.Split(to, ":")
					secondSplit := strings.Split(firstSplit[1], " ")

					feed.FeedSlug = firstSplit[0]
					feed.UserID = secondSplit[0]
					feed.token = secondSplit[1]
					a.To = append(a.To, &feed)
					continue
				}

				match = false
				err = nil

				match, err = regexp.MatchString(`^\w+:\w+$`, to)
				if err != nil {
					continue
				}

				if match {
					firstSplit := strings.Split(to, ":")

					feed.FeedSlug = firstSplit[0]
					feed.UserID = firstSplit[1]
					a.To = append(a.To, &feed)
					continue
				}
			}
		} else {
			var strValue string
			json.Unmarshal(*value, &strValue)
			metadata[key] = strValue
		}
	}

	a.MetaData = metadata
	return nil

}
