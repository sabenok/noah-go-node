package transaction

import (
	"github.com/noah-blockchain/noah-go-node/core/state"
	"github.com/noah-blockchain/noah-go-node/core/types"
	"github.com/noah-blockchain/noah-go-node/crypto"
	"github.com/noah-blockchain/noah-go-node/helpers"
	"github.com/noah-blockchain/noah-go-node/rlp"
	"math/big"
	"sync"
	"testing"
)

func TestDeclareCandidacyTx(t *testing.T) {
	cState := getState()

	privateKey, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(privateKey.PublicKey)

	coin := types.GetBaseCoin()

	cState.AddBalance(addr, coin, helpers.NoahToQNoah(big.NewInt(1000000)))

	pkey, _ := crypto.GenerateKey()
	publicKey := crypto.FromECDSAPub(&pkey.PublicKey)[:32]

	commission := uint(10)

	data := DeclareCandidacyData{
		Address:    addr,
		PubKey:     publicKey,
		Commission: commission,
		Coin:       coin,
		Stake:      helpers.NoahToQNoah(big.NewInt(100)),
	}

	encodedData, err := rlp.EncodeToBytes(data)

	if err != nil {
		t.Fatal(err)
	}

	tx := Transaction{
		Nonce:         1,
		GasPrice:      1,
		ChainID:       types.CurrentChainID,
		GasCoin:       coin,
		Type:          TypeDeclareCandidacy,
		Data:          encodedData,
		SignatureType: SigTypeSingle,
	}

	if err := tx.Sign(privateKey); err != nil {
		t.Fatal(err)
	}

	encodedTx, err := rlp.EncodeToBytes(tx)

	if err != nil {
		t.Fatal(err)
	}

	response := RunTx(cState, false, encodedTx, big.NewInt(0), 0, sync.Map{}, 0)

	if response.Code != 0 {
		t.Fatalf("Response code is not 0. Error %s", response.Log)
	}

	targetBalance, _ := big.NewInt(0).SetString("999890000000000000000000", 10)
	balance := cState.GetBalance(addr, coin)
	if balance.Cmp(targetBalance) != 0 {
		t.Fatalf("Target %s balance is not correct. Expected %s, got %s", coin, targetBalance, balance)
	}

	candidate := cState.GetStateCandidate(publicKey)

	if candidate == nil {
		t.Fatalf("Candidate not found")
	}

	if candidate.OwnerAddress != addr {
		t.Fatalf("Owner address is not correct")
	}

	if candidate.RewardAddress != addr {
		t.Fatalf("Reward address is not correct")
	}

	if candidate.TotalNoahStake != nil && candidate.TotalNoahStake.Cmp(types.Big0) != 0 {
		t.Fatalf("Total stake is not correct")
	}

	if candidate.Commission != commission {
		t.Fatalf("Commission is not correct")
	}

	if candidate.Status != state.CandidateStatusOffline {
		t.Fatalf("Incorrect candidate status")
	}
}
