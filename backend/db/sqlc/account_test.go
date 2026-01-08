package db

import (
	"context"
	"database/sql"
	"testing"

	utils "github.com/project-available/available/utils"
	"github.com/stretchr/testify/require"
)

func createRandomAccount(t *testing.T) Account {
	hashedPassword, err := utils.HashPassword(utils.RandomString(10))
	require.NoError(t, err)

	arg := CreateAccountParams{
		Name:           utils.RandomString(15),
		Role:           utils.RandomRole(),
		Email:          utils.RandomEmail(),
		HashedPassword: hashedPassword,
		Phone:          utils.RandomPhone(),
		StudentID:      utils.RandomStudentID(),
	}
	account, err := testQuery.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Name, account.Name)
	require.Equal(t, arg.Role, account.Role)
	require.Equal(t, arg.Email, account.Email)
	require.Equal(t, arg.Phone, account.Phone)
	require.Equal(t, arg.StudentID, account.StudentID)
	require.Equal(t, arg.HashedPassword, hashedPassword)
	require.NotZero(t, account.ID)
	return account
}

func TestCreateAccount(t *testing.T) {
	createRandomAccount(t)
}

func TestGetAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQuery.GetAccount(context.Background(), account1.StudentID)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Name, account2.Name)
	require.Equal(t, account1.Role, account2.Role)
	require.Equal(t, account1.Email, account2.Email)
	require.Equal(t, account1.Phone, account2.Phone)
	require.Equal(t, account1.HashedPassword, account2.HashedPassword)
	require.Equal(t, account1.StudentID, account2.StudentID)
	require.Equal(t, account1.IsDelete, account2.IsDelete)
}

func TestGetAccountByEmail(t *testing.T) {
	account1 := createRandomAccount(t)
	account2, err := testQuery.GetAccountByEmail(context.Background(), account1.Email)

	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, account1.Name, account2.Name)
	require.Equal(t, account1.Role, account2.Role)
	require.Equal(t, account1.Email, account2.Email)
	require.Equal(t, account1.Phone, account2.Phone)
	require.Equal(t, account1.HashedPassword, account2.HashedPassword)
	require.Equal(t, account1.StudentID, account2.StudentID)
	require.Equal(t, account1.IsDelete, account2.IsDelete)
}

func TestListAccounts(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomAccount(t)
	}

	arg := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQuery.ListAccounts(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, accounts, 5)

	for _, account := range accounts {
		require.NotEmpty(t, account)
	}

}

func TestUpdateAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	newPhoneNumber := utils.RandomPhone()
	newName := utils.RandomString(5)

	arg := UpdateAccountParams{
		ID:    account1.ID,
		Name:  newName,
		Phone: newPhoneNumber,
	}

	account2, err := testQuery.UpdateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, account2)

	require.Equal(t, newName, account2.Name)
	require.Equal(t, account1.Role, account2.Role)
	require.Equal(t, account1.Email, account2.Email)
	require.Equal(t, newPhoneNumber, account2.Phone)
	require.Equal(t, account1.HashedPassword, account2.HashedPassword)
	require.Equal(t, account1.StudentID, account2.StudentID)
	require.Equal(t, account1.IsDelete, account2.IsDelete)
}

func TestDeleteAccount(t *testing.T) {
	account1 := createRandomAccount(t)
	require.False(t, account1.IsDelete)

	err := testQuery.DeleteAccount(context.Background(), account1.StudentID)
	require.NoError(t, err)

	account2, err := testQuery.GetAccount(context.Background(), account1.StudentID)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, account2)

}
