package genesis

import (
	"strconv"
	"strings"
	"testing"

	"github.com/btcsuite/btcutil/bech32"
	ethCommon "github.com/ethereum/go-ethereum/common"
	"github.com/Timestopeofficial/bls/ffi/go/bls"
	"github.com/Timestopeofficial/feechain/internal/common"
)

func TestString(t *testing.T) {
	_ = BeaconAccountPriKey()
}

func TestCommitteeAccounts(test *testing.T) {
	testAccounts(test, FoundationalNodeAccounts)
	testAccounts(test, FeechainAccounts)
	testAccounts(test, TNFeechainAccounts)
	testAccounts(test, TNFoundationalAccounts)
	testAccounts(test, PangaeaAccounts)
}

func testAccounts(test *testing.T, accounts []DeployAccount) {
	index := 0
	for _, account := range accounts {
		accIndex, _ := strconv.Atoi(strings.Trim(account.Index, " "))
		if accIndex != index {
			test.Error("Account index", account.Index, "not in sequence")
		}
		index++

		_, _, err := bech32.Decode(account.Address)
		if err != nil {
			test.Error("Account address", account.Address, "is not valid:", err)
		}

		pubKey := bls.PublicKey{}
		err = pubKey.DeserializeHexStr(account.BLSPublicKey)
		if err != nil {
			test.Error("Account bls public key", account.BLSPublicKey, "is not valid:", err)
		}
	}
}

func testDeployAccounts(t *testing.T, accounts []DeployAccount) {
	indicesByAddress := make(map[ethCommon.Address][]int)
	indicesByKey := make(map[string][]int)
	for index, account := range accounts {
		if strings.TrimSpace(account.Index) != strconv.Itoa(index) {
			t.Errorf("account %+v at index %v has wrong index string",
				account, index)
		}
		if address, err := common.Bech32ToAddress(account.Address); err != nil {
			t.Errorf("account %+v at index %v has invalid address (%s)",
				account, index, err)
		} else {
			indicesByAddress[address] = append(indicesByAddress[address], index)
		}
		pubKey := bls.PublicKey{}
		if err := pubKey.DeserializeHexStr(account.BLSPublicKey); err != nil {
			t.Errorf("account %+v at index %v has invalid public key (%s)",
				account, index, err)
		} else {
			pubKeyStr := pubKey.SerializeToHexStr()
			indicesByKey[pubKeyStr] = append(indicesByKey[pubKeyStr], index)
		}
	}
	for address, indices := range indicesByAddress {
		if len(indices) > 1 {
			t.Errorf("account address %s appears in multiple rows: %v",
				common.MustAddressToBech32(address), indices)
		}
	}
	for pubKey, indices := range indicesByKey {
		if len(indices) > 1 {
			t.Errorf("BLS public key %s appears in multiple rows: %v",
				pubKey, indices)
		}
	}
}
