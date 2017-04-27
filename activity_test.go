package getstream_test

import (
	"errors"
	"testing"

	getstream "github.com/GetStream/stream-go"
	"github.com/pborman/uuid"
)

func TestActivityMarshallJson(t *testing.T) {
	activity := &getstream.Activity{
		Verb:      "post",
		ForeignID: uuid.New(),
		Object:    "flat:eric",
		Actor:     "flat:john",
	}

	_, err := activity.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}
}

func TestActivityBadForeignKeyMarshall(t *testing.T) {
	activity := &getstream.Activity{
		Verb:      "post",
		ForeignID: "not a real foreign id",
		Object:    "flat:eric",
		Actor:     "flat:john",
	}

	_, err := activity.MarshalJSON()
	if err != nil && err.Error() != "invalid ForeignID" {
		t.Fatal(errors.New("Expected activity.MarshalJSON() to fail on non-UUID ForeignID, it failed because of this:" + err.Error()))
	}
}

func TestActivityUnmarshall(t *testing.T) {
	activity := &getstream.Activity{}
	payload := []byte("{\"actor\":\"flat:john\",\"foreign_id\":\"82d2bb81-069d-427b-9238-8d822012e6d7\",\"object\":\"flat:eric\",\"origin\":\"\",\"time\":\"2016-09-22T21:44:58.821577\",\"verb\":\"post\"}")

	err := activity.UnmarshalJSON(payload)
	if err != nil {
		t.Fatal(err)
	}
}

func TestActivityUnmarshallEmptyPayload(t *testing.T) {
	activity := &getstream.Activity{}

	err := activity.UnmarshalJSON([]byte{})
	if err == nil {
		t.Fatal(err)
	}
	if err.Error() != "unexpected end of JSON input" {
		t.Fatal(errors.New("Expected activity.UnmarshalJSON method to fail on a bad payload, it failed because of this:" + err.Error()))
	}
}

func TestActivityUnmarshallBadPayloadTime(t *testing.T) {
	var err error
	activity := &getstream.Activity{}

	// empty json value for "time" should still set "time" to nil
	payload := []byte("{\"actor\":\"flat:john\",\"foreign_id\":\"82d2bb81-069d-427b-9238-8d822012e6d7\",\"object\":\"flat:eric\",\"origin\":\"\",\"time\":{},\"verb\":\"post\"}")
	err = activity.UnmarshalJSON(payload)
	if err != nil {
		t.Fatal(err)
	}
	if activity.TimeStamp != nil {
		t.Fatal("Expected TimeStamp to be nil if it was empty JSON {}")
	}
	// non-Time value should still parse fine and set "time" to nil
	payload = []byte("{\"actor\":\"flat:john\",\"foreign_id\":\"82d2bb81-069d-427b-9238-8d822012e6d7\",\"object\":\"flat:eric\",\"origin\":\"\",\"time\":\"abc\",\"verb\":\"post\"}")
	err = activity.UnmarshalJSON(payload)
	if err != nil {
		t.Fatal(err)
	}
	if activity.TimeStamp != nil {
		t.Fatal("Expected TimeStamp to be nil if it was an unparseable time")
	}
}

func TestActivityUnmarshallBadPayloadTo(t *testing.T) {
	var err error
	activity := &getstream.Activity{}

	// empty json value for "to" should set "to" to nil
	payload := []byte("{\"to\":null,\"actor\":\"flat:john\",\"foreign_id\":\"82d2bb81-069d-427b-9238-8d822012e6d7\",\"object\":\"flat:eric\",\"origin\":\"\",\"time\":\"2016-09-22T21:44:58.821577\",\"verb\":\"post\"}")
	err = activity.UnmarshalJSON(payload)
	if err != nil {
		t.Fatal(err)
	}
	if activity.To != nil {
		t.Fatal("To JSON was null, expected To to be nil afterward, got:", activity.To)
	}

	// empty json value for "to" should set "to" to nil
	payload = []byte("{\"to\":{},\"actor\":\"flat:john\",\"foreign_id\":\"82d2bb81-069d-427b-9238-8d822012e6d7\",\"object\":\"flat:eric\",\"origin\":\"\",\"time\":\"2016-09-22T21:44:58.821577\",\"verb\":\"post\"}")
	err = activity.UnmarshalJSON(payload)
	if err != nil {
		t.Fatal(err)
	}
	if activity.To != nil {
		t.Fatal("To payload was bad JSON, expected To to be nil afterward, got:", activity.To)
	}

	// two-dimensional To should set To to ... something?
	payload = []byte("{\"to\":[\"bob\"],\"actor\":\"flat:john\",\"foreign_id\":\"82d2bb81-069d-427b-9238-8d822012e6d7\",\"object\":\"flat:eric\",\"origin\":\"\",\"time\":\"2016-09-22T21:44:58.821577\",\"verb\":\"post\"}")
	err = activity.UnmarshalJSON(payload)
	if err != nil {
		t.Fatal(err)
	}
	if len(activity.To) != 0 {
		t.Fatal("To payload was bad JSON, expected To to be nil afterward, got:", activity.To)
	}

	// malformed To userID should null out To
	payload = []byte("{\"to\":[{\"bob\"}],\"actor\":\"flat:john\",\"foreign_id\":\"82d2bb81-069d-427b-9238-8d822012e6d7\",\"object\":\"flat:eric\",\"origin\":\"\",\"time\":\"2016-09-22T21:44:58.821577\",\"verb\":\"post\"}")
	err = activity.UnmarshalJSON(payload)
	if err == nil {
		t.Fatal(err)
	}
	if activity.To != nil {
		t.Fatal("To payload was not a value feedslug:userid format, expected To to be nil afterward, got:", activity.To)
	}
}
