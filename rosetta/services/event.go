package services

import (
	"context"

	"github.com/coinbase/rosetta-sdk-go/types"
	fchTypes "github.com/Timestopeofficial/feechain/core/types"
	"github.com/Timestopeofficial/feechain/fch"
)

// EventAPI implements the server.EventsAPIServicer interface.
type EventAPI struct {
	fch *fch.Feechain
}

func NewEventAPI(fch *fch.Feechain) *EventAPI {
	return &EventAPI{fch: fch}
}

// EventsBlocks implements the /events/blocks endpoint
func (e *EventAPI) EventsBlocks(ctx context.Context, request *types.EventsBlocksRequest) (resp *types.EventsBlocksResponse, err *types.Error) {
	cacheItem, cacheHelper, cacheErr := rosettaCacheHelper("EventsBlocks", request)
	if cacheErr == nil {
		if cacheItem != nil {
			return cacheItem.resp.(*types.EventsBlocksResponse), nil
		} else {
			defer cacheHelper(resp, err)
		}
	}

	if err := assertValidNetworkIdentifier(request.NetworkIdentifier, e.fch.ShardID); err != nil {
		return nil, err
	}

	var offset, limit int64

	if request.Limit == nil {
		limit = 10
	} else {
		limit = *request.Limit
		if limit > 1000 {
			limit = 1000
		}
	}

	if request.Offset == nil {
		offset = 0
	} else {
		offset = *request.Offset
	}

	resp = &types.EventsBlocksResponse{
		MaxSequence: e.fch.BlockChain.CurrentHeader().Number().Int64(),
	}

	for i := offset; i < offset+limit; i++ {
		block := e.fch.BlockChain.GetBlockByNumber(uint64(i))
		if block == nil {
			break
		}

		resp.Events = append(resp.Events, buildFromBlock(block))
	}

	return resp, nil
}

func buildFromBlock(block *fchTypes.Block) *types.BlockEvent {
	return &types.BlockEvent{
		Sequence: block.Number().Int64(),
		BlockIdentifier: &types.BlockIdentifier{
			Index: block.Number().Int64(),
			Hash:  block.Hash().Hex(),
		},
		Type: types.ADDED,
	}
}
