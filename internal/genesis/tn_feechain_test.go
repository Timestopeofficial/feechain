package genesis

import "testing"

func TestTNFeechainAccounts(t *testing.T) {
	testDeployAccounts(t, TNFeechainAccounts)
}

func TestTNFoundationalAccounts(t *testing.T) {
	testDeployAccounts(t, TNFoundationalAccounts)
}
