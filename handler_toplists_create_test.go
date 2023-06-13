package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestData struct {
	inputs []toplistItemRequest
	result bool
}

func TestAreRanksInOrder(t *testing.T) {
	testData := []TestData{
		{
			[]toplistItemRequest{
				{Rank: 1, Title: "Item 1", Description: "Description 1"},
				{Rank: 2, Title: "Item 2", Description: "Description 2"},
				{Rank: 3, Title: "Item 3", Description: "Description 3"},
			},
			true,
		},
		{
			[]toplistItemRequest{
				{Rank: 0, Title: "Item 0", Description: "Description 1"},
				{Rank: 1, Title: "Item 2", Description: "Description 2"},
				{Rank: 2, Title: "Item 3", Description: "Description 3"},
			},
			false,
		},
		{
			[]toplistItemRequest{
				{Rank: 0, Title: "Item 1", Description: "Description 1"},
				{Rank: 1, Title: "Item 2", Description: "Description 2"},
				{Rank: 2, Title: "Item 4", Description: "Description 3"},
			},
			false,
		},
	}

	for _, test := range testData {
		result := validateItemRanks(test.inputs)
		if result != test.result {
			t.Errorf("got items %v, expected %t\n", test.inputs, test.result)
		}
	}
}

func TestHandlerToplistsCreate(t *testing.T) {
	testCases := []struct {
		Name          string
		RequestMethod string
		RequestBody   interface{}
		ExpectedCode  int
	}{
		{
			Name:          "Successful creation",
			RequestMethod: http.MethodPost,
			RequestBody: toplistRequest{
				Title:       "test title",
				Description: "test description",
				Items:       []toplistItemRequest{},
			},
			ExpectedCode: http.StatusCreated,
		},
		{
			Name:          "Invalid request",
			RequestMethod: http.MethodPost,
			RequestBody: toplistRequest{
				Title:       "",
				Description: "test description",
				Items:       []toplistItemRequest{},
			},
			ExpectedCode: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {

		body, err := json.Marshal(tc.RequestBody)
		if err != nil {
			t.Fatal(err)
		}
		endpoint := "/api/toplists"
		req, err := http.NewRequest(tc.RequestMethod, endpoint, bytes.NewBuffer(body))
		if err != nil {
			t.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		ctx := context.WithValue(req.Context(), userIDKey, 3)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		apiCfg.handlerToplistsCreate(rr, req)

		if rr.Code != tc.ExpectedCode {
			t.Errorf("Expected %d but got %d", tc.ExpectedCode, rr.Code)
		}
	}
}
