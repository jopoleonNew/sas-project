package yandex

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"encoding/json"

	"reflect"

	"io/ioutil"

	api "github.com/nk2ge5k/goyad"
	"github.com/nk2ge5k/goyad/campaigns"
)

var testAccount []byte

func TestCollectCampaings(t *testing.T) {
	var e error
	testAccount, e = ioutil.ReadFile("testAccount.json")
	if e != nil {
		t.Fatalf("cant read file with test json, error: %v", e)
	}
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Requestid", "12344")
		fmt.Fprintln(w, string(testAccount))
	}))
	defer ts.Close()
	ci := api.NewClient()
	ci.Login = "testlogin"
	ci.Token = api.Token{Value: "123"}
	ci.ApiUrl = ts.URL + "/"
	resCamps, err := collectCampaings(ci)
	if err != nil {
		t.Fatalf("collectCampaings(client) error: %v", err)
	}

	expceted := map[string]campaigns.GetResponse{}
	if err := json.Unmarshal(testAccount, &expceted); err != nil {
		t.Fatalf("bad response from collectCampaings, can't unmarshal: %v", err)
	}
	if expectedRes, ok := expceted["result"]; ok {

		var res, expected []byte
		if res, err = json.Marshal(resCamps.Campaigns); err != nil {
			t.Fatalf("cant marshal : %v", err)
		}
		if expected, err = json.Marshal(expectedRes.Campaigns); err != nil {
			t.Fatalf("cant marshal : %v", err)
		}
		if !reflect.DeepEqual(res, expected) {
			t.Fatalf("result is not equal to expected value")
		}
	} else {
		t.Fatalf("can't find expceted[result] value in result map: ", expectedRes)
	}

}

//func TestHealthCheckHandler(t *testing.T) {
//	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
//	// pass 'nil' as the third parameter.
//	req, err := http.NewRequest("GET", "/health-check", nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
//	rr := httptest.NewRecorder()
//	handler := http.HandlerFunc(HealthCheckHandler)
//
//	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
//	// directly and pass in our Request and ResponseRecorder.
//	handler.ServeHTTP(rr, req)
//
//	// Check the status code is what we expect.
//	if status := rr.Code; status != http.StatusOK {
//		t.Errorf("handler returned wrong status code: got %v want %v",
//			status, http.StatusOK)
//	}
//
//	// Check the response body is what we expect.
//	expected := `{"alive": true}`
//	if rr.Body.String() != expected {
//		t.Errorf("handler returned unexpected body: got %v want %v",
//			rr.Body.String(), expected)
//	}
//}
