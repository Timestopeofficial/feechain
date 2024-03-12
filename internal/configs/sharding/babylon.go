package shardingconfig

import (
	"math/big"

	"github.com/Timestopeofficial/feechain/internal/params"

	"github.com/Timestopeofficial/feechain/numeric"

	"github.com/Timestopeofficial/feechain/internal/genesis"
)

const (
	babylonBlocksPerEpoch     = 28800

	babylonVdfDifficulty = 50000 // This takes about 100s to finish the vdf

	// BabylonHTTPPattern is the http pattern for asadal.
	BabylonHTTPPattern = "https://rpc.s%d.babylon.mojaik.com"
	// BabylonWSPattern is the websocket pattern for asadal.
	BabylonWSPattern = "wss://wss.s%d.babylon.mojaik.com"
)

// BabylonSchedule is the asadal sharding configuration schedule.
var BabylonSchedule babylonSchedule

type babylonSchedule struct{}

func (ms babylonSchedule) InstanceForEpoch(epoch *big.Int) Instance {
	switch {
	case epoch.Cmp(params.BabylonChainConfig.StakingEpoch) >= 0:
		return babylonV1
	default: // genesis
		return babylonV0
	}
}

func (ms babylonSchedule) BlocksPerEpoch() uint64 {
	return babylonBlocksPerEpoch
}

func (ms babylonSchedule) CalcEpochNumber(blockNum uint64) *big.Int {
	return big.NewInt(int64(blockNum / ms.BlocksPerEpoch()))
}

func (ms babylonSchedule) IsLastBlock(blockNum uint64) bool {
	return (blockNum % ms.BlocksPerEpoch() == ms.BlocksPerEpoch() - 1)
}

func (ms babylonSchedule) EpochLastBlock(epochNum uint64) uint64 {
		return ms.BlocksPerEpoch() * (epochNum + 1) - 1
}

func (ms babylonSchedule) VdfDifficulty() int {
	return babylonVdfDifficulty
}

func (ms babylonSchedule) GetNetworkID() NetworkID {
	return Babylon
}

// GetShardingStructure is the sharding structure for asadal.
func (ms babylonSchedule) GetShardingStructure(numShard, shardID int) []map[string]interface{} {
	return genShardingStructure(numShard, shardID, BabylonHTTPPattern, BabylonWSPattern)
}

// IsSkippedEpoch returns if an epoch was skipped on shard due to staking epoch
func (ms babylonSchedule) IsSkippedEpoch(shardID uint32, epoch *big.Int) bool {
	return false
}

var babylonReshardingEpoch = []*big.Int{
	big.NewInt(0),
	params.BabylonChainConfig.StakingEpoch,
}

var (
	babylonV0 = MustNewInstance(2, 10, 10, 	numeric.OneDec(), 								 genesis.BabylonAccounts, genesis.BabylonFnAccounts, babylonReshardingEpoch, BabylonSchedule.BlocksPerEpoch())
	babylonV1 = MustNewInstance(2, 100, 10, numeric.MustNewDecFromStr("0.84"), genesis.BabylonAccounts, genesis.BabylonFnAccounts, babylonReshardingEpoch, BabylonSchedule.BlocksPerEpoch())
)
