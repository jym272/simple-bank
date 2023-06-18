package db

import (
	"context"
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

	// run n concurrent transfer transactions
	n := 5
	amount := int64(10)
	errs := make(chan error)
	results := make(chan TransferTxResult)
	for i := 0; i < n; i++ {
		go func() {
			result, err := s.store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: accountA.ID,
				ToAccountID:   accountB.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
		}()
	}
	for i := 0; i < n; i++ {
		err := <-errs
		s.Require().NoError(err)

		result := <-results
		// print result
		//fmt.Printf("%+v\n", result)
		// TODO
		//s.Equal(accountA.ID, result.FromAccount.ID)
		//s.Equal(accountB.ID, result.ToAccount.ID)
		//s.Equal(accountA.Balance-amount, result.FromAccount.Balance)
		//s.Equal(accountB.Balance+amount, result.ToAccount.Balance)

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

		// TODO: check accounts' balance

	}

}

func TestStoreSuite(t *testing.T) {
	suite.Run(t, new(StoreSuite))
}
