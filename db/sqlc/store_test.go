package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)



func TestTransferTxDeadlock(t *testing.T) {
	store := NewStore(testDB)
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	//debuging : 
	fmt.Println(">> before : " ,"account1Balance is :" , account1.Balance , "account2Balance is :" , account2.Balance)
	//run n concurrent transfer transaction!
	n := 10
	amount := int64(10)
	errs := make(chan error)
	
	// run n concurrent transfer transaction
	for i := 0; i < n; i++ {
		fromAccountID :=account1.ID
		toAccountID := account2.ID

		if i %2==1 {
			fromAccountID= account2.ID
			toAccountID = account1.ID
		}
		go func() {
			_, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: fromAccountID,
				ToAccountID:   toAccountID,
				Amount:        amount,
			})
			errs <- err
		}()
	}
	//cehck results
	for i := 0 ; i < n ; i++ {
		err := <-errs
		require.NoError(t , err)
	}
	//check the final update balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	
	fmt.Println(">> After : " ,"account1Balance is :" , updatedAccount1.Balance , "account2Balance is :" , updatedAccount2.Balance)
	require.Equal(t, account1.Balance, updatedAccount1.Balance)
	require.Equal(t, account2.Balance, updatedAccount2.Balance)


}
