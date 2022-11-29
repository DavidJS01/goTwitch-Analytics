package main

import (
	"errors"
	// "net/http"
	// "net/http/httptest"
	"testing"

	"test.com/m/internal/database"
)

var se database.StreamEvents

type mockStreamEvent struct {
	fakeChannel string
	fakeStatus  bool
}

func TestUpsertResponse(t *testing.T) {
	expectedResponse := Upsert{
		Streamer:    "stream",
		Is_Active:   true,
		Status_Code: 200,
	}
	response := upsertResponse("stream", true, 200)
	if expectedResponse != response {
		t.Errorf("invalid response struct, expected %#v got %#v", expectedResponse, response)
	}
}

func mockInsertStreamer(channel string, is_active bool) error {
	return nil
}

func failingMockInsertStreamer(channel string, is_active bool) error {
	return errors.New("error")
}

// func TestUpsertStreamerHandler(t *testing.T) {
// 	// expectedResponse := `{"streamer":"","is_active":false,"status_code":200}`
// 	req, _ := http.NewRequest("POST", "/stream/upsert", nil)
// 	res := httptest.NewRecorder()
// 	// test for successful 200 code
// 	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { upsertStreamerHandler(w, r, database.InsertStreamer) })
// 	handler.ServeHTTP(res, req)
// 	if status := res.Code; status != http.StatusOK {
// 		t.Errorf("handler returned wrong status code: got %v want %v",
// 			status, http.StatusOK)
// 		// }
// 		// upsertStreamerHandler(res, req, mockInsertStreamer)
// 		// if status := res.Code; status != http.StatusOK {
// 		// 	t.Errorf("Got status code %d, but expected 200", res.Code)
// 		// }
// 		// // if res.Body.String() != `{"streamer":"","is_active":false,"status_code":200}` {
// 		// 	t.Errorf("Expected response %s got %s", expectedResponse, res.Body.String())
// 		// }
// 		// test for 500 error code
// 		// res = httptest.NewRecorder()
// 		// upsertStreamerHandler(res, req, failingMockInsertStreamer)
// 		// if res.Code != 200 {
// 		// 	t.Errorf("Got status code %d, but expected 200", res.Code)
// 		// }
// 	}
// }
