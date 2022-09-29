package v2

import (
	"context"
	"math/big"

	"github.com/Timestopeofficial/feechain/eth/rpc"
	"github.com/Timestopeofficial/feechain/hmy"
	internal_common "github.com/Timestopeofficial/feechain/internal/common"
)

// PublicLegacyService provides an API to access the Harmony blockchain.
// Services here are legacy methods, specific to the V1 RPC that can be deprecated in the future.
type PublicLegacyService struct {
	hmy *hmy.Harmony
}

// NewPublicLegacyAPI creates a new API for the RPC interface
func NewPublicLegacyAPI(hmy *hmy.Harmony, namespace string) rpc.API {
	if namespace == "" {
		namespace = "fchv2"
	}

	return rpc.API{
		Namespace: namespace,
		Version:   "1.0",
		Service:   &PublicLegacyService{hmy},
		Public:    true,
	}
}

// GetBalance returns the amount of Atto for the given address in the state of the
// given block number.
func (s *PublicLegacyService) GetBalance(
	ctx context.Context, address string,
) (*big.Int, error) {
	addr, err := internal_common.ParseAddr(address)
	if err != nil {
		return nil, err
	}
	balance, err := s.hmy.GetBalance(ctx, addr, rpc.BlockNumber(-1))
	if err != nil {
		return nil, err
	}
	return balance, nil
}
