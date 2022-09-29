package slash

import (
	"math/big"

	"github.com/Timestopeofficial/feechain/core/types"
	"github.com/Timestopeofficial/feechain/internal/params"
	"github.com/Timestopeofficial/feechain/shard"
)

// CommitteeReader ..
type CommitteeReader interface {
	Config() *params.ChainConfig
	ReadShardState(epoch *big.Int) (*shard.State, error)
	CurrentBlock() *types.Block
}
