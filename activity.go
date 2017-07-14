package getstream

import (
	"encoding/json"
	"strings"
	"time"
)

// Activity is a getstream Activity
// Use it to post activities to Feeds
// It is also the response from Fetch and List Requests
type Activity struct {
	ID        string
	Actor     string
	Verb      string
	Object    string
	Target    string
	Origin    FeedID
	TimeStamp *time.Time

	ForeignID string
	Data      *json.RawMessage
	MetaData  map[string]string

	To       []FeedID
	signedTo []string
}

// MarshalJSON is the custom marshal function for Activities
// It will be used by json.Marshal()
func (a Activity) MarshalJSON() ([]byte, error) {

	payload := make(map[string]interface{})

	for key, value := range a.MetaData {
		payload[key] = value
	}

	payload["actor"] = a.Actor
	payload["verb"] = a.Verb
	payload["object"] = a.Object
	payload["origin"] = a.Origin.Value()

	if a.ID != "" {
		payload["id"] = a.ID
	}
	if a.Target != "" {
		payload["target"] = a.Target
	}

	if a.Data != nil {
		payload["data"] = a.Data
	}

	if a.ForeignID != "" {
		payload["foreign_id"] = a.ForeignID
	}

	if a.TimeStamp == nil {
		payload["time"] = time.Now().Format("2006-01-02T15:04:05.999999")
	} else {
		payload["time"] = a.TimeStamp.Format("2006-01-02T15:04:05.999999")
	}

	if len(a.signedTo) > 0 {
		payload["to"] = a.signedTo
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
			a.Actor = strValue
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
			a.Object = strValue
		} else if lowerKey == "origin" {
			var strValue string
			json.Unmarshal(*value, &strValue)
			a.Origin = FeedID(strValue)
		} else if lowerKey == "target" {
			var strValue string
			json.Unmarshal(*value, &strValue)
			a.Target = strValue
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
			var to []string
			if err := json.Unmarshal(*value, &to); err == nil {
				for _, to := range to {
					a.To = append(a.To, FeedID(to))
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
