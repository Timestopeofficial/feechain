package fch

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/bloombits"
	"github.com/ethereum/go-ethereum/event"
	"github.com/Timestopeofficial/feechain/block"
	"github.com/Timestopeofficial/feechain/core"
	"github.com/Timestopeofficial/feechain/core/rawdb"
	"github.com/Timestopeofficial/feechain/core/state"
	"github.com/Timestopeofficial/feechain/core/types"
	"github.com/Timestopeofficial/feechain/crypto/bls"
	internal_bls "github.com/Timestopeofficial/feechain/crypto/bls"
	"github.com/Timestopeofficial/feechain/eth/rpc"
	internal_common "github.com/Timestopeofficial/feechain/internal/common"
	"github.com/Timestopeofficial/feechain/internal/params"
	"github.com/Timestopeofficial/feechain/internal/utils"
	"github.com/Timestopeofficial/feechain/shard"
	"github.com/Timestopeofficial/feechain/staking/availability"
	stakingReward "github.com/Timestopeofficial/feechain/staking/reward"
	"github.com/pkg/errors"
)

// ChainConfig ...
func (fch *Feechain) ChainConfig() *params.ChainConfig {
	return fch.BlockChain.Config()
}

// GetShardState ...
func (fch *Feechain) GetShardState() (*shard.State, error) {
	return fch.BlockChain.ReadShardState(fch.BlockChain.CurrentHeader().Epoch())
}

// GetBlockSigners ..
func (fch *Feechain) GetBlockSigners(
	ctx context.Context, blockNum rpc.BlockNumber,
) (shard.SlotList, *internal_bls.Mask, error) {
	blk, err := fch.BlockByNumber(ctx, blockNum)
	if err != nil {
		return nil, nil, err
	}
	blockWithSigners, err := fch.BlockByNumber(ctx, blockNum+1)
	if err != nil {
		return nil, nil, err
	}
	if blockWithSigners == nil {
		return nil, nil, fmt.Errorf("block number %v not found", blockNum+1)
	}
	committee, err := fch.GetValidators(blk.Epoch())
	if err != nil {
		return nil, nil, err
	}
	pubKeys := make([]internal_bls.PublicKeyWrapper, len(committee.Slots))
	for i, validator := range committee.Slots {
		key, err := bls.BytesToBLSPublicKey(validator.BLSPublicKey[:])
		if err != nil {
			return nil, nil, err
		}
		pubKeys[i] = internal_bls.PublicKeyWrapper{
			Bytes:  validator.BLSPublicKey,
			Object: key,
		}
	}
	mask, err := internal_bls.NewMask(pubKeys, nil)
	if err != nil {
		return nil, nil, err
	}
	err = mask.SetMask(blockWithSigners.Header().LastCommitBitmap())
	if err != nil {
		return nil, nil, err
	}
	return committee.Slots, mask, nil
}

// DetailedBlockSignerInfo contains all of the block singing information
type DetailedBlockSignerInfo struct {
	// Signers are all the signers for the block
	Signers shard.SlotList
	// Committee when the block was signed.
	Committee shard.SlotList
	BlockHash common.Hash
}

// GetDetailedBlockSignerInfo fetches the block signer information for any non-genesis block
func (fch *Feechain) GetDetailedBlockSignerInfo(
	ctx context.Context, blk *types.Block,
) (*DetailedBlockSignerInfo, error) {
	parentBlk, err := fch.BlockByNumber(ctx, rpc.BlockNumber(blk.NumberU64()-1))
	if err != nil {
		return nil, err
	}
	parentShardState, err := fch.BlockChain.ReadShardState(parentBlk.Epoch())
	if err != nil {
		return nil, err
	}
	committee, signers, _, err := availability.BallotResult(
		parentBlk.Header(), blk.Header(), parentShardState, blk.ShardID(),
	)
	return &DetailedBlockSignerInfo{
		Signers:   signers,
		Committee: committee,
		BlockHash: blk.Hash(),
	}, nil
}

// PreStakingBlockRewards are the rewards for a block in the pre-staking era (epoch < staking epoch).
type PreStakingBlockRewards map[common.Address]*big.Int

// GetPreStakingBlockRewards for the given block number.
// Calculated rewards are done exactly like chain.AccumulateRewardsAndCountSigs.
func (fch *Feechain) GetPreStakingBlockRewards(
	ctx context.Context, blk *types.Block,
) (PreStakingBlockRewards, error) {
	if fch.IsStakingEpoch(blk.Epoch()) {
		return nil, fmt.Errorf("block %v is in staking era", blk.Number())
	}

	if cachedReward, ok := fch.preStakingBlockRewardsCache.Get(blk.Hash()); ok {
		return cachedReward.(PreStakingBlockRewards), nil
	}
	rewards := PreStakingBlockRewards{}

	sigInfo, err := fch.GetDetailedBlockSignerInfo(ctx, blk)
	if err != nil {
		return nil, err
	}
	last := big.NewInt(0)
	count := big.NewInt(int64(len(sigInfo.Signers)))
	for i, slot := range sigInfo.Signers {
		rewardsForThisAddr, ok := rewards[slot.EcdsaAddress]
		if !ok {
			rewardsForThisAddr = big.NewInt(0)
		}
		cur := big.NewInt(0)
		cur.Mul(stakingReward.PreStakedBlocks, big.NewInt(int64(i+1))).Div(cur, count)
		reward := big.NewInt(0).Sub(cur, last)
		rewards[slot.EcdsaAddress] = new(big.Int).Add(reward, rewardsForThisAddr)
		last = cur
	}

	// Report tx fees of the coinbase (== leader)
	receipts, err := fch.GetReceipts(ctx, blk.Hash())
	if err != nil {
		return nil, err
	}
	txFees := big.NewInt(0)
	for _, tx := range blk.Transactions() {
		txnHash := tx.HashByType()
		dbTx, _, _, receiptIndex := rawdb.ReadTransaction(fch.ChainDb(), txnHash)
		if dbTx == nil {
			return nil, fmt.Errorf("could not find receipt for tx: %v", txnHash.String())
		}
		if len(receipts) <= int(receiptIndex) {
			return nil, fmt.Errorf("invalid receipt indext %v (>= num receipts: %v) for tx: %v",
				receiptIndex, len(receipts), txnHash.String())
		}
		txFee := new(big.Int).Mul(tx.GasPrice(), big.NewInt(int64(receipts[receiptIndex].GasUsed)))
		txFees = new(big.Int).Add(txFee, txFees)
	}

	if amt, ok := rewards[blk.Header().Coinbase()]; ok {
		rewards[blk.Header().Coinbase()] = new(big.Int).Add(amt, txFees)
	} else {
		rewards[blk.Header().Coinbase()] = txFees
	}

	fch.preStakingBlockRewardsCache.Add(blk.Hash(), rewards)
	return rewards, nil
}

// GetLatestChainHeaders ..
func (fch *Feechain) GetLatestChainHeaders() *block.HeaderPair {
	return &block.HeaderPair{
		BeaconHeader: fch.BeaconChain.CurrentHeader(),
		ShardHeader:  fch.BlockChain.CurrentHeader(),
	}
}

// GetLastCrossLinks ..
func (fch *Feechain) GetLastCrossLinks() ([]*types.CrossLink, error) {
	crossLinks := []*types.CrossLink{}
	for i := uint32(1); i < shard.Schedule.InstanceForEpoch(fch.CurrentBlock().Epoch()).NumShards(); i++ {
		link, err := fch.BlockChain.ReadShardLastCrossLink(i)
		if err != nil {
			return nil, err
		}
		crossLinks = append(crossLinks, link)
	}

	return crossLinks, nil
}

// CurrentBlock ...
func (fch *Feechain) CurrentBlock() *types.Block {
	return types.NewBlockWithHeader(fch.BlockChain.CurrentHeader())
}

// GetBlock ...
func (fch *Feechain) GetBlock(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return fch.BlockChain.GetBlockByHash(hash), nil
}

// GetCurrentBadBlocks ..
func (fch *Feechain) GetCurrentBadBlocks() []core.BadBlock {
	return fch.BlockChain.BadBlocks()
}

// GetBalance returns balance of an given address.
func (fch *Feechain) GetBalance(ctx context.Context, address common.Address, blockNum rpc.BlockNumber) (*big.Int, error) {
	s, _, err := fch.StateAndHeaderByNumber(ctx, blockNum)
	if s == nil || err != nil {
		return nil, err
	}
	return s.GetBalance(address), s.Error()
}

// BlockByNumber ...
func (fch *Feechain) BlockByNumber(ctx context.Context, blockNum rpc.BlockNumber) (*types.Block, error) {
	// Pending block is only known by the miner
	if blockNum == rpc.PendingBlockNumber {
		return nil, errors.New("not implemented")
	}
	// Otherwise resolve and return the block
	if blockNum == rpc.LatestBlockNumber {
		return fch.BlockChain.CurrentBlock(), nil
	}
	return fch.BlockChain.GetBlockByNumber(uint64(blockNum)), nil
}

// HeaderByNumber ...
func (fch *Feechain) HeaderByNumber(ctx context.Context, blockNum rpc.BlockNumber) (*block.Header, error) {
	// Pending block is only known by the miner
	if blockNum == rpc.PendingBlockNumber {
		return nil, errors.New("not implemented")
	}
	// Otherwise resolve and return the block
	if blockNum == rpc.LatestBlockNumber {
		return fch.BlockChain.CurrentBlock().Header(), nil
	}
	return fch.BlockChain.GetHeaderByNumber(uint64(blockNum)), nil
}

// HeaderByHash ...
func (fch *Feechain) HeaderByHash(ctx context.Context, blockHash common.Hash) (*block.Header, error) {
	header := fch.BlockChain.GetHeaderByHash(blockHash)
	if header == nil {
		return nil, errors.New("Header is not found")
	}
	return header, nil
}

// StateAndHeaderByNumber ...
func (fch *Feechain) StateAndHeaderByNumber(ctx context.Context, blockNum rpc.BlockNumber) (*state.DB, *block.Header, error) {
	// Pending state is only known by the miner
	if blockNum == rpc.PendingBlockNumber {
		return nil, nil, errors.New("not implemented")
	}
	// Otherwise resolve the block number and return its state
	header, err := fch.HeaderByNumber(ctx, blockNum)
	if header == nil || err != nil {
		return nil, nil, err
	}
	stateDb, err := fch.BlockChain.StateAt(header.Root())
	return stateDb, header, err
}

// GetLeaderAddress returns the one address of the leader, given the coinbaseAddr.
// Note that the coinbaseAddr is overloaded with the BLS pub key hash in staking era.
func (fch *Feechain) GetLeaderAddress(coinbaseAddr common.Address, epoch *big.Int) string {
	if fch.IsStakingEpoch(epoch) {
		if leader, exists := fch.leaderCache.Get(coinbaseAddr); exists {
			bech32, _ := internal_common.AddressToBech32(leader.(common.Address))
			return bech32
		}
		committee, err := fch.GetValidators(epoch)
		if err != nil {
			return ""
		}
		for _, val := range committee.Slots {
			addr := utils.GetAddressFromBLSPubKeyBytes(val.BLSPublicKey[:])
			fch.leaderCache.Add(addr, val.EcdsaAddress)
			if addr == coinbaseAddr {
				bech32, _ := internal_common.AddressToBech32(val.EcdsaAddress)
				return bech32
			}
		}
		return "" // Did not find matching address
	}
	bech32, _ := internal_common.AddressToBech32(coinbaseAddr)
	return bech32
}

// Filter related APIs

// GetLogs ...
func (fch *Feechain) GetLogs(ctx context.Context, blockHash common.Hash, isEth bool) ([][]*types.Log, error) {
	receipts := fch.BlockChain.GetReceiptsByHash(blockHash)
	if receipts == nil {
		return nil, errors.New("Missing receipts")
	}
	if isEth {
		block := fch.BlockChain.GetBlockByHash(blockHash)
		if block == nil {
			return nil, errors.New("Missing block data")
		}
		txns := block.Transactions()
		for i, _ := range receipts {
			if i < len(txns) {
				ethHash := txns[i].ConvertToEth().Hash()
				receipts[i].TxHash = ethHash
				for j, _ := range receipts[i].Logs {
					// Override log txHash with receipt's
					receipts[i].Logs[j].TxHash = ethHash
				}
			}
		}
	}

	logs := make([][]*types.Log, len(receipts))
	for i, receipt := range receipts {
		logs[i] = receipt.Logs
	}
	return logs, nil
}

// ServiceFilter ...
func (fch *Feechain) ServiceFilter(ctx context.Context, session *bloombits.MatcherSession) {
	// TODO(dm): implement
}

// SubscribeNewTxsEvent subscribes new tx event.
// TODO: this is not implemented or verified yet for feechain.
func (fch *Feechain) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	return fch.TxPool.SubscribeNewTxsEvent(ch)
}

// SubscribeChainEvent subscribes chain event.
// TODO: this is not implemented or verified yet for feechain.
func (fch *Feechain) SubscribeChainEvent(ch chan<- core.ChainEvent) event.Subscription {
	return fch.BlockChain.SubscribeChainEvent(ch)
}

// SubscribeChainHeadEvent subcribes chain head event.
// TODO: this is not implemented or verified yet for feechain.
func (fch *Feechain) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return fch.BlockChain.SubscribeChainHeadEvent(ch)
}

// SubscribeChainSideEvent subcribes chain side event.
// TODO: this is not implemented or verified yet for feechain.
func (fch *Feechain) SubscribeChainSideEvent(ch chan<- core.ChainSideEvent) event.Subscription {
	return fch.BlockChain.SubscribeChainSideEvent(ch)
}

// SubscribeRemovedLogsEvent subcribes removed logs event.
// TODO: this is not implemented or verified yet for feechain.
func (fch *Feechain) SubscribeRemovedLogsEvent(ch chan<- core.RemovedLogsEvent) event.Subscription {
	return fch.BlockChain.SubscribeRemovedLogsEvent(ch)
}

// SubscribeLogsEvent subcribes log event.
// TODO: this is not implemented or verified yet for feechain.
func (fch *Feechain) SubscribeLogsEvent(ch chan<- []*types.Log) event.Subscription {
	return fch.BlockChain.SubscribeLogsEvent(ch)
}
