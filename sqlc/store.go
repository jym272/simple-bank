package db

import (
	"context"
	"database/sql"
)

type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	q := New(tx)
	err = fn(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return rbErr
		}
		return err
	}
	return tx.Commit()
}

type TransferTxParams struct {
	FromAccountID int64 `json:"from_account_id"`
	ToAccountID   int64 `json:"to_account_id"`
	Amount        int64 `json:"amount"`
}

type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

var txKey = struct{}{}

// Perform a money transfer between two accounts. The `amount` value must be positive,
// and the `from_account_id` must have sufficient funds to complete the transfer.
// This function returns the newly-created transfer as well as the updated
// account entries.

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult
	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}

		// prevent deadlocks
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		if err != nil {
			return err
		}

		//result.FromAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		//	Amount: -arg.Amount,
		//	ID:     arg.FromAccountID,
		//})
		//if err != nil {
		//	return err
		//}
		//
		//result.ToAccount, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		//	Amount: arg.Amount,
		//	ID:     arg.ToAccountID,
		//})
		//if err != nil {
		//	return err
		//}

		return nil
	})

	return result, err
}

func addMoney(ctx context.Context, q *Queries, accountIDA int64, amountA int64, accountIDB int64, amountB int64) (accountA Account,
	accountB Account, err error) {
	accountA, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amountA,
		ID:     accountIDA,
	})
	if err != nil {
		return
	}

	accountB, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		Amount: amountB,
		ID:     accountIDB,
	})
	if err != nil {
		return
	}

	return
}
