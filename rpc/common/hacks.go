package common

import (
	"github.com/Timestopeofficial/feechain/consensus/quorum"
	"github.com/Timestopeofficial/feechain/crypto/bls"
	"github.com/Timestopeofficial/feechain/numeric"
)

type setRawStakeHack interface {
	SetRawStake(key bls.SerializedPublicKey, d numeric.Dec)
}

// SetRawStake is a hack, return value is if was successful or not at setting
func SetRawStake(q quorum.Decider, key bls.SerializedPublicKey, d numeric.Dec) bool {
	if setter, ok := q.(setRawStakeHack); ok {
		setter.SetRawStake(key, d)
		return true
	}
	return false
}
