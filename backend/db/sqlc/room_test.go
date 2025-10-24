package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/project-available/available.git/utils"
	"github.com/stretchr/testify/require"
)

func createRandomRoom(t *testing.T) Room {
	arg := CreateRoomParams{
		Name:     utils.RandomString(10),
		Location: utils.RandomString(15),
		Image:    utils.RandomString(20),
	}
	room, err := testQuery.CreateRoom(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, room)

	require.Equal(t, arg.Name, room.Name)
	require.Equal(t, arg.Location, room.Location)
	require.Equal(t, arg.Image, room.Image)
	require.NotZero(t, room.ID)
	return room
}
func TestCreateRoom(t *testing.T) {
	createRandomRoom(t)
}
func TestListRooms(t *testing.T) {
	for i := 0; i < 6; i++ {
		createRandomRoom(t)
	}
	arg := ListRoomsParams{
		Limit:  5,
		Offset: 5,
	}
	rooms, err := testQuery.ListRooms(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, rooms, 5)
	for _, room := range rooms {
		require.NotEmpty(t, room)
	}
}

func TestUpdateRoom (t *testing.T){
	room1 := createRandomRoom(t);
	newName := utils.RandomString(10)
	newLocation := utils.RandomString(15)
	newImage := utils.RandomString(10)

	arg := UpdateRoomParams{
		ID: room1.ID,
		Name: newName,
		Location: newLocation,
		Image: newImage,
	}
	room2, err := testQuery.UpdateRoom(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, room2)

	require.Equal(t, room1.ID, room2.ID)
	require.Equal(t, newName, room2.Name)
	require.Equal(t, newLocation,  room2.Location)
	require.Equal(t, newImage, room2.Image)

}

func TestDeleteRoom(t *testing.T){
	room1 := createRandomRoom(t);
	require.False(t,room1.IsDelete)
	err := testQuery.DeleteRoom( context.Background(), room1.ID)
	require.NoError(t, err)
	room2, err := testQuery.GetRoom(context.Background(), room1.ID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, room2)
}
