package rpc

import (
	"context"

	"github.com/Timestopeofficial/feechain/eth/rpc"
	"github.com/Timestopeofficial/feechain/fch"
)

// PrivateDebugService Internal JSON RPC for debugging purpose
type PrivateDebugService struct {
	fch     *fch.Feechain
	version Version
}

// NewPrivateDebugAPI creates a new API for the RPC interface
// TODO(dm): expose public via config
func NewPrivateDebugAPI(fch *fch.Feechain, version Version) rpc.API {
	return rpc.API{
		Namespace: version.Namespace(),
		Version:   APIVersion,
		Service:   &PrivateDebugService{fch, version},
		Public:    false,
	}
}

// ConsensusViewChangingID return the current view changing ID to RPC
func (s *PrivateDebugService) ConsensusViewChangingID(
	ctx context.Context,
) uint64 {
	return s.fch.NodeAPI.GetConsensusViewChangingID()
}

// ConsensusCurViewID return the current view ID to RPC
func (s *PrivateDebugService) ConsensusCurViewID(
	ctx context.Context,
) uint64 {
	return s.fch.NodeAPI.GetConsensusCurViewID()
}

// GetConsensusMode return the current consensus mode
func (s *PrivateDebugService) GetConsensusMode(
	ctx context.Context,
) string {
	return s.fch.NodeAPI.GetConsensusMode()
}

// GetConsensusPhase return the current consensus mode
func (s *PrivateDebugService) GetConsensusPhase(
	ctx context.Context,
) string {
	return s.fch.NodeAPI.GetConsensusPhase()
}

// GetConfig get feechain config
func (s *PrivateDebugService) GetConfig(
	ctx context.Context,
) (StructuredResponse, error) {
	return NewStructuredResponse(s.fch.NodeAPI.GetConfig())
}

// GetLastSigningPower get last signed power
func (s *PrivateDebugService) GetLastSigningPower(
	ctx context.Context,
) (float64, error) {
	return s.fch.NodeAPI.GetLastSigningPower()
}
