package transaction

import (
	"github.com/noah-blockchain/noah-go-node/core/types"
	"github.com/noah-blockchain/noah-go-node/crypto"
	"github.com/noah-blockchain/noah-go-node/helpers"
	"github.com/noah-blockchain/noah-go-node/rlp"
	"math/big"
	"math/rand"
	"sync"
	"testing"
)

func TestEditCandidateTx(t *testing.T) {
	cState := getState()

	privateKey, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(privateKey.PublicKey)
	coin := types.GetBaseCoin()
	cState.AddBalance(addr, coin, helpers.NoahToQNoah(big.NewInt(1000000)))

	pubkey := make([]byte, 32)
	rand.Read(pubkey)

	cState.CreateCandidate(addr, addr, pubkey, 10, 0, types.GetBaseCoin(), helpers.NoahToQNoah(big.NewInt(1)))
	cState.CreateValidator(addr, pubkey, 10, 0, types.GetBaseCoin(), helpers.NoahToQNoah(big.NewInt(1)))

	newRewardAddress := types.Address{1}
	newOwnerAddress := types.Address{2}

	data := EditCandidateData{
		PubKey:        pubkey,
		RewardAddress: newRewardAddress,
		OwnerAddress:  newOwnerAddress,
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
		Type:          TypeEditCandidate,
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

	targetBalance, _ := big.NewInt(0).SetString("999990000000000000000000", 10)
	balance := cState.GetBalance(addr, coin)
	if balance.Cmp(targetBalance) != 0 {
		t.Fatalf("Target %s balance is not correct. Expected %s, got %s", coin, targetBalance, balance)
	}

	candidate := cState.GetStateCandidate(pubkey)

	if candidate == nil {
		t.Fatalf("Candidate not found")
	}

	if candidate.OwnerAddress != newOwnerAddress {
		t.Fatalf("OwnerAddress has not changed")
	}

	if candidate.RewardAddress != newRewardAddress {
		t.Fatalf("RewardAddress has not changed")
	}
}
