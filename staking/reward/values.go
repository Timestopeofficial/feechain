package reward

import (
	"fmt"
	"math/big"

	"github.com/Timestopeofficial/feechain/common/denominations"
	"github.com/Timestopeofficial/feechain/consensus/engine"
	shardingconfig "github.com/Timestopeofficial/feechain/internal/configs/sharding"
	"github.com/Timestopeofficial/feechain/internal/params"
	"github.com/Timestopeofficial/feechain/numeric"
	"github.com/Timestopeofficial/feechain/shard"
)

var (
	// PreStakedBlocks is the block reward, to be split evenly among block signers in pre-staking era.
	// 99 FEE per block
	PreStakedBlocks = new(big.Int).Mul(big.NewInt(99), big.NewInt(denominations.One))
	// StakedBlocks is the flat-rate block reward for epos staking launch.
	// 100 FEE per block.
	StakedBlocks = numeric.NewDecFromBigInt(new(big.Int).Mul(
		big.NewInt(100), big.NewInt(denominations.One),
	))
	// FiveSecStakedBlocks is the flat-rate block reward.
	// 99 FEE per block
	FiveSecStakedBlocks = numeric.NewDecFromBigInt(new(big.Int).Mul(
		big.NewInt(99*denominations.Nano), big.NewInt(denominations.Nano),
	))

	// TotalInitialTokens is the total amount of tokens (in ONE) at block 0 of the network.
	// This should be set/change on the node's init according to the core.GenesisSpec.
	TotalInitialTokens = numeric.Dec{Int: big.NewInt(0)}

	// None ..
	None = big.NewInt(0)

	// ErrInvalidBeaconChain if given chain is not beacon chain
	ErrInvalidBeaconChain = fmt.Errorf("given chain is not beaconchain")
)

// getPreStakingRewardsFromBlockNumber returns the number of tokens injected into the network
// in the pre-staking era (epoch < staking epoch) in ATTO.
//
// If the block number is > than the last block of an epoch, the last block of the epoch is
// used for the calculation by default.
//
// WARNING: This assumes beacon chain is at most the same block height as another shard in the
// transition from pre-staking to staking era/epoch.
func getPreStakingRewardsFromBlockNumber(id shardingconfig.NetworkID, blockNum *big.Int) *big.Int {
	if blockNum.Cmp(big.NewInt(2)) == -1 {
		// block 0 & 1 does not contain block rewards
		return big.NewInt(0)
	}

	lastBlockInEpoch := blockNum

	switch id {
	case shardingconfig.MainNet:
		lastBlockInEpoch = new(big.Int).SetUint64(shardingconfig.MainnetSchedule.EpochLastBlock(
			params.MainnetChainConfig.StakingEpoch.Uint64() - 1,
		))
	case shardingconfig.Babylon:
		lastBlockInEpoch = new(big.Int).SetUint64(shardingconfig.BabylonSchedule.EpochLastBlock(
			params.BabylonChainConfig.StakingEpoch.Uint64() - 1,
		))
	case shardingconfig.TestNet:
		lastBlockInEpoch = new(big.Int).SetUint64(shardingconfig.TestnetSchedule.EpochLastBlock(
			params.TestnetChainConfig.StakingEpoch.Uint64() - 1,
		))
	case shardingconfig.LocalNet:
		lastBlockInEpoch = new(big.Int).SetUint64(shardingconfig.LocalnetSchedule.EpochLastBlock(
			params.LocalnetChainConfig.StakingEpoch.Uint64() - 1,
		))
	}

	if blockNum.Cmp(lastBlockInEpoch) == 1 {
		blockNum = lastBlockInEpoch
	}

	return new(big.Int).Mul(PreStakedBlocks, new(big.Int).Sub(blockNum, big.NewInt(1)))
}

// WARNING: the data collected here are calculated from a consumer of the Rosetta API.
// If data becomes mission critical, implement a cross-link based approach.
//
// Data Source: https://github.com/Timestopeofficial/jupyter
//
// TODO (dm): use first crosslink of all shards to compute rewards on network instead of relying on constants.
var (
	totalPreStakingNetworkRewardsInAtto = map[shardingconfig.NetworkID][]*big.Int{
		shardingconfig.MainNet: {
			// Below are all of the last blocks of pre-staking era for asadal.
			getPreStakingRewardsFromBlockNumber(shardingconfig.MainNet, big.NewInt(28800)),
			getPreStakingRewardsFromBlockNumber(shardingconfig.MainNet, big.NewInt(28800)),
			getPreStakingRewardsFromBlockNumber(shardingconfig.MainNet, big.NewInt(28800)),
			getPreStakingRewardsFromBlockNumber(shardingconfig.MainNet, big.NewInt(28800)),
		},
		shardingconfig.Babylon: {
			// Below are all of the last blocks of pre-staking era for babylon.
			getPreStakingRewardsFromBlockNumber(shardingconfig.Babylon, big.NewInt(28800)),
			getPreStakingRewardsFromBlockNumber(shardingconfig.Babylon, big.NewInt(28800)),
			getPreStakingRewardsFromBlockNumber(shardingconfig.Babylon, big.NewInt(28800)),
			getPreStakingRewardsFromBlockNumber(shardingconfig.Babylon, big.NewInt(28800)),
		},
		shardingconfig.TestNet: {
			// Below are all of the placeholders 'last blocks' of pre-staking era for testnet.
			getPreStakingRewardsFromBlockNumber(shardingconfig.TestNet, big.NewInt(999999)),
			getPreStakingRewardsFromBlockNumber(shardingconfig.TestNet, big.NewInt(999999)),
			getPreStakingRewardsFromBlockNumber(shardingconfig.TestNet, big.NewInt(999999)),
			getPreStakingRewardsFromBlockNumber(shardingconfig.TestNet, big.NewInt(999999)),
		},
		shardingconfig.LocalNet: {
			// Below are all of the placeholders 'last blocks' of pre-staking era for localnet.
			getPreStakingRewardsFromBlockNumber(shardingconfig.LocalNet, big.NewInt(999999)),
			getPreStakingRewardsFromBlockNumber(shardingconfig.LocalNet, big.NewInt(999999)),
		},
	}
)

// getTotalPreStakingNetworkRewards in ATTO for given NetworkID
func getTotalPreStakingNetworkRewards(id shardingconfig.NetworkID) *big.Int {
	totalRewards := big.NewInt(0)
	if allRewards, ok := totalPreStakingNetworkRewardsInAtto[id]; ok {
		for _, reward := range allRewards {
			totalRewards = new(big.Int).Add(reward, totalRewards)
		}
	}
	return totalRewards
}

// GetTotalTokens in the network for all shards in ONE.
// This can only be computed with beaconchain if in staking era.
// If not in staking era, returns the rewards given out by the start of staking era.
func GetTotalTokens(chain engine.ChainReader) (numeric.Dec, error) {
	currHeader := chain.CurrentHeader()
	if !chain.Config().IsStaking(currHeader.Epoch()) {
		return GetTotalPreStakingTokens(), nil
	}
	if chain.ShardID() != shard.BeaconChainShardID {
		return numeric.Dec{}, ErrInvalidBeaconChain
	}

	stakingRewards, err := chain.ReadBlockRewardAccumulator(currHeader.Number().Uint64())
	if err != nil {
		return numeric.Dec{}, err
	}
	return GetTotalPreStakingTokens().Add(numeric.NewDecFromBigIntWithPrec(stakingRewards, 18)), nil
}

// GetTotalPreStakingTokens returns the total amount of tokens (in ONE) in the
// network at the the last block of the pre-staking era (epoch < staking epoch).
func GetTotalPreStakingTokens() numeric.Dec {
	preStakingRewards := numeric.NewDecFromBigIntWithPrec(
		getTotalPreStakingNetworkRewards(shard.Schedule.GetNetworkID()), 18,
	)
	return TotalInitialTokens.Add(preStakingRewards)
}

// SetTotalInitialTokens with the given initial tokens (from genesis in ATTO).
func SetTotalInitialTokens(initTokensAsAtto *big.Int) {
	TotalInitialTokens = numeric.NewDecFromBigIntWithPrec(initTokensAsAtto, 18)
}
