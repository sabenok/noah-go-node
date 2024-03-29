package rewards

import (
	"math/big"

	"github.com/noah-blockchain/noah-go-node/core/types"
	"github.com/noah-blockchain/noah-go-node/helpers"
)

const lastBlock = 43702611 //todo
const firstReward = 333    // todo
const lastReward = 68      // todo

var startHeight uint64 = 0
var BeforeGenesis = big.NewInt(0)

func GetRewardForBlock(blockHeight uint64) *big.Int {
	blockHeight += startHeight

	if blockHeight > lastBlock {
		return big.NewInt(0)
	}

	if blockHeight == lastBlock {
		return helpers.NoahToQNoah(big.NewInt(lastReward))
	}

	reward := big.NewInt(firstReward)
	reward.Sub(reward, big.NewInt(int64(blockHeight/200000)))

	if reward.Cmp(types.Big0) < 1 {
		return helpers.NoahToQNoah(big.NewInt(1))
	}

	return helpers.NoahToQNoah(reward)
}

func SetStartHeight(sHeight uint64) {
	for i := uint64(1); i <= sHeight; i++ {
		BeforeGenesis.Add(BeforeGenesis, GetRewardForBlock(i))
	}

	startHeight = sHeight
}
