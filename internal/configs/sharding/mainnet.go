package shardingconfig

import (
	"math/big"

	"github.com/Timestopeofficial/feechain/internal/params"

	"github.com/Timestopeofficial/feechain/numeric"

	"github.com/Timestopeofficial/feechain/internal/genesis"
)

const (
	blocksPerEpoch     = 28800

	mainnetVdfDifficulty = 50000 // This takes about 100s to finish the vdf

	// MainNetHTTPPattern is the http pattern for asadal.
	MainNetHTTPPattern = "https://api.s%d.asadal.timestope.net"
	// MainNetWSPattern is the websocket pattern for asadal.
	MainNetWSPattern = "wss://ws.s%d.asadal.timestope.net"
)

// MainnetSchedule is the asadal sharding configuration schedule.
var MainnetSchedule mainnetSchedule

type mainnetSchedule struct{}

func (ms mainnetSchedule) InstanceForEpoch(epoch *big.Int) Instance {
	switch {
	case epoch.Cmp(params.MainnetChainConfig.StakingEpoch) >= 0:
		return mainnetV1
	default: // genesis
		return mainnetV0
	}
}

func (ms mainnetSchedule) BlocksPerEpoch() uint64 {
	return blocksPerEpoch
}

func (ms mainnetSchedule) CalcEpochNumber(blockNum uint64) *big.Int {
	return big.NewInt(int64(blockNum / ms.BlocksPerEpoch()))
}

func (ms mainnetSchedule) IsLastBlock(blockNum uint64) bool {
	return (blockNum % ms.BlocksPerEpoch() == ms.BlocksPerEpoch() - 1)
}

func (ms mainnetSchedule) EpochLastBlock(epochNum uint64) uint64 {
		return ms.BlocksPerEpoch() * (epochNum + 1) - 1
}

func (ms mainnetSchedule) VdfDifficulty() int {
	return mainnetVdfDifficulty
}

func (ms mainnetSchedule) GetNetworkID() NetworkID {
	return MainNet
}

// GetShardingStructure is the sharding structure for asadal.
func (ms mainnetSchedule) GetShardingStructure(numShard, shardID int) []map[string]interface{} {
	return genShardingStructure(numShard, shardID, MainNetHTTPPattern, MainNetWSPattern)
}

// IsSkippedEpoch returns if an epoch was skipped on shard due to staking epoch
func (ms mainnetSchedule) IsSkippedEpoch(shardID uint32, epoch *big.Int) bool {
	return false
}

var mainnetReshardingEpoch = []*big.Int{
	big.NewInt(0),
	params.MainnetChainConfig.StakingEpoch,
}

var (
	mainnetV0 = MustNewInstance(4, 6, 6, 	numeric.OneDec(), 								 genesis.HarmonyAccounts, genesis.FoundationalNodeAccounts, mainnetReshardingEpoch, MainnetSchedule.BlocksPerEpoch())
	mainnetV1 = MustNewInstance(4, 85, 6, numeric.MustNewDecFromStr("0.84"), genesis.HarmonyAccounts, genesis.FoundationalNodeAccounts, mainnetReshardingEpoch, MainnetSchedule.BlocksPerEpoch())
)
