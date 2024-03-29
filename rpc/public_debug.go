package rpc

import (
	"context"

	"github.com/ethereum/go-ethereum/log"
	"github.com/Timestopeofficial/feechain/eth/rpc"
	"github.com/Timestopeofficial/feechain/fch"
	"github.com/Timestopeofficial/feechain/internal/utils"
)

// PublicDebugService Internal JSON RPC for debugging purpose
type PublicDebugService struct {
	fch     *fch.Feechain
	version Version
}

// NewPublicDebugAPI creates a new API for the RPC interface
func NewPublicDebugAPI(fch *fch.Feechain, version Version) rpc.API {
	return rpc.API{
		Namespace: version.Namespace(),
		Version:   APIVersion,
		Service:   &PublicDebugService{fch, version},
		Public:    false,
	}
}

// SetLogVerbosity Sets log verbosity on runtime
// curl -H "Content-Type: application/json" -d '{"method":"fch_setLogVerbosity","params":[5],"id":1}' http://127.0.0.1:9500
func (s *PublicDebugService) SetLogVerbosity(ctx context.Context, level int) (map[string]interface{}, error) {
	if level < int(log.LvlCrit) || level > int(log.LvlTrace) {
		return nil, ErrInvalidLogLevel
	}

	verbosity := log.Lvl(level)
	utils.SetLogVerbosity(verbosity)
	return map[string]interface{}{"verbosity": verbosity.String()}, nil
}
