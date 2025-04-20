package polygon_test

import (
	"crypto/ecdsa"
	"math/big"
	"os"
	"testing"

	"bridgebot/internal/blockchain/polygon"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/stretchr/testify/assert"
)

func setupTestClient(t *testing.T) *ethclient.Client {
	t.Setenv("POLYGON_RPC_URL", "https://polygon-rpc.com")
	client, err := polygon.NewPolygonClient()
	assert.NoError(t, err, "Failed to create Polygon client")
	return client
}

func getTestPrivateKey(t *testing.T) *ecdsa.PrivateKey {
	key := os.Getenv("PRIVATE_KEY")
	if key == "" {
		key = "dde619e9c94141eb5c60cf3c52e812f95db0a593543767a59e6b12e133a40c6d" // ⚠️ For testing only
	}
	pk, err := crypto.HexToECDSA(key)
	assert.NoError(t, err, "Invalid private key")
	return pk
}

func TestIsApprovalNeeded(t *testing.T) {
	client := setupTestClient(t)
	token := common.HexToAddress(polygon.TokenAddress)
	user := common.HexToAddress("0x7d0F13148e85A53227c65Ed013E7961A67839858")
	amount := big.NewInt(1000000)

	needed, err := polygon.IsApprovalNeeded(client, token, user, user, amount)
	assert.NoError(t, err, "IsApprovalNeeded returned error")
	t.Logf("Approval needed? %v", needed)
}

func TestApproveContract(t *testing.T) {
	client := setupTestClient(t)
	pk := getTestPrivateKey(t)
	token := common.HexToAddress(polygon.TokenAddress)
	spender := crypto.PubkeyToAddress(pk.PublicKey)
	amount := big.NewInt(5000000)

	txHash, err := polygon.ApproveContract(client, token, spender, amount, pk)
	assert.NoError(t, err, "ApproveContract failed")
	assert.NotEmpty(t, txHash, "Transaction hash should not be empty")
	t.Logf("ApproveContract tx hash: %s", txHash)
}

func TestRevokeApproval(t *testing.T) {
	client := setupTestClient(t)
	pk := getTestPrivateKey(t)
	spender := crypto.PubkeyToAddress(pk.PublicKey)
	token := common.HexToAddress(polygon.TokenAddress)

	txHash, err := polygon.RevokeApproval(client, token, spender, pk)
	assert.NoError(t, err, "RevokeApproval failed")
	assert.NotEmpty(t, txHash, "Transaction hash should not be empty")
	t.Logf("RevokeApproval tx hash: %s", txHash)
}
