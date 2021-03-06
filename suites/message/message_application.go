package message

import (
	"context"
	"testing"

	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	paych_spec "github.com/filecoin-project/specs-actors/actors/builtin/paych"
	crypto_spec "github.com/filecoin-project/specs-actors/actors/crypto"

	"github.com/filecoin-project/specs-actors/actors/runtime/exitcode"

	"github.com/filecoin-project/chain-validation/chain"
	"github.com/filecoin-project/chain-validation/drivers"
	"github.com/filecoin-project/chain-validation/state"
	"github.com/filecoin-project/chain-validation/suites/utils"
)

func MessageTest_MessageApplicationEdgecases(t *testing.T, factory state.Factories) {
	builder := drivers.NewBuilder(context.Background(), factory).
		WithDefaultGasLimit(1_000_000).
		WithDefaultGasPrice(big_spec.NewInt(1)).
		WithActorState(drivers.DefaultBuiltinActorsState)

	var aliceBal = abi_spec.NewTokenAmount(1_000_000_000)
	var transferAmnt = abi_spec.NewTokenAmount(10)

	t.Run("fail to cover gas cost for message receipt on chain", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		alice, _ := td.NewAccountActor(drivers.SECP, aliceBal)
		td.ApplyFailure(
			td.MessageProducer.Transfer(alice, alice, chain.Value(transferAmnt), chain.Nonce(0), chain.GasPrice(1), chain.GasLimit(8)),
			exitcode.SysErrOutOfGas)
	})

	t.Run("not enough gas to pay message on-chain-size cost", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		alice, _ := td.NewAccountActor(drivers.SECP, aliceBal)
		// Expect Message application to fail due to lack of gas
		td.ApplyFailure(
			td.MessageProducer.Transfer(alice, alice, chain.Value(transferAmnt), chain.Nonce(0), chain.GasPrice(10), chain.GasLimit(1)),
			exitcode.SysErrOutOfGas)

		// Expect Message application to fail due to lack of gas when sender address is unknown
		unknown := utils.NewIDAddr(t, 10000000)
		td.ApplyFailure(
			td.MessageProducer.Transfer(alice, unknown, chain.Value(transferAmnt), chain.Nonce(0), chain.GasPrice(10), chain.GasLimit(1)),
			exitcode.SysErrOutOfGas)
	})

	t.Run("invalid actor CallSeqNum", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		alice, _ := td.NewAccountActor(drivers.SECP, aliceBal)

		// Expect Message application to fail due to callseqnum being invalid: 1 instead of 0
		td.ApplyFailure(
			td.MessageProducer.Transfer(alice, alice, chain.Value(transferAmnt), chain.Nonce(1)),
			exitcode.SysErrSenderStateInvalid)

		// Expect message application to fail due to unknow actor when call seq num is also incorrect
		unknown := utils.NewIDAddr(t, 10000000)
		td.ApplyFailure(
			td.MessageProducer.Transfer(alice, unknown, chain.Value(transferAmnt), chain.Nonce(1)),
			exitcode.SysErrSenderInvalid)
	})

	t.Run("abort during actor execution", func(t *testing.T) {
		td := builder.Build(t)
		defer td.Complete()

		const pcTimeLock = abi_spec.ChainEpoch(10)
		const pcLane = uint64(123)
		const pcNonce = uint64(1)
		var pcAmount = big_spec.NewInt(10)
		var initialBal = abi_spec.NewTokenAmount(200_000_000_000)
		var toSend = abi_spec.NewTokenAmount(10_000)
		var pcSig = &crypto_spec.Signature{
			Type: crypto_spec.SigTypeBLS,
			Data: []byte("Grrr im an invalid signature, I cause panics in the payment channel actor"),
		}

		// will create and send on payment channel
		sender, _ := td.NewAccountActor(drivers.SECP, initialBal)
		// will be receiver on paych
		receiver, receiverID := td.NewAccountActor(drivers.SECP, initialBal)

		// the _expected_ address of the payment channel
		paychAddr := utils.NewIDAddr(t, utils.IdFromAddress(receiverID)+1)
		createRet := td.ComputeInitActorExecReturn(sender, 0, 0, paychAddr)
		td.ApplyExpect(
			td.MessageProducer.CreatePaymentChannelActor(receiver, sender, chain.Value(toSend), chain.Nonce(0)),
			chain.MustSerialize(&createRet))

		// message application fails due to invalid argument (signature).
		td.ApplyFailure(
			td.MessageProducer.PaychUpdateChannelState(paychAddr, sender, &paych_spec.UpdateChannelStateParams{
				Sv: paych_spec.SignedVoucher{
					TimeLockMin: pcTimeLock,
					TimeLockMax: pcTimeLock,
					Lane:        pcLane,
					Nonce:       pcNonce,
					Amount:      pcAmount,
					Signature:   pcSig, // construct with invalid signature
				},
			}, chain.Nonce(1), chain.Value(big_spec.Zero())),
			exitcode.ErrIllegalArgument)
	})

	// TODO more tests:
	// - receiver ID/Actor address does not exist
	// - invalid method for receiver
	// - missing/mismatched params for receiver
	// - various out-of-gas cases
}
