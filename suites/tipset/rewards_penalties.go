package tipset

import (
	"context"
	"testing"

	miner_spec "github.com/filecoin-project/specs-actors/actors/builtin/miner"

	addr "github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	"github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"
	"github.com/stretchr/testify/assert"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

// Test for semantically in/valid messages, including miner penalties.
func TipSetTest_MinerRewardsAndPenalties(t *testing.T, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState)

	acctDefaultBalance := abi.NewTokenAmount(1_000_000_000)
	sendValue := abi.NewTokenAmount(1)

	t.Run("ok simple send", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		tipB := drivers.NewTipSetMessageBuilder(td)
		miner := td.ExeCtx.Miner

		alicePk, aliceId := td.NewAccountActor(drivers.SECP, acctDefaultBalance)
		bobPk, bobId := td.NewAccountActor(drivers.SECP, acctDefaultBalance)

		// Exercise all combinations of ID and PK address for the sender.
		callSeq := uint64(0)
		for _, alice := range []addr.Address{alicePk, aliceId} {
			for _, bob := range []addr.Address{bobPk, bobId} {
				aBal := td.GetBalance(aliceId)
				bBal := td.GetBalance(bobId)
				prevRewards := td.GetRewardSummary()
				prevMinerBal := td.GetBalance(miner)

				// Process a block with two messages, a simple send back and forth between accounts.
				result := tipB.WithBlockBuilder(
					drivers.NewBlockBuilder(td, td.ExeCtx.Miner).
						WithBLSMessageOk(
							td.MessageProducer.Transfer(bob, alice, chain.Value(sendValue), chain.Nonce(callSeq)),
						).
						WithBLSMessageOk(
							td.MessageProducer.Transfer(alice, bob, chain.Value(sendValue), chain.Nonce(callSeq)),
						),
				).ApplyAndValidate()
				tipB.Clear()

				// Each account has paid gas fees.
				td.AssertBalance(aliceId, big.Sub(aBal, result.Receipts[0].GasUsed.Big()))
				td.AssertBalance(bobId, big.Sub(bBal, result.Receipts[1].GasUsed.Big()))
				gasSum := big.Add(result.Receipts[0].GasUsed.Big(), result.Receipts[1].GasUsed.Big()) // Exploit gas price = 1

				// Validate rewards are paid directly to miner
				newRewards := td.GetRewardSummary()

				// total supply should decrease by the last reward amount
				assert.Equal(t, big.Sub(prevRewards.Treasury, newRewards.LastPerEpochReward), newRewards.Treasury)

				// the miners balance should have increased by the reward amount
				thisReward := big.Add(newRewards.LastPerEpochReward, gasSum)
				assert.Equal(t, big.Add(prevMinerBal, thisReward), td.GetBalance(miner))

				// no money was burnt
				assert.Equal(t, big.Zero(), td.GetBalance(builtin.BurntFundsActorAddr))

				callSeq++
			}
		}
	})

	t.Run("penalize sender does't exist", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()
		bb := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)
		miner := td.ExeCtx.Miner

		_, receiver := td.NewAccountActor(drivers.SECP, acctDefaultBalance)
		badSenders := []addr.Address{
			utils.NewIDAddr(t, 1234),
			utils.NewSECP256K1Addr(t, "1234"),
			utils.NewBLSAddr(t, 1234),
			utils.NewActorAddr(t, "1234"),
		}

		for _, s := range badSenders {
			bb.WithBLSMessageAndCode(td.MessageProducer.Transfer(receiver, s, chain.Value(sendValue)),
				exitcode.SysErrSenderInvalid,
			)
		}

		prevRewards := td.GetRewardSummary()
		drivers.NewTipSetMessageBuilder(td).WithBlockBuilder(bb).ApplyAndValidate()

		// Nothing received, no actors created.
		td.AssertBalance(receiver, acctDefaultBalance)
		for _, s := range badSenders {
			td.AssertNoActor(s)
		}

		newRewards := td.GetRewardSummary()
		// The penalty charged to the miner is not present in the receipt so we just have to hardcode it here.
		gasPenalty := big.NewInt(350)

		// The penalty amount has been burnt by the reward actor, and subtracted from the miner's block reward
		validateRewards(td, prevRewards, newRewards, miner, big.Zero(), gasPenalty)
		td.AssertBalance(builtin.BurntFundsActorAddr, gasPenalty)
	})

	t.Run("penalize sender non account", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		miner := td.ExeCtx.Miner
		tb := drivers.NewTipSetMessageBuilder(td)
		bb := drivers.NewBlockBuilder(td, miner)

		_, receiver := td.NewAccountActor(drivers.SECP, acctDefaultBalance)
		// Various non-account actors that can't be top-level senders.
		senders := []addr.Address{
			builtin.SystemActorAddr,
			builtin.InitActorAddr,
			builtin.CronActorAddr,
			miner,
		}

		for _, sender := range senders {
			bb.WithBLSMessageAndCode(td.MessageProducer.Transfer(receiver, sender, chain.Value(sendValue)),
				exitcode.SysErrSenderInvalid)
		}
		prevRewards := td.GetRewardSummary()
		tb.WithBlockBuilder(bb).ApplyAndValidate()
		td.AssertBalance(receiver, acctDefaultBalance)

		newRewards := td.GetRewardSummary()
		// The penalty charged to the miner is not present in the receipt so we just have to hardcode it here.
		gasPenalty := big.NewInt(176)

		// The penalty amount has been burnt by the reward actor, and subtracted from the miner's block reward.
		validateRewards(td, prevRewards, newRewards, miner, big.Zero(), gasPenalty)
		td.AssertBalance(builtin.BurntFundsActorAddr, gasPenalty)
	})

	t.Run("penalize wrong callseqnum", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		miner := td.ExeCtx.Miner
		tb := drivers.NewTipSetMessageBuilder(td)
		bb := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		_, aliceId := td.NewAccountActor(drivers.BLS, acctDefaultBalance)

		bb.WithBLSMessageAndCode(
			td.MessageProducer.Transfer(builtin.BurntFundsActorAddr, aliceId, chain.Nonce(1)),
			exitcode.SysErrSenderStateInvalid,
		)

		prevRewards := td.GetRewardSummary()
		tb.WithBlockBuilder(bb).ApplyAndValidate()

		newRewards := td.GetRewardSummary()
		// The penalty charged to the miner is not present in the receipt so we just have to hardcode it here.
		gasPenalty := big.NewInt(40)
		validateRewards(td, prevRewards, newRewards, miner, big.Zero(), gasPenalty)
		td.AssertBalance(builtin.BurntFundsActorAddr, gasPenalty)
	})

	t.Run("miner penalty exceeds declared gas limit for BLS message", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		miner := td.ExeCtx.Miner
		tb := drivers.NewTipSetMessageBuilder(td)
		bb := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		alice, _ := td.NewAccountActor(drivers.BLS, acctDefaultBalance)

		gasPrice := int64(2)
		gasPenalty := int64(260)
		gasLimit := gasPenalty - 130

		bb.WithBLSMessageAndCode(
			td.MessageProducer.Transfer(builtin.BurntFundsActorAddr, alice,
				chain.Nonce(1), // cause the message application to fail resulting in a miner penalty.
				chain.GasPrice(gasPrice), chain.GasLimit(gasLimit)),
			exitcode.SysErrSenderStateInvalid,
		)

		prevRewards := td.GetRewardSummary()
		tb.WithBlockBuilder(bb).ApplyAndValidate()

		newRewards := td.GetRewardSummary()
		// The penalty charged to the miner is not present in the receipt so we just have to hardcode it here.
		validateRewards(td, prevRewards, newRewards, miner, big.Zero(), big.NewInt(gasPenalty))
		td.AssertBalance(builtin.BurntFundsActorAddr, big.NewInt(gasPenalty))
		td.AssertBalance(alice, acctDefaultBalance)
	})

	t.Run("miner penalty exceeds declared gas limit for SECP message", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		miner := td.ExeCtx.Miner
		tb := drivers.NewTipSetMessageBuilder(td)
		bb := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		alice, _ := td.NewAccountActor(drivers.SECP, acctDefaultBalance)

		gasPrice := int64(2)
		gasPenalty := int64(420)
		gasLimit := gasPenalty - 210
		bb.WithSECPMessageAndCode(
			td.MessageProducer.Transfer(builtin.BurntFundsActorAddr, alice,
				chain.Nonce(1), // cause the message application to fail resulting in a miner penalty.
				chain.GasPrice(gasPrice), chain.GasLimit(gasLimit)),
			exitcode.SysErrSenderStateInvalid,
		)

		prevRewards := td.GetRewardSummary()
		tb.WithBlockBuilder(bb).ApplyAndValidate()

		newRewards := td.GetRewardSummary()
		// The penalty charged to the miner is not present in the receipt so we just have to hardcode it here.
		validateRewards(td, prevRewards, newRewards, miner, big.Zero(), big.NewInt(gasPenalty))
		td.AssertBalance(builtin.BurntFundsActorAddr, big.NewInt(gasPenalty))
		td.AssertBalance(alice, acctDefaultBalance)
	})

	t.Run("penalize sender insufficient balance", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		miner := td.ExeCtx.Miner
		tb := drivers.NewTipSetMessageBuilder(td)
		bb := drivers.NewBlockBuilder(td, td.ExeCtx.Miner)

		halfBalance := abi.NewTokenAmount(10_000_000)
		_, aliceId := td.NewAccountActor(drivers.BLS, big.Add(halfBalance, halfBalance))

		// Attempt to whole balance, in two parts.
		// The second message should fail (insufficient balance to pay fees).
		bb.WithBLSMessageOk(
			td.MessageProducer.Transfer(builtin.BurntFundsActorAddr, aliceId, chain.Value(halfBalance)),
		).WithBLSMessageAndCode(
			td.MessageProducer.Transfer(builtin.BurntFundsActorAddr, aliceId, chain.Value(halfBalance), chain.Nonce(1)),
			exitcode.SysErrSenderStateInvalid,
		)

		prevRewards := td.GetRewardSummary()
		result := tb.WithBlockBuilder(bb).ApplyAndValidate()

		newRewards := td.GetRewardSummary()
		// The penalty charged to the miner is not present in the receipt so we just have to hardcode it here.
		gasPenalty := big.NewInt(48)
		validateRewards(td, prevRewards, newRewards, miner, result.Receipts[0].GasUsed.Big(), gasPenalty)
		td.AssertBalance(builtin.BurntFundsActorAddr, big.Add(halfBalance, gasPenalty))
	})

	t.Run("insufficient gas to cover return value", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		miner := td.ExeCtx.Miner
		tb := drivers.NewTipSetMessageBuilder(td)

		alice, _ := td.NewAccountActor(drivers.BLS, acctDefaultBalance)

		// get a successful result so we can determine how much gas it costs. We'll reduce this by 1 in a subsequent call
		// to test insufficient gas to cover return value.
		tracerResult := tb.WithBlockBuilder(
			drivers.NewBlockBuilder(td, td.ExeCtx.Miner).
				WithBLSMessageAndRet(td.MessageProducer.MinerControlAddresses(miner, alice, nil, chain.Nonce(0)),
					// required to satisfy testing methods, unrelated to current test.
					chain.MustSerialize(&miner_spec.GetControlAddressesReturn{
						Owner:  td.StateDriver.BuiltinMinerInfo().OwnerID,
						Worker: td.StateDriver.BuiltinMinerInfo().WorkerID,
					}),
				),
		).ApplyAndValidate()
		requiredGasLimit := tracerResult.Receipts[0].GasUsed

		/* now the test */
		tb.Clear()
		rewardsBefore := td.GetRewardSummary()
		minerBalanceBefore := td.GetBalance(miner)
		senderBalanceBefore := td.GetBalance(alice)
		td.ExeCtx.Epoch++

		// Apply the message again with a reduced gas limit
		// A value just one less than the required limit for success ensures that the gas limit will be reached
		// at the last possible gas charge, which is that for the return value size.
		gasLimit := requiredGasLimit - 1
		result := tb.WithBlockBuilder(
			drivers.NewBlockBuilder(td, td.ExeCtx.Miner).
				WithBLSMessageAndCode(td.MessageProducer.MinerControlAddresses(miner, alice, nil, chain.Nonce(1), chain.GasLimit(int64(gasLimit))),
					exitcode.SysErrOutOfGas,
				),
		).ApplyAndValidate()
		gasUsed := result.Receipts[0].GasUsed
		gasCost := gasLimit.Big() // Gas price is 1
		newRewards := td.GetRewardSummary()

		// Check the actual gas charged is equal to the gas limit rather than the amount consumed up to but excluding
		// the return value which is smaller than the gas limit.
		assert.Equal(t, gasLimit, gasUsed)

		// Check sender charged exactly the max cost.
		assert.Equal(td.T, big.Sub(senderBalanceBefore, gasCost), td.GetBalance(alice))

		// Check the miner earned exactly the max cost (plus block reward).
		thisRwd := big.Add(newRewards.LastPerEpochReward, gasCost)
		assert.Equal(td.T, big.Add(minerBalanceBefore, thisRwd), td.GetBalance(miner))
		assert.Equal(td.T, big.Sub(rewardsBefore.Treasury, newRewards.LastPerEpochReward), newRewards.Treasury)
	})

	// TODO more tests:
	// - miner penalty causes subsequent otherwise-valid message to have wrong nonce (another miner penalty)
	// - miner penalty followed by non-miner penalty with same nonce (in different block)
}

func validateRewards(td *drivers.TestDriver, prevRewards *drivers.RewardSummary, newRewards *drivers.RewardSummary, miner addr.Address, gasReward big.Int, gasPenalty big.Int) {
	rwd := big.Add(big.Sub(newRewards.LastPerEpochReward, gasPenalty), gasReward)
	assert.Equal(td.T, big.Add(prevRewards.LastPerEpochReward, rwd), td.GetBalance(miner))
	assert.Equal(td.T, big.Sub(prevRewards.Treasury, newRewards.LastPerEpochReward), newRewards.Treasury)
}
