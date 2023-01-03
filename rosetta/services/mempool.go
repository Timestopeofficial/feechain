package services

import (
	"context"
	"math/big"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	ethCommon "github.com/ethereum/go-ethereum/common"
	fchTypes "github.com/Timestopeofficial/feechain/core/types"
	"github.com/Timestopeofficial/feechain/fch"
	"github.com/Timestopeofficial/feechain/rosetta/common"
	"github.com/Timestopeofficial/feechain/staking"
)

// MempoolAPI implements the server.MempoolAPIServicer interface
type MempoolAPI struct {
	fch *fch.Feechain
}

// NewMempoolAPI creates a new instance of MempoolAPI
func NewMempoolAPI(fch *fch.Feechain) server.MempoolAPIServicer {
	return &MempoolAPI{
		fch: fch,
	}
}

// Mempool implements the /mempool endpoint.
func (s *MempoolAPI) Mempool(
	ctx context.Context, req *types.NetworkRequest,
) (*types.MempoolResponse, *types.Error) {
	if err := assertValidNetworkIdentifier(req.NetworkIdentifier, s.fch.ShardID); err != nil {
		return nil, err
	}

	pool, err := s.fch.GetPoolTransactions()
	if err != nil {
		return nil, common.NewError(common.CatchAllError, map[string]interface{}{
			"message": "unable to fetch pool transactions",
		})
	}
	txIDs := make([]*types.TransactionIdentifier, pool.Len())
	for i, tx := range pool {
		txIDs[i] = &types.TransactionIdentifier{
			Hash: tx.Hash().String(),
		}
	}
	return &types.MempoolResponse{
		TransactionIdentifiers: txIDs,
	}, nil
}

// MempoolTransaction implements the /mempool/transaction endpoint.
func (s *MempoolAPI) MempoolTransaction(
	ctx context.Context, req *types.MempoolTransactionRequest,
) (*types.MempoolTransactionResponse, *types.Error) {
	if err := assertValidNetworkIdentifier(req.NetworkIdentifier, s.fch.ShardID); err != nil {
		return nil, err
	}

	hash := ethCommon.HexToHash(req.TransactionIdentifier.Hash)
	poolTx := s.fch.GetPoolTransaction(hash)
	if poolTx == nil {
		return nil, &common.TransactionNotFoundError
	}

	senderAddr, _ := poolTx.SenderAddress()
	estLog := &fchTypes.Log{
		Address:     senderAddr,
		Topics:      []ethCommon.Hash{staking.CollectRewardsTopic},
		Data:        big.NewInt(0).Bytes(),
		BlockNumber: s.fch.CurrentBlock().NumberU64(),
	}

	// Contract related information for pending transactions is not reported
	estReceipt := &fchTypes.Receipt{
		PostState:         []byte{},
		Status:            fchTypes.ReceiptStatusSuccessful, // Assume transaction will succeed
		CumulativeGasUsed: poolTx.GasLimit(),
		Bloom:             [256]byte{},
		Logs:              []*fchTypes.Log{estLog},
		TxHash:            poolTx.Hash(),
		ContractAddress:   ethCommon.Address{},
		GasUsed:           poolTx.GasLimit(),
	}

	respTx, err := FormatTransaction(poolTx, estReceipt, &ContractInfo{}, true)
	if err != nil {
		return nil, err
	}

	return &types.MempoolTransactionResponse{
		Transaction: respTx,
	}, nil
}
