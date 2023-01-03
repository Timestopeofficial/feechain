package rpc

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/Timestopeofficial/feechain/eth/rpc"
	"github.com/Timestopeofficial/feechain/fch"
)

// PublicFeechainService provides an API to access Feechain related information.
// It offers only methods that operate on public data that is freely available to anyone.
type PublicFeechainService struct {
	fch     *fch.Feechain
	version Version
}

// NewPublicFeechainAPI creates a new API for the RPC interface
func NewPublicFeechainAPI(fch *fch.Feechain, version Version) rpc.API {
	return rpc.API{
		Namespace: version.Namespace(),
		Version:   APIVersion,
		Service:   &PublicFeechainService{fch, version},
		Public:    true,
	}
}

// ProtocolVersion returns the current Feechain protocol version this node supports
// Note that the return type is an interface to account for the different versions
func (s *PublicFeechainService) ProtocolVersion(
	ctx context.Context,
) (interface{}, error) {
	// Format response according to version
	switch s.version {
	case V1, Eth:
		return hexutil.Uint(s.fch.ProtocolVersion()), nil
	case V2:
		return s.fch.ProtocolVersion(), nil
	default:
		return nil, ErrUnknownRPCVersion
	}
}

// Syncing returns false in case the node is currently not syncing with the network. It can be up to date or has not
// yet received the latest block headers from its pears. In case it is synchronizing:
// - startingBlock: block number this node started to synchronise from
// - currentBlock:  block number this node is currently importing
// - highestBlock:  block number of the highest block header this node has received from peers
// - pulledStates:  number of state entries processed until now
// - knownStates:   number of known state entries that still need to be pulled
func (s *PublicFeechainService) Syncing(
	ctx context.Context,
) (interface{}, error) {
	// TODO(dm): find our Downloader module for syncing blocks
	return false, nil
}

// GasPrice returns a suggestion for a gas price.
// Note that the return type is an interface to account for the different versions
func (s *PublicFeechainService) GasPrice(ctx context.Context) (interface{}, error) {
	price, err := s.fch.SuggestPrice(ctx)
	if err != nil || price.Cmp(big.NewInt(1e12)) < 0 {
		price = big.NewInt(1e12)
	}
	// Format response according to version
	switch s.version {
	case V1, Eth:
		return (*hexutil.Big)(price), nil
	case V2:
		return price.Uint64(), nil
	default:
		return nil, ErrUnknownRPCVersion
	}
}

// GetNodeMetadata produces a NodeMetadata record, data is from the answering RPC node
func (s *PublicFeechainService) GetNodeMetadata(
	ctx context.Context,
) (StructuredResponse, error) {
	// Response output is the same for all versions
	return NewStructuredResponse(s.fch.GetNodeMetadata())
}

// GetPeerInfo produces a NodePeerInfo record
func (s *PublicFeechainService) GetPeerInfo(
	ctx context.Context,
) (StructuredResponse, error) {
	// Response output is the same for all versions
	return NewStructuredResponse(s.fch.GetPeerInfo())
}

// GetNumPendingCrossLinks returns length of fch.BlockChain.ReadPendingCrossLinks()
func (s *PublicFeechainService) GetNumPendingCrossLinks() (int, error) {
	links, err := s.fch.BlockChain.ReadPendingCrossLinks()
	if err != nil {
		return 0, err
	}

	return len(links), nil
}
