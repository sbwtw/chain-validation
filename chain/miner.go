package chain

import (
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/specs-actors/actors/abi/big"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/miner"
	"github.com/filecoin-project/specs-actors/actors/builtin/power"
	"github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/chain/types"
)

func (mp *MessageProducer) MinerConstructor(to, from address.Address, params *power.MinerConstructorParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, builtin_spec.MethodsMiner.Constructor, ser, opts...)
}
func (mp *MessageProducer) MinerControlAddresses(to, from address.Address, params *adt.EmptyValue, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, builtin_spec.MethodsMiner.ControlAddresses, ser, opts...)
}
func (mp *MessageProducer) MinerChangeWorkerAddress(to, from address.Address, params *miner.ChangeWorkerAddressParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, builtin_spec.MethodsMiner.ChangeWorkerAddress, ser, opts...)
}
func (mp *MessageProducer) MinerChangePeerID(to, from address.Address, params *miner.ChangePeerIDParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, builtin_spec.MethodsMiner.ChangePeerID, ser, opts...)
}
func (mp *MessageProducer) MinerSubmitWindowedPoSt(to, from address.Address, params *miner.SubmitWindowedPoStParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, builtin_spec.MethodsMiner.SubmitWindowedPoSt, ser, opts...)
}
func (mp *MessageProducer) MinerPreCommitSector(to, from address.Address, params *miner.SectorPreCommitInfo, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, builtin_spec.MethodsMiner.PreCommitSector, ser, opts...)
}
func (mp *MessageProducer) MinerProveCommitSector(to, from address.Address, params *miner.ProveCommitSectorParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, builtin_spec.MethodsMiner.ProveCommitSector, ser, opts...)
}
func (mp *MessageProducer) MinerExtendSectorExpiration(to, from address.Address, params *miner.ExtendSectorExpirationParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, builtin_spec.MethodsMiner.ExtendSectorExpiration, ser, opts...)
}
func (mp *MessageProducer) MinerTerminateSectors(to, from address.Address, params *miner.TerminateSectorsParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, builtin_spec.MethodsMiner.TerminateSectors, ser, opts...)
}
func (mp *MessageProducer) MinerDeclareFaults(to, from address.Address, params *miner.DeclareFaultsParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, builtin_spec.MethodsMiner.DeclareFaults, ser, opts...)
}
func (mp *MessageProducer) MinerDeclareFaultsRecovered(to, from address.Address, params *miner.DeclareFaultsRecoveredParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, builtin_spec.MethodsMiner.DeclareFaultsRecovered, ser, opts...)
}
func (mp *MessageProducer) MinerOnDeferredCronEvent(to, from address.Address, params *miner.CronEventPayload, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, builtin_spec.MethodsMiner.OnDeferredCronEvent, ser, opts...)
}
func (mp *MessageProducer) MinerCheckSectorProven(to, from address.Address, params *miner.CheckSectorProvenParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, builtin_spec.MethodsMiner.CheckSectorProven, ser, opts...)
}
func (mp *MessageProducer) MinerAddLockedFund(to, from address.Address, params *big.Int, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, builtin_spec.MethodsMiner.AddLockedFund, ser, opts...)
}
func (mp *MessageProducer) MinerReportConsensusFault(to, from address.Address, params *miner.ReportConsensusFaultParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, builtin_spec.MethodsMiner.ReportConsensusFault, ser, opts...)
}
func (mp *MessageProducer) MinerWithdrawBalance(to, from address.Address, params *miner.WithdrawBalanceParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(params)
	return mp.Build(to, from, builtin_spec.MethodsMiner.WithdrawBalance, ser, opts...)
}
