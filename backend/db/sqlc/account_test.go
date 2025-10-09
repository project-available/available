package db

import (
	"context"
	"testing"

	utils "github.com/project-available/available.git/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account{
	arg := CreateAccountParams{
		Name: utils.RandomString(15),
		Role: utils.RandomRole(),
		Email: utils.RandomEmail(),
		Password: "12345",
		Phone: utils.RandomPhone(),
		StudentID: utils.RandomStudentID(),
	}
	account, err := testQuery.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Name, account.Name)
	require.Equal(t, arg.Role, account.Role)
	require.Equal(t, arg.Email, account.Email)
	require.Equal(t, arg.Phone, account.Phone)
	require.Equal(t, arg.StudentID, account.StudentID)
	require.NotZero(t,account.ID)
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQuery.GetAccount(context.Background(), account1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Name, account2.Name)
	require.Equal(t, account1.Role, account2.Role)
	require.Equal(t, account1.Email, account2.Email)
	require.Equal(t, account1.Phone, account2.Phone)
	require.Equal(t, account1.StudentID, account2.StudentID)
}
