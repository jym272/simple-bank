package db

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
)

type StoreSuite struct {
	Suite
	store *Store
}

func (s *StoreSuite) SetupSuite() {
	s.store = NewStore(connection.GetConnection())
}

func (s *StoreSuite) SetupTest() {
	connection.DeleteAccountTable()

}
func (s *StoreSuite) TestTransferTx() {

	accountA, _ := s.createRandomAccount()
	accountB, _ := s.createRandomAccount()

	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := s.store.TransferTx(ctx, TransferTxParams{
				FromAccountID: accountA.ID,
				ToAccountID:   accountB.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}

	var existed = make(map[int]bool)

	for i := 0; i < n; i++ {
		err := <-errs
		s.Require().NoError(err)

		result := <-results
		// print result
		//fmt.Printf("%+v\n", result)
		s.Require().Equal(accountA.ID, result.FromAccount.ID)
		s.Require().Equal(accountB.ID, result.ToAccount.ID)
		diffA := accountA.Balance - result.FromAccount.Balance
		diffB := result.ToAccount.Balance - accountB.Balance
		s.Require().Equal(diffB, diffA)
		s.Require().True(diffB > 0)
		s.Require().True(diffB%amount == 0)
		k := int(diffB / amount)
		s.Require().True(k >= 1 && k <= n)
		s.Require().False(existed[k])
		existed[k] = true

		// check transfer
		transfer := result.Transfer
		s.Equal(accountA.ID, transfer.FromAccountID)
		s.Equal(accountB.ID, transfer.ToAccountID)
		s.Equal(amount, transfer.Amount)

		_, err = s.store.GetTransfer(context.Background(), transfer.ID)
		s.Require().NoError(err)

		// check entries
		fromEntry := result.FromEntry
		s.Equal(accountA.ID, fromEntry.AccountID)
		s.Equal(-amount, fromEntry.Amount)

		_, err = s.store.GetEntry(context.Background(), fromEntry.ID)
		s.Require().NoError(err)

		toEntry := result.ToEntry
		s.Equal(accountB.ID, toEntry.AccountID)
		s.Equal(amount, toEntry.Amount)

		_, err = s.store.GetEntry(context.Background(), toEntry.ID)
		s.Require().NoError(err)

	}

	updatedAccountA, err := s.store.GetAccount(context.Background(), accountA.ID)
	s.Require().NoError(err)

	updatedAccountB, err := s.store.GetAccount(context.Background(), accountB.ID)
	s.Require().NoError(err)

	s.Equal(accountA.Balance-int64(n)*amount, updatedAccountA.Balance)
	s.Equal(accountB.Balance+int64(n)*amount, updatedAccountB.Balance)

}

func (s *StoreSuite) TestTransferTxDeadlock() {

	accountA, _ := s.createRandomAccount()
	accountB, _ := s.createRandomAccount()

	n := 10
	amount := int64(10)
	errs := make(chan error)

	for i := 0; i < n; i++ {
		fromAccountID := accountA.ID
		toAccountID := accountB.ID
		if i%2 == 1 {
			fromAccountID = accountB.ID
			toAccountID = accountA.ID
		}
		go func() {
			_, err := s.store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}

	for i := 0; i < n; i++ {
		err := <-errs
		s.Require().NoError(err)
	}

	updatedAccountA, err := s.store.GetAccount(context.Background(), accountA.ID)
	s.Require().NoError(err)

	updatedAccountB, err := s.store.GetAccount(context.Background(), accountB.ID)
	s.Require().NoError(err)

	s.Equal(accountA.Balance, updatedAccountA.Balance)
	s.Equal(accountB.Balance, updatedAccountB.Balance)

}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(StoreSuite))
}
