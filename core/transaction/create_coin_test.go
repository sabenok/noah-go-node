package transaction

import (
	"github.com/noah-blockchain/noah-go-node/core/types"
	"github.com/noah-blockchain/noah-go-node/crypto"
	"github.com/noah-blockchain/noah-go-node/helpers"
	"github.com/noah-blockchain/noah-go-node/rlp"
	"math/big"
	"sync"
	"testing"
)

func TestCreateCoinTx(t *testing.T) {
	cState := getState()

	privateKey, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(privateKey.PublicKey)
	coin := types.GetBaseCoin()

	cState.AddBalance(addr, coin, helpers.NoahToQNoah(big.NewInt(1000000)))

	var toCreate types.CoinSymbol
	copy(toCreate[:], []byte("ABCDEF"))

	reserve := helpers.NoahToQNoah(big.NewInt(1000))
	amount := helpers.NoahToQNoah(big.NewInt(100))
	crr := uint(50)
	name := "My Test Coin"

	data := CreateCoinData{
		Name:                 name,
		Symbol:               toCreate,
		InitialAmount:        amount,
		InitialReserve:       reserve,
		ConstantReserveRatio: crr,
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
		Type:          TypeCreateCoin,
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

	targetBalance, _ := big.NewInt(0).SetString("998000000000000000000000", 10)
	balance := cState.GetBalance(addr, coin)
	if balance.Cmp(targetBalance) != 0 {
		t.Fatalf("Target %s balance is not correct. Expected %s, got %s", coin, targetBalance, balance)
	}

	stateCoin := cState.GetStateCoin(toCreate)

	if stateCoin == nil {
		t.Fatalf("Coin %s not found in state", toCreate)
	}

	if stateCoin.ReserveBalance().Cmp(reserve) != 0 {
		t.Fatalf("Reserve balance in state is not correct. Expected %s, got %s", reserve, stateCoin.ReserveBalance())
	}

	if stateCoin.Volume().Cmp(amount) != 0 {
		t.Fatalf("Volume in state is not correct. Expected %s, got %s", amount, stateCoin.Volume())
	}

	if stateCoin.Crr() != crr {
		t.Fatalf("Crr in state is not correct. Expected %d, got %d", crr, stateCoin.Crr())
	}

	if stateCoin.Name() != name {
		t.Fatalf("Name in state is not correct. Expected %s, got %s", name, stateCoin.Name())
	}
}
