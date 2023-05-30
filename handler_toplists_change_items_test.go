package main

import (
	"testing"

	"github.com/emilmalmsten/my_top_xyz/internal/database"
)

type TestData struct {
	inputs []database.ToplistItem
	result bool
}

func TestAreRanksInOrder(t *testing.T) {
	testData := []TestData{
		{
			[]database.ToplistItem{
				{ID: 1, ListId: 1, Rank: 1, Title: "Item 1", Description: "Description 1"},
				{ID: 2, ListId: 1, Rank: 2, Title: "Item 2", Description: "Description 2"},
				{ID: 3, ListId: 1, Rank: 3, Title: "Item 3", Description: "Description 3"},
			},
			true,
		},
		{
			[]database.ToplistItem{
				{ID: 1, ListId: 1, Rank: 0, Title: "Item 0", Description: "Description 1"},
				{ID: 2, ListId: 1, Rank: 1, Title: "Item 2", Description: "Description 2"},
				{ID: 3, ListId: 1, Rank: 2, Title: "Item 3", Description: "Description 3"},
			},
			false,
		},
		{
			[]database.ToplistItem{
				{ID: 1, ListId: 1, Rank: 0, Title: "Item 1", Description: "Description 1"},
				{ID: 2, ListId: 1, Rank: 1, Title: "Item 2", Description: "Description 2"},
				{ID: 3, ListId: 1, Rank: 2, Title: "Item 4", Description: "Description 3"},
			},
			false,
		},
	}

	for _, test := range testData {
		result := areRanksInOrder(test.inputs)
		if result != test.result {
			t.Errorf("got items %v, expected %t\n", test.inputs, test.result)
		}
	}
}
