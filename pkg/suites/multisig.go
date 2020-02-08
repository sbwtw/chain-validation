package suites

import (
	"testing"

	address "github.com/filecoin-project/go-address"
	require "github.com/stretchr/testify/require"

	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	big_spec "github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	multisig_spec "github.com/filecoin-project/specs-actors/actors/builtin/multisig"
	exitcode_spec "github.com/filecoin-project/specs-actors/actors/runtime/exitcode"

	chain "github.com/filecoin-project/chain-validation/pkg/chain"
)

func MultiSigActorConstructor(t testing.TB, factory Factories) {
	var aliceBal = abi_spec.NewTokenAmount(200000000000)
	const numApprovals = 3
	const unlockDuration = 10

	var valueSend = abi_spec.NewTokenAmount(10)

	td := NewTestDriver(t, factory, map[address.Address]big_spec.Int{
		builtin_spec.InitActorAddr:         big_spec.NewInt(0),
		builtin_spec.BurntFundsActorAddr:   big_spec.NewInt(0),
		builtin_spec.StoragePowerActorAddr: big_spec.NewInt(0),
		builtin_spec.RewardActorAddr:       TotalNetworkBalance,
	})

	// creator of the multisig actor
	alice := td.Driver.NewAccountActor(aliceBal)

	// expected address of the actor
	multisigAddr, err := address.NewIDAddress(102)
	require.NoError(t, err)

	td.MustCreateAndVerifyMultisigActor(0, valueSend, multisigAddr, alice, &multisig_spec.ConstructorParams{
		Signers:               []address.Address{alice},
		NumApprovalsThreshold: numApprovals,
		UnlockDuration:        unlockDuration,
	})
}

func MultiSigActorProposeApprove(t testing.TB, factory Factories) {
	var initialBal = abi_spec.NewTokenAmount(200000000000)
	const numApprovals = 2
	const unlockDuration = 10
	var valueSend = abi_spec.NewTokenAmount(10)

	td := NewTestDriver(t, factory, map[address.Address]big_spec.Int{
		builtin_spec.InitActorAddr:         big_spec.NewInt(0),
		builtin_spec.BurntFundsActorAddr:   big_spec.NewInt(0),
		builtin_spec.StoragePowerActorAddr: big_spec.NewInt(0),
		builtin_spec.RewardActorAddr:       TotalNetworkBalance,
	})

	// Signers
	alice := td.Driver.NewAccountActor(initialBal)
	bob := td.Driver.NewAccountActor(initialBal)

	// Not Signer
	outsider := td.Driver.NewAccountActor(initialBal)

	// Multisig actor address
	multisigAddr, err := address.NewIDAddress(104)
	require.NoError(t, err)

	// create the multisig actor
	td.MustCreateAndVerifyMultisigActor(0, valueSend, multisigAddr, alice, &multisig_spec.ConstructorParams{
		Signers:               []address.Address{alice, bob},
		NumApprovalsThreshold: numApprovals,
		UnlockDuration:        unlockDuration,
	})

	// setup propose expected values and params
	txID0 := multisig_spec.TxnID(0)
	pparams := &multisig_spec.ProposeParams{
		To:     outsider,
		Value:  valueSend,
		Method: builtin_spec.MethodSend,
		Params: []byte{},
	}

	// propose the transaction and assert it exists in the actor state
	td.MustProposeMultisigTransfer(1, big_spec.Zero(), txID0, multisigAddr, alice, pparams)
	td.Driver.AssertMultisigTransaction(multisigAddr, txID0, multisig_spec.MultiSigTransaction{
		To:       pparams.To,
		Value:    pparams.Value,
		Method:   pparams.Method,
		Params:   pparams.Params,
		Approved: []address.Address{alice},
	})

	// outsider proposes themselves to receive 'valueSend' FIL. This fails as they are not a signer.
	td.ApplyMessageExpectReceipt(
		func() (*chain.Message, error) {
			return td.Producer.MultiSigPropose(multisigAddr, outsider, 0, &multisig_spec.ProposeParams{
				To:     outsider,
				Value:  valueSend,
				Method: builtin_spec.MethodSend,
				Params: []byte{},
			}, chain.Value(big_spec.Zero()))
		},
		chain.MessageReceipt{
			ExitCode:    exitcode_spec.ErrForbidden,
			ReturnValue: nil,
			GasUsed:     big_spec.NewInt(1000000),
		},
	)

	// outsider approves the value transfer alice sent. This fails as they are not a signer.
	td.ApplyMessageExpectReceipt(
		func() (*chain.Message, error) {
			return td.Producer.MultiSigApprove(multisigAddr, outsider, 1, txID0, chain.Value(big_spec.Zero()))
		},
		chain.MessageReceipt{
			ExitCode:    exitcode_spec.ErrForbidden,
			ReturnValue: nil,
			GasUsed:     big_spec.NewInt(1000000),
		},
	)

	// bob approves transfer of 'valueSend' FIL to outsider.
	txID1 := multisig_spec.TxnID(1)
	td.MustApproveMultisigActor(0, big_spec.Zero(), multisigAddr, bob, txID0)
	td.Driver.AssertMultisigState(multisigAddr, multisig_spec.MultiSigActorState{
		Signers:               []address.Address{alice, bob},
		NumApprovalsThreshold: numApprovals,
		NextTxnID:             txID1,
		InitialBalance:        valueSend,
		StartEpoch:            1,
		UnlockDuration:        unlockDuration,
	})
	td.Driver.AssertBalance(multisigAddr, big_spec.Zero())
	// TODO assert the pendingtxns is empty
}

func MultiSigActorProposeCancel(t testing.TB, factory Factories) {
	var initialBal = abi_spec.NewTokenAmount(200000000000)
	const numApprovals = 2
	const unlockDuration = 10
	var valueSend = abi_spec.NewTokenAmount(10)

	td := NewTestDriver(t, factory, map[address.Address]big_spec.Int{
		builtin_spec.InitActorAddr:         big_spec.NewInt(0),
		builtin_spec.BurntFundsActorAddr:   big_spec.NewInt(0),
		builtin_spec.StoragePowerActorAddr: big_spec.NewInt(0),
		builtin_spec.RewardActorAddr:       TotalNetworkBalance,
	})

	alice := td.Driver.NewAccountActor(initialBal)
	bob := td.Driver.NewAccountActor(initialBal)
	outsider := td.Driver.NewAccountActor(initialBal)

	multisigAddr, err := address.NewIDAddress(104)
	require.NoError(t, err)

	// create the multisig actor
	td.MustCreateAndVerifyMultisigActor(0, valueSend, multisigAddr, alice, &multisig_spec.ConstructorParams{
		Signers:               []address.Address{alice, bob},
		NumApprovalsThreshold: numApprovals,
		UnlockDuration:        unlockDuration,
	})

	// alice proposes that outsider should receive 'valueSend' FIL.
	txID0 := multisig_spec.TxnID(0)
	pparams := &multisig_spec.ProposeParams{
		To:     outsider,
		Value:  valueSend,
		Method: builtin_spec.MethodSend,
		Params: []byte{},
	}

	// propose the transaction and assert it exists in the actor state
	td.MustProposeMultisigTransfer(1, big_spec.Zero(), txID0, multisigAddr, alice, pparams)
	td.Driver.AssertMultisigTransaction(multisigAddr, txID0, multisig_spec.MultiSigTransaction{
		To:       pparams.To,
		Value:    pparams.Value,
		Method:   pparams.Method,
		Params:   pparams.Params,
		Approved: []address.Address{alice},
	})

	// bob cancels alice's transaction. This fails as bob did not create alice's transaction.
	td.ApplyMessageExpectReceipt(
		func() (*chain.Message, error) {
			return td.Producer.MultiSigCancel(multisigAddr, bob, 0, txID0, chain.Value(big_spec.Zero()))
		},
		chain.MessageReceipt{
			ExitCode:    exitcode_spec.ErrForbidden,
			ReturnValue: nil,
			GasUsed:     big_spec.NewInt(1000000),
		},
	)

	// alice cancels their transaction. The outsider doesn't receive any FIL, the multisig actor's balance is empty, and the
	// transaction is canceled.
	td.MustCancelMultisigActor(2, big_spec.Zero(), multisigAddr, alice, txID0)
	td.Driver.AssertMultisigState(multisigAddr, multisig_spec.MultiSigActorState{
		Signers:               []address.Address{alice, bob},
		NumApprovalsThreshold: numApprovals,
		NextTxnID:             1,
		InitialBalance:        valueSend,
		StartEpoch:            1,
		UnlockDuration:        unlockDuration,
	})

	td.Driver.AssertBalance(multisigAddr, valueSend)
	td.Driver.AssertBalance(outsider, initialBal)
}
