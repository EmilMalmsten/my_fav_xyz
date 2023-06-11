package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/emilmalmsten/my_top_xyz/internal/database"
	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

type TestData struct {
	inputs []createToplistItemRequest
	result bool
}

func TestAreRanksInOrder(t *testing.T) {
	testData := []TestData{
		{
			[]createToplistItemRequest{
				{Rank: 1, Title: "Item 1", Description: "Description 1"},
				{Rank: 2, Title: "Item 2", Description: "Description 2"},
				{Rank: 3, Title: "Item 3", Description: "Description 3"},
			},
			true,
		},
		{
			[]createToplistItemRequest{
				{Rank: 0, Title: "Item 0", Description: "Description 1"},
				{Rank: 1, Title: "Item 2", Description: "Description 2"},
				{Rank: 2, Title: "Item 3", Description: "Description 3"},
			},
			false,
		},
		{
			[]createToplistItemRequest{
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
			RequestBody: createToplistRequest{
				Title:       "test title",
				Description: "test description",
				UserID:      3,
				Items:       []createToplistItemRequest{},
			},
			ExpectedCode: http.StatusCreated,
		},
		{
			Name:          "Invalid request",
			RequestMethod: http.MethodPost,
			RequestBody: createToplistRequest{
				Title:       "",
				Description: "test description",
				UserID:      3,
				Items:       []createToplistItemRequest{},
			},
			ExpectedCode: http.StatusBadRequest,
		},
	}

	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL env var is not set")
	}

	db, err := database.CreateDatabaseConnection(dbURL)
	if err != nil {
		log.Fatalf("unable to initialize database: %v", err)
	}

	apiCfg := apiConfig{
		DB: db,
	}

	r := chi.NewRouter()
	r.Put("/api/toplists", apiCfg.handlerToplistsCreate)

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

		rr := httptest.NewRecorder()
		apiCfg.handlerToplistsCreate(rr, req)

		if rr.Code != tc.ExpectedCode {
			t.Errorf("Expected %d but got %d", tc.ExpectedCode, rr.Code)
		}

	}
}
