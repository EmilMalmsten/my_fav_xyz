package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/go-chi/chi"
)

func TestHandlerToplistsDelete(t *testing.T) {
	testCases := []struct {
		Name             string
		RequestMethod    string
		RequestToplistID int
		RequestUserID    int
		ExpectedCode     int
	}{
		{
			Name:             "Successful deletion",
			RequestMethod:    http.MethodDelete,
			RequestToplistID: insertedTestToplists[0].ToplistID,
			RequestUserID:    insertedTestUser.ID,
			ExpectedCode:     http.StatusOK,
		},
		{
			Name:             "Unauthorized deletion",
			RequestMethod:    http.MethodDelete,
			RequestToplistID: insertedTestToplists[1].ToplistID,
			RequestUserID:    15,
			ExpectedCode:     http.StatusUnauthorized,
		},
		{
			Name:             "Delete non existent toplist",
			RequestMethod:    http.MethodDelete,
			RequestToplistID: insertedTestToplists[0].ToplistID,
			RequestUserID:    insertedTestUser.ID,
			ExpectedCode:     http.StatusNotFound,
		},
	}

	for _, tc := range testCases {

		endpoint := fmt.Sprintf("/api/toplists/%d", tc.RequestToplistID)
		fmt.Println(tc)
		req, err := http.NewRequest(tc.RequestMethod, endpoint, nil)
		if err != nil {
			t.Fatal(err)
		}

		toplistIDString := strconv.Itoa(tc.RequestToplistID)

		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("toplistID", toplistIDString)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))
		req = req.WithContext(context.WithValue(req.Context(), userIDKey, tc.RequestUserID))

		rr := httptest.NewRecorder()
		apiCfg.handlerToplistsDelete(rr, req)

		if rr.Code != tc.ExpectedCode {
			t.Errorf("Expected %d but got %d", tc.ExpectedCode, rr.Code)
		}
	}
}
