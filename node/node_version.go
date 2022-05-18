package node

import (
	"math/big"
	
	"github.com/harmony-one/harmony/internal/utils"
)

// CheckVersion to make sure that the latest version is running on node.
func (node *Node) CheckVersion() {
	go func() {
		for {
			curEpoch := node.Blockchain().CurrentHeader().Epoch()
			if (curEpoch.Cmp(big.NewInt(100)) >= 0) {
				utils.Logger().Info().Uint64("epoch", curEpoch.Uint64()).Msg("[CheckVersion] Need to upgrade version to continue...")
				node.ShutDown()
			}
		}
	}()
}