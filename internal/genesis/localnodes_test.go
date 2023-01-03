package genesis

import "testing"

func TestLocalTestAccounts(t *testing.T) {
	for name, accounts := range map[string][]DeployAccount{
		"FeechainV0":      LocalFeechainAccounts,
		"FeechainV1":      LocalFeechainAccountsV1,
		"FeechainV2":      LocalFeechainAccountsV2,
		"FoundationalV0": LocalFnAccounts,
		"FoundationalV1": LocalFnAccountsV1,
		"FoundationalV2": LocalFnAccountsV2,
	} {
		t.Run(name, func(t *testing.T) { testDeployAccounts(t, accounts) })
	}
}
