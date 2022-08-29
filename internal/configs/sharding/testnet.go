package shardingconfig

import (
	"math/big"

	"github.com/harmony-one/harmony/internal/genesis"
	"github.com/harmony-one/harmony/internal/params"
	"github.com/harmony-one/harmony/numeric"
)

// TestnetSchedule is the long-running public testnet sharding
// configuration schedule.
var TestnetSchedule testnetSchedule

type testnetSchedule struct{}

const (
	// ~304 sec epochs for P2 of open staking
	testnetBlocksPerEpoch = 4500

	testnetVdfDifficulty = 10000 // This takes about 20s to finish the vdf

	// TestNetHTTPPattern is the http pattern for testnet.
	TestNetHTTPPattern = "https://api.s%d.t.timestope.net"
	// TestNetWSPattern is the websocket pattern for testnet.
	TestNetWSPattern = "wss://ws.s%d.t.timestope.net"
)

func (ts testnetSchedule) InstanceForEpoch(epoch *big.Int) Instance {
	switch {
	case epoch.Cmp(params.TestnetChainConfig.StakingEpoch) >= 0:
		return testnetV1
	default: // genesis
		return testnetV0
	}
}

func (ts testnetSchedule) BlocksPerEpoch() uint64 {
	return testnetBlocksPerEpoch
}

func (ts testnetSchedule) CalcEpochNumber(blockNum uint64) *big.Int {
	return big.NewInt(int64(blockNum / ts.BlocksPerEpoch()))
}

func (ts testnetSchedule) IsLastBlock(blockNum uint64) bool {
	return (blockNum+1)%ts.BlocksPerEpoch() == 0
}

func (ts testnetSchedule) EpochLastBlock(epochNum uint64) uint64 {
	return ts.BlocksPerEpoch()*(epochNum+1) - 1
}

func (ts testnetSchedule) VdfDifficulty() int {
	return testnetVdfDifficulty
}

func (ts testnetSchedule) GetNetworkID() NetworkID {
	return TestNet
}

// GetShardingStructure is the sharding structure for testnet.
func (ts testnetSchedule) GetShardingStructure(numShard, shardID int) []map[string]interface{} {
	return genShardingStructure(numShard, shardID, TestNetHTTPPattern, TestNetWSPattern)
}

// IsSkippedEpoch returns if an epoch was skipped on shard due to staking epoch
func (ts testnetSchedule) IsSkippedEpoch(shardID uint32, epoch *big.Int) bool {
	return false
}

var testnetReshardingEpoch = []*big.Int{
	big.NewInt(0),
	params.TestnetChainConfig.StakingEpoch,
}

var testnetV0 = MustNewInstance(2, 4, 4, numeric.OneDec(), genesis.TNHarmonyAccounts, genesis.TNFoundationalAccounts, testnetReshardingEpoch, TestnetSchedule.BlocksPerEpoch())
var testnetV1 = MustNewInstance(2, 20, 4, numeric.MustNewDecFromStr("0.80"), genesis.TNHarmonyAccounts, genesis.TNFoundationalAccounts, testnetReshardingEpoch, TestnetSchedule.BlocksPerEpoch())
