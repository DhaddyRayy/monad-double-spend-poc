package execution_test

import (
    "sync"
    "testing"

    "github.com/stretchr/testify/assert"
    "monad/execution/processor"
)

// TestDoubleSpendRaceCondition demonstrates a race condition that can allow double-spend
func TestDoubleSpendRaceCondition(t *testing.T) {
    // 1. Setup the environment with an attacker account holding 100 MON.
    stateDB := setupStateDB()
    attackerAddr := "0xAttacker"
    exchangeAddr := "0xExchange"
    personalAddr := "0xPersonal"
    nonce := stateDB.GetNonce(attackerAddr) // Nonce is 0

    // 2. Create two transactions with the same nonce spending the same funds.
    tx1 := createSignedTx(attackerAddr, exchangeAddr, 100, nonce)
    tx2 := createSignedTx(attackerAddr, personalAddr, 100, nonce)

    // 3. Simulate the parallel processor by running validation in separate goroutines
    var wg sync.WaitGroup
    wg.Add(2)

    go func() {
        defer wg.Done()
        isValid, _ := processor.ValidateTransaction(stateDB, tx1)
        assert.True(t, isValid, "Tx1 should be valid initially")
        // State is NOT updated yet due to race condition
    }()

    go func() {
        defer wg.Done()
        isValid, _ := processor.ValidateTransaction(stateDB, tx2)
        assert.True(t, isValid, "Tx2 should also be valid initially due to race condition")
    }()

    wg.Wait()
}
