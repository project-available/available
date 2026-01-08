package db

import (
	"context"
	"testing"

	"github.com/project-available/available/utils"
	"github.com/stretchr/testify/require"
)

func createRandomCustomField(t *testing.T) CustomField {
	key := utils.RandomString(10)
	customField, err := testQuery.CreateCustomField(context.Background(), key)
	require.NoError(t, err)
	require.NotEmpty(t, customField)
	require.Equal(t, key, customField.Key)
	require.True(t, customField.Shown)
	require.NotZero(t, customField.ID)
	return customField
}

func TestCreateCustomField(t *testing.T) {
	createRandomCustomField(t)
}

func TestListCustomFields(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomCustomField(t)
	}
	customFields, err := testQuery.ListCustomFields(context.Background())
	require.NoError(t, err)
	for _, cf := range customFields {
		require.NotEmpty(t, cf)
	}
}

func TestUpdateCustomFieldShown(t *testing.T) {
	cf1 := createRandomCustomField(t)
	newName := utils.RandomString(10)
	arg := UpdateCustomFieldShownParams{
		ID:    cf1.ID,
		Key:   newName,
		Shown: false,
	}
	err := testQuery.UpdateCustomFieldShown(context.Background(), arg)
	require.NoError(t, err)
	// Verification needs a Get but we don't have GetCustomField in the provided sql.go,
	// relying on List or just standard error check which implies success in unit test context for generated code usually.
	// But let's see if we can use List to find it.
	cfs, err := testQuery.ListCustomFields(context.Background())
	require.NoError(t, err)
	var found bool
	for _, cf := range cfs {
		if cf.ID == cf1.ID {
			require.Equal(t, newName, cf.Key)
			require.False(t, cf.Shown)
			found = true
			break
		}
	}
	require.True(t, found)
}

func TestUpsertRoomCustomFieldValue(t *testing.T) {
	cf := createRandomCustomField(t)
	room := createRandomRoom(t)
	val1 := utils.RandomString(10)

	arg := UpsertRoomCustomFieldValueParams{
		RoomID:        room.ID,
		CustomfieldID: cf.ID,
		Value:         val1,
	}

	err := testQuery.UpsertRoomCustomFieldValue(context.Background(), arg)
	require.NoError(t, err)

	vals, err := testQuery.ListRoomCustomFieldValues(context.Background(), room.ID)
	require.NoError(t, err)
	require.Len(t, vals, 1)
	require.Equal(t, val1, vals[0].Value)

	// Update
	val2 := utils.RandomString(10)
	arg.Value = val2
	err = testQuery.UpsertRoomCustomFieldValue(context.Background(), arg)
	require.NoError(t, err)

	vals, err = testQuery.ListRoomCustomFieldValues(context.Background(), room.ID)
	require.NoError(t, err)
	require.Len(t, vals, 1)
	require.Equal(t, val2, vals[0].Value)
}

func TestListRoomCustomFieldValuesBatch(t *testing.T) {
	cf := createRandomCustomField(t)
	room1 := createRandomRoom(t)
	room2 := createRandomRoom(t)
	val1 := utils.RandomString(10)
	val2 := utils.RandomString(10)

	err := testQuery.UpsertRoomCustomFieldValue(context.Background(), UpsertRoomCustomFieldValueParams{
		RoomID:        room1.ID,
		CustomfieldID: cf.ID,
		Value:         val1,
	})
	require.NoError(t, err)

	err = testQuery.UpsertRoomCustomFieldValue(context.Background(), UpsertRoomCustomFieldValueParams{
		RoomID:        room2.ID,
		CustomfieldID: cf.ID,
		Value:         val2,
	})
	require.NoError(t, err)

	vals, err := testQuery.ListRoomCustomFieldValuesBatch(context.Background(), []int64{room1.ID, room2.ID})
	require.NoError(t, err)
	require.Len(t, vals, 2)
}
