package database

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func insertToplist(t *testing.T) Toplist {
	user := insertUser(t)
	toplist := Toplist{
		Title:       "My Toplist",
		Description: "This is a mock toplist",
		UserID:      user.ID,
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
		require.Equal(t, insertedToplist.Items[i].ListID, insertedToplist.ToplistID)
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
	toplist2, err := dbTestConfig.GetToplist(toplist1.ToplistID)
	require.NoError(t, err)
	require.NotEmpty(t, toplist2)

	require.Equal(t, toplist1.ToplistID, toplist2.ToplistID)
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

func TestListToplists(t *testing.T) {
	for i := 0; i < 4; i++ {
		_ = insertToplist(t)
	}

	limit := 10
	offset := 1
	toplists, err := dbTestConfig.ListToplists(limit, offset)
	require.NoError(t, err)
	require.LessOrEqual(t, len(toplists), limit)

	for i := range toplists {
		require.NotZero(t, toplists[i].ToplistID)
		require.NotZero(t, toplists[i].Title)
		require.NotZero(t, toplists[i].UserID)
		require.NotZero(t, toplists[i].CreatedAt)
	}
}

func TestListRecentToplists(t *testing.T) {
	for i := 0; i < 4; i++ {
		_ = insertToplist(t)
	}

	limit := 10
	toplistFilterProp := "date"
	toplists, err := dbTestConfig.ListToplistsByProperty(limit, toplistFilterProp)
	require.NoError(t, err)
	require.LessOrEqual(t, len(toplists), limit)

	for i := 1; i < len(toplists); i++ {
		fmt.Println(toplists[i].ToplistID)
		require.NotZero(t, toplists[i].ToplistID)
		require.NotZero(t, toplists[i].Title)
		require.NotZero(t, toplists[i].UserID)
		require.NotZero(t, toplists[i].CreatedAt)
		if toplists[i-1].CreatedAt.Before(toplists[i].CreatedAt) {
			t.Error("Toplists are not sorted correctly")
		}
	}
}

func TestListPopularToplists(t *testing.T) {
	for i := 0; i < 4; i++ {
		_ = insertToplist(t)
	}

	limit := 10
	toplistFilterProp := "views"
	toplists, err := dbTestConfig.ListToplistsByProperty(limit, toplistFilterProp)
	require.NoError(t, err)
	require.LessOrEqual(t, len(toplists), limit)

	for i := 1; i < len(toplists); i++ {
		fmt.Println(toplists[i].ToplistID)
		require.NotZero(t, toplists[i].ToplistID)
		require.NotZero(t, toplists[i].Title)
		require.NotZero(t, toplists[i].UserID)
		require.NotZero(t, toplists[i].CreatedAt)
		require.NotNil(t, toplists[i].Views)
		if toplists[i-1].Views < toplists[i].Views {
			t.Error("Toplists are not sorted correctly")
		}
	}
}

func TestUpdateToplist(t *testing.T) {
	toplist1 := insertToplist(t)

	toplist2 := Toplist{
		ToplistID:   toplist1.ToplistID,
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
	require.Equal(t, toplist1.ToplistID, toplist2.ToplistID)

}

func TestUpdateToplistLonger(t *testing.T) {
	toplist1 := insertToplist(t)

	toplist2 := Toplist{
		ToplistID:   toplist1.ToplistID,
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
		ToplistID:   toplist1.ToplistID,
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

	err := dbTestConfig.DeleteToplist(toplist.ToplistID)
	require.NoError(t, err)

	_, err = dbTestConfig.GetToplist(toplist.ToplistID)
	require.ErrorIs(t, err, ErrNotExist)
}

func TestUpdateToplistViews(t *testing.T) {
	toplist := insertToplist(t)

	updatedToplist, err := dbTestConfig.UpdateToplistViews(toplist.ToplistID)
	require.NoError(t, err)

	require.Equal(t, toplist.Views+1, updatedToplist.Views)
	require.Equal(t, toplist.ToplistID, updatedToplist.ToplistID)
	require.Equal(t, toplist.UserID, updatedToplist.UserID)
}
