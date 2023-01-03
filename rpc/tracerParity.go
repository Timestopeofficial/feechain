package rpc

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/Timestopeofficial/feechain/eth/rpc"
	"github.com/Timestopeofficial/feechain/fch"
)

var (
	parityTraceGO = "ParityBlockTracer"
)

type PublicParityTracerService struct {
	*PublicTracerService
}

func (s *PublicParityTracerService) Transaction(ctx context.Context, hash common.Hash) (interface{}, error) {
	timer := DoMetricRPCRequest(Transaction)
	defer DoRPCRequestDuration(Transaction, timer)
	return s.TraceTransaction(ctx, hash, &fch.TraceConfig{Tracer: &parityTraceGO})
}

// trace_block RPC
func (s *PublicParityTracerService) Block(ctx context.Context, number rpc.BlockNumber) (interface{}, error) {
	timer := DoMetricRPCRequest(Block)
	defer DoRPCRequestDuration(Block, timer)

	block := s.fch.BlockChain.GetBlockByNumber(uint64(number))
	if block == nil {
		return nil, nil
	}
	results, err := s.fch.TraceBlock(ctx, block, &fch.TraceConfig{Tracer: &parityTraceGO})
	if err != nil {
		return results, err
	}
	var resultArray = make([]json.RawMessage, 0)
	for _, result := range results {
		raw, ok := result.Result.([]json.RawMessage)
		if !ok {
			return results, errors.New("tracer bug:expected []json.RawMessage")
		}
		resultArray = append(resultArray, raw...)
	}
	return resultArray, nil
}
