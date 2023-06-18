package db

import (
	"context"
	"database/sql"
	"github.com/stretchr/testify/suite"
	"simple_bank/utils"
	"testing"
)

type AccountSuite struct {
	suite.Suite
}

//func (s *AccountSuite) SetupSuite() {
//	// start the server
//	println("setup suite")
//
//}

func (s *AccountSuite) SetupTest() {
	connection.DeleteAccountTable()

}

//	func (s *AccountSuite) TearDownSuite() {
//		// stop the server
//		println("teardown suite")
//	}
func (s *AccountSuite) TestGetAccount() {
	newAccount, _ := s.createRandomAccount()
	getAccount, err := testQueries.GetAccount(context.Background(), newAccount.ID)
	s.Require().NoError(err)

	s.Equal(newAccount.ID, getAccount.ID)
	s.Equal(newAccount.Owner, getAccount.Owner)
	s.Equal(newAccount.Balance, getAccount.Balance)
	s.Equal(newAccount.Currency, getAccount.Currency)
	s.WithinDuration(newAccount.CreatedAt, getAccount.CreatedAt, 0)
}
func (s *AccountSuite) TestCreateAccount() {
	account, randomArgs := s.createRandomAccount()

	s.Equal(randomArgs.Owner, account.Owner)
	s.Equal(randomArgs.Balance, account.Balance)
	s.Equal(randomArgs.Currency, account.Currency)
	s.NotZero(account.ID)
	s.NotZero(account.CreatedAt)
}

func (s *AccountSuite) createRandomAccount() (*Account, *CreateAccountParams) {
	randomArgs := CreateAccountParams{
		Owner:    utils.RandomOwner(),
		Balance:  utils.RandomMoney(),
		Currency: utils.RandomCurrency(),
	}

	account, err := testQueries.CreateAccount(context.Background(), randomArgs)
	s.Require().NoError(err)
	s.Require().NotEmpty(account)

	return &account, &randomArgs
}

func (s *AccountSuite) TestUpdateAccount() {
	newAccount, _ := s.createRandomAccount()
	updateArgs := UpdateAccountParams{
		ID:      newAccount.ID,
		Balance: utils.RandomMoney(),
	}

	updatedAccount, err := testQueries.UpdateAccount(context.Background(), updateArgs)
	s.Require().NoError(err)

	s.Equal(newAccount.ID, updatedAccount.ID)
	s.Equal(newAccount.Owner, updatedAccount.Owner)
	s.Equal(updateArgs.Balance, updatedAccount.Balance)
	s.Equal(newAccount.Currency, updatedAccount.Currency)
	s.WithinDuration(newAccount.CreatedAt, updatedAccount.CreatedAt, 0)
}

func (s *AccountSuite) TestDeleteAccount() {
	newAccount, _ := s.createRandomAccount()
	err := testQueries.DeleteAccount(context.Background(), newAccount.ID)
	s.Require().NoError(err)

	getAccount, err := testQueries.GetAccount(context.Background(), newAccount.ID)
	s.Require().Error(err)
	s.EqualError(err, sql.ErrNoRows.Error())
	s.Empty(getAccount)
}

func (s *AccountSuite) TestListAccounts() {
	for i := 0; i < 10; i++ {
		s.createRandomAccount()
	}

	args := ListAccountsParams{
		Limit:  5,
		Offset: 5,
	}

	accounts, err := testQueries.ListAccounts(context.Background(), args)
	s.Require().NoError(err)
	s.Require().Len(accounts, 5)

	for _, account := range accounts {
		s.NotEmpty(account)
	}
}

func TestAccountSuite(t *testing.T) {
	suite.Run(t, new(AccountSuite))
}
