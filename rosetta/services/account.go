package services

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Timestopeofficial/feechain/core/vm"

	"github.com/coinbase/rosetta-sdk-go/server"
	"github.com/coinbase/rosetta-sdk-go/types"
	ethCommon "github.com/ethereum/go-ethereum/common"

	fchTypes "github.com/Timestopeofficial/feechain/core/types"
	"github.com/Timestopeofficial/feechain/eth/rpc"
	"github.com/Timestopeofficial/feechain/fch"
	internalCommon "github.com/Timestopeofficial/feechain/internal/common"
	"github.com/Timestopeofficial/feechain/rosetta/common"
)

// AccountAPI implements the server.AccountAPIServicer interface.
type AccountAPI struct {
	fch *fch.Feechain
}

func (s *AccountAPI) AccountCoins(ctx context.Context, request *types.AccountCoinsRequest) (*types.AccountCoinsResponse, *types.Error) {
	panic("implement me")
}

// NewAccountAPI creates a new instance of a BlockAPI.
func NewAccountAPI(fch *fch.Feechain) server.AccountAPIServicer {
	return &AccountAPI{
		fch: fch,
	}
}

// AccountBalance implements the /account/balance endpoint
func (s *AccountAPI) AccountBalance(
	ctx context.Context, request *types.AccountBalanceRequest,
) (*types.AccountBalanceResponse, *types.Error) {
	if err := assertValidNetworkIdentifier(request.NetworkIdentifier, s.fch.ShardID); err != nil {
		return nil, err
	}

	var block *fchTypes.Block
	var rosettaError *types.Error
	if request.BlockIdentifier == nil {
		block = s.fch.CurrentBlock()
	} else {
		block, rosettaError = getBlock(ctx, s.fch, request.BlockIdentifier)
		if rosettaError != nil {
			return nil, rosettaError
		}
	}

	addr, err := getAddress(request.AccountIdentifier)
	if err != nil {
		return nil, common.NewError(common.SanityCheckError, map[string]interface{}{
			"message": err.Error(),
		})
	}
	blockNum := rpc.BlockNumber(block.Header().Header.Number().Int64())
	balance := new(big.Int)

	if request.AccountIdentifier.SubAccount != nil {
		// indicate it may be a request for delegated balance
		balance, rosettaError = s.getStakingBalance(request.AccountIdentifier.SubAccount, addr, block)
		if rosettaError != nil {
			return nil, rosettaError
		}
	} else {
		balance, err = s.fch.GetBalance(ctx, addr, blockNum)
		if err != nil {
			return nil, common.NewError(common.SanityCheckError, map[string]interface{}{
				"message": "invalid address",
			})
		}
	}

	amount := types.Amount{
		Value:    balance.String(),
		Currency: &common.NativeCurrency,
	}

	respBlock := types.BlockIdentifier{
		Index: blockNum.Int64(),
		Hash:  block.Header().Hash().String(),
	}

	return &types.AccountBalanceResponse{
		BlockIdentifier: &respBlock,
		Balances:        []*types.Amount{&amount},
	}, nil
}

// getStakingBalance used for get delegated balance with sub account identifier
func (s *AccountAPI) getStakingBalance(
	subAccount *types.SubAccountIdentifier, addr ethCommon.Address, block *fchTypes.Block,
) (*big.Int, *types.Error) {
	balance := new(big.Int)
	ty, exist := subAccount.Metadata["type"]

	if !exist {
		return nil, common.NewError(common.SanityCheckError, map[string]interface{}{
			"message": "invalid sub account",
		})
	}

	switch ty.(string) {
	case Delegation:
		validatorAddr := subAccount.Address
		validators, delegations := s.fch.GetDelegationsByDelegatorByBlock(addr, block)
		for index, validator := range validators {
			if validatorAddr == internalCommon.MustAddressToBech32(validator) {
				balance = new(big.Int).Add(balance, delegations[index].Amount)
			}
		}
	case UnDelegation:
		validatorAddr := subAccount.Address
		validators, delegations := s.fch.GetDelegationsByDelegatorByBlock(addr, block)
		for index, validator := range validators {
			if validatorAddr == internalCommon.MustAddressToBech32(validator) {
				undelegations := delegations[index].Undelegations
				for _, undelegate := range undelegations {
					balance = new(big.Int).Add(balance, undelegate.Amount)
				}
			}
		}
	default:
		return nil, common.NewError(common.SanityCheckError, map[string]interface{}{
			"message": "invalid sub account type",
		})
	}

	return balance, nil
}

// AccountMetadata used for account identifiers
type AccountMetadata struct {
	Address string `json:"hex_address"`
}

// newAccountIdentifier ..
func newAccountIdentifier(
	address ethCommon.Address,
) (*types.AccountIdentifier, *types.Error) {
	b32Address, err := internalCommon.AddressToBech32(address)
	if err != nil {
		return nil, common.NewError(common.SanityCheckError, map[string]interface{}{
			"message": err.Error(),
		})
	}
	metadata, err := types.MarshalMap(AccountMetadata{Address: address.String()})
	if err != nil {
		return nil, common.NewError(common.CatchAllError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	return &types.AccountIdentifier{
		Address:  b32Address,
		Metadata: metadata,
	}, nil
}

// newAccountIdentifier ..
func newRosettaAccountIdentifier(address *vm.RosettaLogAddressItem) (*types.AccountIdentifier, *types.Error) {
	if address == nil || address.Account == nil {
		return nil, nil
	}

	b32Address, err := internalCommon.AddressToBech32(*address.Account)
	if err != nil {
		return nil, common.NewError(common.SanityCheckError, map[string]interface{}{
			"message": err.Error(),
		})
	}
	metadata, err := types.MarshalMap(AccountMetadata{Address: address.Account.String()})
	if err != nil {
		return nil, common.NewError(common.CatchAllError, map[string]interface{}{
			"message": err.Error(),
		})
	}

	ai := &types.AccountIdentifier{
		Address:  b32Address,
		Metadata: metadata,
	}

	if address.SubAccount != nil {
		b32Address, err := internalCommon.AddressToBech32(*address.SubAccount)
		if err != nil {
			return nil, common.NewError(common.SanityCheckError, map[string]interface{}{
				"message": err.Error(),
			})
		}

		ai.SubAccount = &types.SubAccountIdentifier{
			Address:  b32Address,
			Metadata: address.Metadata,
		}
	}

	return ai, nil
}

func newSubAccountIdentifier(
	address ethCommon.Address, metadata map[string]interface{},
) (*types.SubAccountIdentifier, *types.Error) {
	b32Address, err := internalCommon.AddressToBech32(address)
	if err != nil {
		return nil, common.NewError(common.SanityCheckError, map[string]interface{}{
			"message": err.Error(),
		})
	}
	return &types.SubAccountIdentifier{
		Address:  b32Address,
		Metadata: metadata,
	}, nil
}

func newAccountIdentifierWithSubAccount(
	address, subAddress ethCommon.Address, metadata map[string]interface{},
) (*types.AccountIdentifier, *types.Error) {
	accountIdentifier, err := newAccountIdentifier(address)
	if err != nil {
		return nil, err
	}

	subAccountIdentifier, err := newSubAccountIdentifier(subAddress, metadata)
	if err != nil {
		return nil, err
	}

	accountIdentifier.SubAccount = subAccountIdentifier
	return accountIdentifier, nil
}

// getAddress ..
func getAddress(
	identifier *types.AccountIdentifier,
) (ethCommon.Address, error) {
	if identifier == nil {
		return ethCommon.Address{}, fmt.Errorf("identifier cannot be nil")
	}
	return internalCommon.Bech32ToAddress(identifier.Address)
}
