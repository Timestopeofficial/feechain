package shardingconfig

import (
	"math/big"

	"github.com/harmony-one/harmony/internal/params"

	"github.com/harmony-one/harmony/numeric"

	"github.com/harmony-one/harmony/internal/genesis"
)

const (
	mainnetEpochBlock1 = 8640
	blocksPerEpoch     = 8640

	mainnetVdfDifficulty = 50000 // This takes about 100s to finish the vdf

	// MainNetHTTPPattern is the http pattern for mainnet.
	MainNetHTTPPattern = "https://api.s%d.b.timestope.net"
	// MainNetWSPattern is the websocket pattern for mainnet.
	MainNetWSPattern = "wss://ws.s%d.b.timestope.net"
)

// MainnetSchedule is the mainnet sharding configuration schedule.
var MainnetSchedule mainnetSchedule

type mainnetSchedule struct{}

func (ms mainnetSchedule) InstanceForEpoch(epoch *big.Int) Instance {
	switch {
	case epoch.Cmp(params.MainnetChainConfig.StakingEpoch) >= 0:
		return mainnetV2
	default: // genesis
		return mainnetV1
	}
}

func (ms mainnetSchedule) BlocksPerEpoch() uint64 {
	return blocksPerEpoch
}

func (ms mainnetSchedule) CalcEpochNumber(blockNum uint64) *big.Int {
	switch {
	case blockNum >= mainnetEpochBlock1:
		return big.NewInt(int64((blockNum-mainnetEpochBlock1)/ms.BlocksPerEpoch()) + 1)
	default:
		return big.NewInt(0)
	}
}

func (ms mainnetSchedule) IsLastBlock(blockNum uint64) bool {
	switch {
	case blockNum < mainnetEpochBlock1-1:
		return false
	case blockNum == mainnetEpochBlock1-1:
		return true
	default:
		return ((blockNum-mainnetEpochBlock1)%ms.BlocksPerEpoch() == ms.BlocksPerEpoch()-1)
	}
}

func (ms mainnetSchedule) EpochLastBlock(epochNum uint64) uint64 {
	switch {
	case epochNum == 0:
		return mainnetEpochBlock1 - 1
	default:
		return mainnetEpochBlock1 - 1 + ms.BlocksPerEpoch()*epochNum
	}
}

func (ms mainnetSchedule) VdfDifficulty() int {
	return mainnetVdfDifficulty
}

func (ms mainnetSchedule) GetNetworkID() NetworkID {
	return MainNet
}

// GetShardingStructure is the sharding structure for mainnet.
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
	mainnetV1 = MustNewInstance(4, 4, 3, numeric.MustNewDecFromStr("0.68"), genesis.HarmonyAccounts, genesis.FoundationalNodeAccountsV1_5, mainnetReshardingEpoch, MainnetSchedule.BlocksPerEpoch())
	mainnetV2 = MustNewInstance(4, 70, 3, numeric.MustNewDecFromStr("0.68"), genesis.HarmonyAccounts, genesis.FoundationalNodeAccountsV1_5, mainnetReshardingEpoch, MainnetSchedule.BlocksPerEpoch())
)
