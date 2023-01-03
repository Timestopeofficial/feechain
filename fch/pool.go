package fch

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/Timestopeofficial/feechain/core/types"
)

// GetPoolStats returns the number of pending and queued transactions
func (fch *Feechain) GetPoolStats() (pendingCount, queuedCount int) {
	return fch.TxPool.Stats()
}

// GetPoolNonce ...
func (fch *Feechain) GetPoolNonce(ctx context.Context, addr common.Address) (uint64, error) {
	return fch.TxPool.State().GetNonce(addr), nil
}

// GetPoolTransaction ...
func (fch *Feechain) GetPoolTransaction(hash common.Hash) types.PoolTransaction {
	return fch.TxPool.Get(hash)
}

// GetPendingCXReceipts ..
func (fch *Feechain) GetPendingCXReceipts() []*types.CXReceiptsProof {
	return fch.NodeAPI.PendingCXReceipts()
}

// GetPoolTransactions returns pool transactions.
func (fch *Feechain) GetPoolTransactions() (types.PoolTransactions, error) {
	pending, err := fch.TxPool.Pending()
	if err != nil {
		return nil, err
	}
	queued, err := fch.TxPool.Queued()
	if err != nil {
		return nil, err
	}
	var txs types.PoolTransactions
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	for _, batch := range queued {
		txs = append(txs, batch...)
	}
	return txs, nil
}

func (fch *Feechain) SuggestPrice(ctx context.Context) (*big.Int, error) {
	return fch.gpo.SuggestPrice(ctx)
}
