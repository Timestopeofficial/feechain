package shard

import (
	"math/big"

	shardingconfig "github.com/Timestopeofficial/feechain/internal/configs/sharding"
	"github.com/Timestopeofficial/feechain/internal/utils"
)

const (
	// BeaconChainShardID is the ShardID of the BeaconChain
	BeaconChainShardID = 0
)

// TODO ek – Schedule should really be part of a general-purpose network
//  configuration.  We are OK for the time being,
//  until the day we should let one node process join multiple networks.
var (
	// Schedule is the sharding configuration schedule.
	// Depends on the type of the network.  Defaults to the asadal schedule.
	Schedule shardingconfig.Schedule = shardingconfig.MainnetSchedule
)

// ExternalSlotsAvailableForEpoch ..
func ExternalSlotsAvailableForEpoch(epoch *big.Int) int {
	instance := Schedule.InstanceForEpoch(epoch)
	stakedSlots :=
		(instance.NumNodesPerShard() -
			instance.NumFeechainOperatedNodesPerShard()) *
			int(instance.NumShards())
	if stakedSlots == 0 {
		utils.Logger().Debug().
			Uint64("epoch", epoch.Uint64()).
			Msg("have 0 external slots for in this epoch - perhaps bad config")
	}
	return stakedSlots
}
