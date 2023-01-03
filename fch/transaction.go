package fch

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/Timestopeofficial/feechain/core"
	"github.com/Timestopeofficial/feechain/core/rawdb"
	"github.com/Timestopeofficial/feechain/core/types"
	"github.com/Timestopeofficial/feechain/eth/rpc"
)

// SendTx ...
func (fch *Feechain) SendTx(ctx context.Context, signedTx *types.Transaction) error {
	tx, _, _, _ := rawdb.ReadTransaction(fch.chainDb, signedTx.Hash())
	if tx == nil {
		return fch.NodeAPI.AddPendingTransaction(signedTx)
	}
	return ErrFinalizedTransaction
}

// ResendCx retrieve blockHash from txID and add blockHash to CxPool for resending
// Note that cross shard txn is only for regular txns, not for staking txns, so the input txn hash
// is expected to be regular txn hash
func (fch *Feechain) ResendCx(ctx context.Context, txID common.Hash) (uint64, bool) {
	blockHash, blockNum, index := fch.BlockChain.ReadTxLookupEntry(txID)
	if blockHash == (common.Hash{}) {
		return 0, false
	}

	blk := fch.BlockChain.GetBlockByHash(blockHash)
	if blk == nil {
		return 0, false
	}

	txs := blk.Transactions()
	// a valid index is from 0 to len-1
	if int(index) > len(txs)-1 {
		return 0, false
	}
	tx := txs[int(index)]

	// check whether it is a valid cross shard tx
	if tx.ShardID() == tx.ToShardID() || blk.Header().ShardID() != tx.ShardID() {
		return 0, false
	}
	entry := core.CxEntry{blockHash, tx.ToShardID()}
	success := fch.CxPool.Add(entry)
	return blockNum, success
}

// GetReceipts ...
func (fch *Feechain) GetReceipts(ctx context.Context, hash common.Hash) (types.Receipts, error) {
	return fch.BlockChain.GetReceiptsByHash(hash), nil
}

// GetTransactionsHistory returns list of transactions hashes of address.
func (fch *Feechain) GetTransactionsHistory(address, txType, order string) ([]common.Hash, error) {
	return fch.NodeAPI.GetTransactionsHistory(address, txType, order)
}

// GetAccountNonce returns the nonce value of the given address for the given block number
func (fch *Feechain) GetAccountNonce(
	ctx context.Context, address common.Address, blockNum rpc.BlockNumber) (uint64, error) {
	state, _, err := fch.StateAndHeaderByNumber(ctx, blockNum)
	if state == nil || err != nil {
		return 0, err
	}
	return state.GetNonce(address), state.Error()
}

// GetTransactionsCount returns the number of regular transactions of address.
func (fch *Feechain) GetTransactionsCount(address, txType string) (uint64, error) {
	return fch.NodeAPI.GetTransactionsCount(address, txType)
}

// GetCurrentTransactionErrorSink ..
func (fch *Feechain) GetCurrentTransactionErrorSink() types.TransactionErrorReports {
	return fch.NodeAPI.ReportPlainErrorSink()
}
