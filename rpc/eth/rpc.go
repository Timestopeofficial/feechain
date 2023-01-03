package eth

import (
	"context"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/Timestopeofficial/feechain/eth/rpc"
	"github.com/Timestopeofficial/feechain/fch"
	internal_common "github.com/Timestopeofficial/feechain/internal/common"
)

// PublicEthService provides an API to access to the Eth endpoints for the Feechain blockchain.
type PublicEthService struct {
	fch *fch.Feechain
}

// NewPublicEthService creates a new API for the RPC interface
func NewPublicEthService(fch *fch.Feechain, namespace string) rpc.API {
	if namespace == "" {
		namespace = "eth"
	}

	return rpc.API{
		Namespace: namespace,
		Version:   "1.0",
		Service:   &PublicEthService{fch},
		Public:    true,
	}
}

// GetBalance returns the amount of Atto for the given address in the state of the
// given block number. The rpc.LatestBlockNumber and rpc.PendingBlockNumber meta
// block numbers are also allowed.
func (s *PublicEthService) GetBalance(
	ctx context.Context, address string, blockNr rpc.BlockNumber,
) (*hexutil.Big, error) {
	addr, err := internal_common.ParseAddr(address)
	if err != nil {
		return nil, err
	}
	balance, err := s.fch.GetBalance(ctx, addr, blockNr)
	if err != nil {
		return nil, err
	}
	return (*hexutil.Big)(balance), nil
}
