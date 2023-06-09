package main

import (
	"testing"
)

type TestData struct {
	inputs []createToplistItemRequest
	result bool
}

func TestAreRanksInOrder(t *testing.T) {
	testData := []TestData{
		{
			[]createToplistItemRequest{
				{ListID: 1, Rank: 1, Title: "Item 1", Description: "Description 1"},
				{ListID: 1, Rank: 2, Title: "Item 2", Description: "Description 2"},
				{ListID: 1, Rank: 3, Title: "Item 3", Description: "Description 3"},
			},
			true,
		},
		{
			[]createToplistItemRequest{
				{ListID: 1, Rank: 0, Title: "Item 0", Description: "Description 1"},
				{ListID: 1, Rank: 1, Title: "Item 2", Description: "Description 2"},
				{ListID: 1, Rank: 2, Title: "Item 3", Description: "Description 3"},
			},
			false,
		},
		{
			[]createToplistItemRequest{
				{ListID: 1, Rank: 0, Title: "Item 1", Description: "Description 1"},
				{ListID: 1, Rank: 1, Title: "Item 2", Description: "Description 2"},
				{ListID: 1, Rank: 2, Title: "Item 4", Description: "Description 3"},
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
