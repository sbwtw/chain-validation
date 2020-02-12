package chain

import (
	"github.com/filecoin-project/go-address"
	builtin_spec "github.com/filecoin-project/specs-actors/actors/builtin"
	"github.com/filecoin-project/specs-actors/actors/builtin/cron"
	"github.com/filecoin-project/specs-actors/actors/util/adt"

	"github.com/filecoin-project/chain-validation/pkg/chain/types"
)

func (mp *MessageProducer) CronConstructor(to, from address.Address, params cron.ConstructorParams, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsCron.Constructor, ser, opts...)
}
func (mp *MessageProducer) CronEpochTick(to, from address.Address, params adt.EmptyValue, opts ...MsgOpt) *types.Message {
	ser := MustSerialize(&params)
	return mp.Build(to, from, builtin_spec.MethodsCron.EpochTick, ser, opts...)
}
