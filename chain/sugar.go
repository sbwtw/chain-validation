package chain

import (
	"github.com/filecoin-project/go-address"
	abi_spec "github.com/filecoin-project/specs-actors/actors/abi"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	init_spec "github.com/filecoin-project/specs-actors/actors/builtin/init"
	multisig_spec "github.com/filecoin-project/specs-actors/actors/builtin/multisig"
	paych_spec "github.com/filecoin-project/specs-actors/actors/builtin/paych"
	power_spec "github.com/filecoin-project/specs-actors/actors/builtin/power"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/filecoin-project/chain-validation/chain/types"
)

var noParams []byte

// Transfer builds a simple value transfer message and returns it.
func (mp *MessageProducer) Transfer(to, from address.Address, opts ...MsgOpt) *types.Message {
	return mp.Build(to, from, builtin_spec.MethodSend, noParams, opts...)
}

func (mp *MessageProducer) CreatePaymentChannelActor(to, from address.Address, opts ...MsgOpt) *types.Message {
	return mp.InitExec(builtin_spec.InitActorAddr, from, &init_spec.ExecParams{
		CodeCID: builtin_spec.PaymentChannelActorCodeID,
		ConstructorParams: MustSerialize(&paych_spec.ConstructorParams{
			From: from,
			To:   to,
		}),
	}, opts...)
}

func (mp *MessageProducer) CreateMultisigActor(from address.Address, signers []address.Address, unlockDuration abi_spec.ChainEpoch, numApprovals int64, opts ...MsgOpt) *types.Message {
	return mp.InitExec(builtin_spec.InitActorAddr, from, &init_spec.ExecParams{
		CodeCID: builtin_spec.MultisigActorCodeID,
		ConstructorParams: MustSerialize(&multisig_spec.ConstructorParams{
			Signers:               signers,
			NumApprovalsThreshold: numApprovals,
			UnlockDuration:        unlockDuration,
		}),
	}, opts...)
}

func (mp *MessageProducer) CreateMinerActor(owner, worker address.Address, sealProofType abi_spec.RegisteredProof, pid peer.ID, opts ...MsgOpt) *types.Message {
	return mp.PowerCreateMiner(builtin_spec.StoragePowerActorAddr, owner, &power_spec.CreateMinerParams{
		Worker:        worker,
		Owner:         owner,
		SealProofType: sealProofType,
		Peer:          pid,
	}, opts...)
}
