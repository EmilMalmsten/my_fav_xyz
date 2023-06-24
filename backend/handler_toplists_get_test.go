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

func TestHandlerToplistsGetOne(t *testing.T) {
	testCases := []struct {
		Name             string
		RequestMethod    string
		RequestToplistID int
		ExpectedCode     int
	}{
		{
			Name:             "Successful get one toplist",
			RequestMethod:    http.MethodGet,
			RequestToplistID: insertedTestToplists[2].ToplistID,
			ExpectedCode:     http.StatusOK,
		},
		{
			Name:             "Fail get invalid toplist id",
			RequestMethod:    http.MethodGet,
			RequestToplistID: -1,
			ExpectedCode:     http.StatusNotFound,
		},
	}

	for _, tc := range testCases {

		endpoint := fmt.Sprintf("/api/toplists/%d", tc.RequestToplistID)
		req, err := http.NewRequest(tc.RequestMethod, endpoint, nil)
		if err != nil {
			t.Fatal(err)
		}

		toplistIDString := strconv.Itoa(tc.RequestToplistID)

		reqContext := chi.NewRouteContext()
		reqContext.URLParams.Add("toplistID", toplistIDString)
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, reqContext))

		rr := httptest.NewRecorder()
		apiCfg.handlerToplistsGetOne(rr, req)

		if rr.Code != tc.ExpectedCode {
			t.Errorf("Expected %d but got %d", tc.ExpectedCode, rr.Code)
		}
	}
}

func TestHandlerToplistsGetMany(t *testing.T) {
	testCases := []struct {
		Name            string
		RequestMethod   string
		RequestPageID   int
		RequestPageSize int
		ExpectedCode    int
	}{
		{
			Name:            "Successful get toplists",
			RequestMethod:   http.MethodGet,
			RequestPageID:   1,
			RequestPageSize: 3,
			ExpectedCode:    http.StatusOK,
		},
		{
			Name:            "Fail invalid page ID",
			RequestMethod:   http.MethodGet,
			RequestPageID:   0,
			RequestPageSize: 3,
			ExpectedCode:    http.StatusBadRequest,
		},
		{
			Name:            "Fail invalid page size",
			RequestMethod:   http.MethodGet,
			RequestPageID:   1,
			RequestPageSize: 99,
			ExpectedCode:    http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		endpoint := fmt.Sprintf("/api/toplists?page_id=%d&page_size=%d", tc.RequestPageID, tc.RequestPageSize)
		req, err := http.NewRequest(tc.RequestMethod, endpoint, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()
		apiCfg.handlerToplistsGetMany(rr, req)

		if rr.Code != tc.ExpectedCode {
			t.Errorf("Expected %d but got %d", tc.ExpectedCode, rr.Code)
		}
	}
}
