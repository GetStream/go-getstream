package getstream_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	getstream "github.com/GetStream/stream-go"
)

func TestActivityMarshallJson(t *testing.T) {
	activity := &getstream.Activity{
		Verb:      "post",
		ForeignID: RandString(8),
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
	payload := []byte("{\"actor\":\"flat:john\",\"foreign_id\":\"82d2bb81-069d-427b-9238-8d822012e6d7\",\"object\":\"flat:eric\",\"origin\":\"\",\"time\":\"2016-09-22T21:44:58.821577\",\"verb\":\"post\"}")

	// empty json value for "time" should still set "time" to nil
	payload = []byte("{\"actor\":\"flat:john\",\"foreign_id\":\"82d2bb81-069d-427b-9238-8d822012e6d7\",\"object\":\"flat:eric\",\"origin\":\"\",\"time\":{},\"verb\":\"post\"}")
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
	payload := []byte("{\"actor\":\"flat:john\",\"foreign_id\":\"82d2bb81-069d-427b-9238-8d822012e6d7\",\"object\":\"flat:eric\",\"origin\":\"\",\"time\":\"2016-09-22T21:44:58.821577\",\"verb\":\"post\"}")
	// empty json value for "to" should set "to" to nil
	payload = []byte("{\"to\":null,\"actor\":\"flat:john\",\"foreign_id\":\"82d2bb81-069d-427b-9238-8d822012e6d7\",\"object\":\"flat:eric\",\"origin\":\"\",\"time\":\"2016-09-22T21:44:58.821577\",\"verb\":\"post\"}")
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

func TestActivityUnmarshalScore(t *testing.T) {
	activity := &getstream.Activity{}

	expected := 1.123
	payload := []byte(fmt.Sprintf(`{ "score": %f }`, expected))

	err := activity.UnmarshalJSON(payload)
	if err != nil {
		t.Fatal(err)
	}

	if *activity.Score != expected {
		t.Fatalf("expected %f, got %f", expected, *activity.Score)
	}
}

func TestActivityMarshalScore(t *testing.T) {
	expectedScore := 1.123
	activity := &getstream.Activity{
		Score: &expectedScore,
	}

	out, err := activity.MarshalJSON()
	if err != nil {
		t.Fatal(err)
	}

	var dest map[string]interface{}
	if err := json.Unmarshal(out, &dest); err != nil {
		t.Fatal(err)
	}

	if score, ok := dest["score"].(float64); !ok || score != expectedScore {
		t.Fatalf("expected score %f, got %f", expectedScore, score)
	}
}
