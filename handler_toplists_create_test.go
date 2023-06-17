package main

import (
	"bytes"
	"context"
	"encoding/json"
	_ "errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestValidateToplistValues(t *testing.T) {
	type TestData struct {
		inputs     toplistRequest
		exptectErr bool
	}
	testData := []TestData{
		{
			toplistRequest{
				Title:       "Test toplist",
				Description: "Test description",
				Items:       []toplistItemRequest{},
			},
			false,
		},
		{
			toplistRequest{
				Title:       "",
				Description: "Test description",
				Items:       []toplistItemRequest{},
			},
			true,
		},
		{
			toplistRequest{
				Title:       "Test toplist",
				Description: "Test description",
			},
			false,
		},
	}

	for _, test := range testData {
		var gotErr bool
		err := validateToplistValues(test.inputs)
		if err != nil {
			gotErr = true
		} else {
			gotErr = false
		}
		if gotErr != test.exptectErr {
			t.Errorf("got items %v, expected err to be %t but got %t\n", test.inputs, test.exptectErr, gotErr)
		}
	}
}

func TestValidateListItemValues(t *testing.T) {
	type TestData struct {
		inputs     []toplistItemRequest
		exptectErr bool
	}

	testData := []TestData{
		{
			[]toplistItemRequest{
				{Rank: 1, Title: "Item 1", Description: "Description 1"},
				{Rank: 2, Title: "Item 2", Description: "Description 2"},
				{Rank: 3, Title: "Item 3", Description: "Description 3"},
			},
			false,
		},
		{
			[]toplistItemRequest{
				{Rank: 0, Title: "Item 0", Description: "Description 1"},
				{Rank: 1, Title: "Item 2", Description: "Description 2"},
				{Rank: 2, Title: "Item 3", Description: "Description 3"},
			},
			true,
		},
		{
			[]toplistItemRequest{
				{Rank: 0, Title: "Item 1", Description: "Description 1"},
				{Rank: 1, Title: "Item 2", Description: "Description 2"},
				{Rank: 2, Title: "Item 4", Description: "Description 3"},
			},
			true,
		},
		{
			[]toplistItemRequest{},
			false,
		},
	}

	for _, test := range testData {
		var gotErr bool
		err := validateListItemValues(test.inputs)
		if err != nil {
			gotErr = true
		} else {
			gotErr = false
		}
		if gotErr != test.exptectErr {
			t.Errorf("got items %v, expected %t but got %t\n", test.inputs, test.exptectErr, gotErr)
		}
	}
}

func TestHandlerToplistsCreate(t *testing.T) {
	testCases := []struct {
		Name          string
		RequestMethod string
		RequestBody   interface{}
		RequestUserID interface{}
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
			RequestUserID: 1,
			ExpectedCode:  http.StatusCreated,
		},
		{
			Name:          "Invalid request",
			RequestMethod: http.MethodPost,
			RequestBody: toplistRequest{
				Title:       "",
				Description: "test description",
				Items:       []toplistItemRequest{},
			},
			RequestUserID: 1,
			ExpectedCode:  http.StatusBadRequest,
		},
		{
			Name:          "Invalid userID",
			RequestMethod: http.MethodPost,
			RequestBody: toplistRequest{
				Title:       "test title",
				Description: "test description",
				Items:       []toplistItemRequest{},
			},
			RequestUserID: nil,
			ExpectedCode:  http.StatusBadRequest,
		},
		{
			Name:          "Invalid request body",
			RequestMethod: http.MethodPost,
			RequestBody:   "test",
			RequestUserID: 1,
			ExpectedCode:  http.StatusInternalServerError,
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
		ctx := context.WithValue(req.Context(), userIDKey, tc.RequestUserID)
		req = req.WithContext(ctx)

		rr := httptest.NewRecorder()
		apiCfg.handlerToplistsCreate(rr, req)

		if rr.Code != tc.ExpectedCode {
			t.Errorf("Expected %d but got %d", tc.ExpectedCode, rr.Code)
		}
	}
}
