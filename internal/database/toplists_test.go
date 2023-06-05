package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func insertToplist(t *testing.T) Toplist {
	toplist := Toplist{
		Title:       "My Toplist",
		Description: "This is a mock toplist",
		UserID:      1,
		Items: []ToplistItem{
			{
				Rank:        1,
				Title:       "Item 1",
				Description: "Description 1",
			},
			{
				Rank:        2,
				Title:       "Item 2",
				Description: "Description 2",
			},
			{
				Rank:        3,
				Title:       "Item 3",
				Description: "Description 3",
			},
		},
	}

	insertedToplist, err := dbTestConfig.InsertToplist(toplist)
	require.NoError(t, err)
	require.NotZero(t, insertedToplist)

	require.Equal(t, toplist.Title, insertedToplist.Title)
	require.Equal(t, toplist.Description, insertedToplist.Description)
	require.Equal(t, toplist.UserID, insertedToplist.UserID)

	require.Equal(t, len(toplist.Items), len(insertedToplist.Items))

	for i := range insertedToplist.Items {
		require.Equal(t, insertedToplist.Items[i].ListID, insertedToplist.ID)
		require.Equal(t, toplist.Items[i].Rank, insertedToplist.Items[i].Rank)
		require.Equal(t, toplist.Items[i].Title, insertedToplist.Items[i].Title)
		require.Equal(t, toplist.Items[i].Description, insertedToplist.Items[i].Description)

	}

	return insertedToplist
}

func TestInsertToplist(t *testing.T) {
	insertToplist(t)
}

func TestGetToplist(t *testing.T) {
	toplist1 := insertToplist(t)
	toplist2, err := dbTestConfig.GetToplist(toplist1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, toplist2)

	require.Equal(t, toplist1.ID, toplist2.ID)
	require.Equal(t, toplist1.Title, toplist2.Title)
	require.Equal(t, toplist1.Description, toplist2.Description)

	require.Equal(t, len(toplist1.Items), len(toplist2.Items))

	for i := range toplist2.Items {
		require.Equal(t, toplist1.Items[i].ListID, toplist2.Items[i].ListID)
		require.Equal(t, toplist1.Items[i].Rank, toplist2.Items[i].Rank)
		require.Equal(t, toplist1.Items[i].Title, toplist2.Items[i].Title)
		require.Equal(t, toplist1.Items[i].Description, toplist2.Items[i].Description)
	}
}

func TestUpdateToplist(t *testing.T) {
	toplist1 := insertToplist(t)

	toplist2 := Toplist{
		ID:          toplist1.ID,
		Title:       toplist1.Title,
		Description: toplist1.Description,
		UserID:      toplist1.UserID,
		Items:       make([]ToplistItem, len(toplist1.Items)),
	}
	copy(toplist2.Items, toplist1.Items)

	toplist2.Title = "Updated My Toplist"

	toplist2, err := dbTestConfig.UpdateToplist(toplist2)
	require.NoError(t, err)
	require.NotEmpty(t, toplist2)

	require.NotEqual(t, toplist1.Title, toplist2.Title)
	require.Equal(t, toplist1.Description, toplist2.Description)
	require.Equal(t, toplist1.UserID, toplist2.UserID)
	require.Equal(t, toplist1.ID, toplist2.ID)

}

func TestUpdateToplistLonger(t *testing.T) {
	toplist1 := insertToplist(t)

	toplist2 := Toplist{
		ID:          toplist1.ID,
		Title:       toplist1.Title,
		Description: toplist1.Description,
		UserID:      toplist1.UserID,
		Items:       make([]ToplistItem, len(toplist1.Items)),
	}
	copy(toplist2.Items, toplist1.Items)

	toplist2.Items = append(toplist2.Items, ToplistItem{
		Rank:        3,
		Title:       "Item 3",
		Description: "Description 3",
	})

	toplist2, err := dbTestConfig.UpdateToplist(toplist2)
	require.NoError(t, err)

	require.Greater(t, len(toplist2.Items), len(toplist1.Items))

	for i := range toplist2.Items {
		if i < len(toplist2.Items)-1 {
			require.Equal(t, toplist1.Items[i].ListID, toplist2.Items[i].ListID)
			require.Equal(t, toplist1.Items[i].Title, toplist2.Items[i].Title)
			require.Equal(t, toplist1.Items[i].Description, toplist2.Items[i].Description)
		} else {
			require.NotEmpty(t, toplist2.Items[i])
		}
	}
}

func TestUpdateToplistShorter(t *testing.T) {
	toplist1 := insertToplist(t)

	toplist2 := Toplist{
		ID:          toplist1.ID,
		Title:       toplist1.Title,
		Description: toplist1.Description,
		UserID:      toplist1.UserID,
		Items:       make([]ToplistItem, len(toplist1.Items)),
	}
	copy(toplist2.Items, toplist1.Items)

	newItems := []ToplistItem{
		{
			Rank:        1,
			Title:       "New Item 1",
			Description: "New Description 1",
		},
		{
			Rank:        2,
			Title:       "New Item 2",
			Description: "New Description 2",
		},
	}

	toplist2.Items = newItems

	toplist2, err := dbTestConfig.UpdateToplist(toplist2)
	require.NoError(t, err)

	require.Less(t, len(toplist2.Items), len(toplist1.Items))

	for i := range toplist1.Items {
		if i < len(toplist1.Items)-1 {
			require.Equal(t, toplist1.Items[i].ListID, toplist2.Items[i].ListID)
			require.NotEqual(t, toplist1.Items[i].Title, toplist2.Items[i].Title)
			require.NotEqual(t, toplist1.Items[i].Description, toplist2.Items[i].Description)
		} else {
			require.NotEmpty(t, toplist1.Items[i])
		}
	}
}

func TestDeleteToplist(t *testing.T) {
	toplist := insertToplist(t)

	err := dbTestConfig.DeleteToplist(toplist.ID)
	require.NoError(t, err)

	_, err = dbTestConfig.GetToplist(toplist.ID)
	require.ErrorIs(t, err, ErrNotExist)
}
