package paych

import (
	"fmt"
	"io"
	"sort"

	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"

	"github.com/filecoin-project/chain-validation/pkg/state/types"
)

/* This file was generated by github.com/whyrusleeping/cbor-gen */

var _ = xerrors.Errorf

func (t *PaymentInfo) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{132}); err != nil {
		return err
	}

	// t.t.PayChActor (address.Address) (struct)
	if err := t.PayChActor.MarshalCBOR(w); err != nil {
		return err
	}

	// t.t.Payer (address.Address) (struct)
	if err := t.Payer.MarshalCBOR(w); err != nil {
		return err
	}

	// t.t.ChannelMessage (cid.Cid) (struct)

	if t.ChannelMessage == nil {
		if _, err := w.Write(cbg.CborNull); err != nil {
			return err
		}
	} else {
		if err := cbg.WriteCid(w, *t.ChannelMessage); err != nil {
			return xerrors.Errorf("failed to write cid field t.ChannelMessage: %w", err)
		}
	}

	// t.t.Vouchers ([]*types.SignedVoucher) (slice)
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajArray, uint64(len(t.Vouchers)))); err != nil {
		return err
	}
	for _, v := range t.Vouchers {
		if err := v.MarshalCBOR(w); err != nil {
			return err
		}
	}
	return nil
}

func (t *PaymentInfo) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 4 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.t.PayChActor (address.Address) (struct)

	{

		if err := t.PayChActor.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	// t.t.Payer (address.Address) (struct)

	{

		if err := t.Payer.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	// t.t.ChannelMessage (cid.Cid) (struct)

	{

		pb, err := br.PeekByte()
		if err != nil {
			return err
		}
		if pb == cbg.CborNull[0] {
			var nbuf [1]byte
			if _, err := br.Read(nbuf[:]); err != nil {
				return err
			}
		} else {

			c, err := cbg.ReadCid(br)
			if err != nil {
				return xerrors.Errorf("failed to read cid field t.ChannelMessage: %w", err)
			}

			t.ChannelMessage = &c
		}

	}
	// t.t.Vouchers ([]*types.SignedVoucher) (slice)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if extra > 8192 {
		return fmt.Errorf("t.Vouchers: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}
	if extra > 0 {
		t.Vouchers = make([]*types.SignedVoucher, extra)
	}
	for i := 0; i < int(extra); i++ {

		var v types.SignedVoucher
		if err := v.UnmarshalCBOR(br); err != nil {
			return err
		}

		t.Vouchers[i] = &v
	}

	return nil
}

func (t *PaymentChannelActorState) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{134}); err != nil {
		return err
	}

	// t.t.From (address.Address) (struct)
	if err := t.From.MarshalCBOR(w); err != nil {
		return err
	}

	// t.t.To (address.Address) (struct)
	if err := t.To.MarshalCBOR(w); err != nil {
		return err
	}

	// t.t.ToSend (types.BigInt) (struct)
	if err := t.ToSend.MarshalCBOR(w); err != nil {
		return err
	}

	// t.t.ClosingAt (uint64) (uint64)
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.ClosingAt))); err != nil {
		return err
	}

	// t.t.MinCloseHeight (uint64) (uint64)
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.MinCloseHeight))); err != nil {
		return err
	}

	// t.t.LaneStates (map[string]*paych.LaneState) (map)
	{
		if err := cbg.CborWriteHeader(w, cbg.MajMap, uint64(len(t.LaneStates))); err != nil {
			return err
		}

		keys := make([]string, 0, len(t.LaneStates))
		for k := range t.LaneStates {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := t.LaneStates[k]

			if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajTextString, uint64(len(k)))); err != nil {
				return err
			}
			if _, err := w.Write([]byte(k)); err != nil {
				return err
			}

			if err := v.MarshalCBOR(w); err != nil {
				return err
			}

		}
	}
	return nil
}

func (t *PaymentChannelActorState) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 6 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.t.From (address.Address) (struct)

	{

		if err := t.From.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	// t.t.To (address.Address) (struct)

	{

		if err := t.To.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	// t.t.ToSend (types.BigInt) (struct)

	{

		if err := t.ToSend.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	// t.t.ClosingAt (uint64) (uint64)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajUnsignedInt {
		return fmt.Errorf("wrong type for uint64 field")
	}
	t.ClosingAt = uint64(extra)
	// t.t.MinCloseHeight (uint64) (uint64)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajUnsignedInt {
		return fmt.Errorf("wrong type for uint64 field")
	}
	t.MinCloseHeight = uint64(extra)
	// t.t.LaneStates (map[string]*paych.LaneState) (map)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajMap {
		return fmt.Errorf("expected a map (major type 5)")
	}
	if extra > 4096 {
		return fmt.Errorf("t.LaneStates: map too large")
	}

	t.LaneStates = make(map[string]*LaneState, extra)

	for i, l := 0, int(extra); i < l; i++ {

		var k string

		{
			sval, err := cbg.ReadString(br)
			if err != nil {
				return err
			}

			k = string(sval)
		}

		var v *LaneState

		{

			pb, err := br.PeekByte()
			if err != nil {
				return err
			}
			if pb == cbg.CborNull[0] {
				var nbuf [1]byte
				if _, err := br.Read(nbuf[:]); err != nil {
					return err
				}
			} else {
				v = new(LaneState)
				if err := v.UnmarshalCBOR(br); err != nil {
					return err
				}
			}

		}

		t.LaneStates[k] = v

	}
	return nil
}

func (t *LaneState) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{131}); err != nil {
		return err
	}

	// t.t.Closed (bool) (bool)
	if err := cbg.WriteBool(w, t.Closed); err != nil {
		return err
	}

	// t.t.Redeemed (types.BigInt) (struct)
	if err := t.Redeemed.MarshalCBOR(w); err != nil {
		return err
	}

	// t.t.Nonce (uint64) (uint64)
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.Nonce))); err != nil {
		return err
	}
	return nil
}

func (t *LaneState) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 3 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.t.Closed (bool) (bool)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.Closed = false
	case 21:
		t.Closed = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	// t.t.Redeemed (types.BigInt) (struct)

	{

		if err := t.Redeemed.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	// t.t.Nonce (uint64) (uint64)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajUnsignedInt {
		return fmt.Errorf("wrong type for uint64 field")
	}
	t.Nonce = uint64(extra)
	return nil
}