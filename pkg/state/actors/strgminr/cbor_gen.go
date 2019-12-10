package strgminr

import (
	"fmt"
	"io"
	"sort"

	"github.com/libp2p/go-libp2p-core/peer"
	cbg "github.com/whyrusleeping/cbor-gen"
	xerrors "golang.org/x/xerrors"
)

/* This file was generated by github.com/whyrusleeping/cbor-gen */

var _ = xerrors.Errorf

func (t *StorageMinerActorState) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{139}); err != nil {
		return err
	}

	// t.t.PreCommittedSectors (map[string]*strgminr.PreCommittedSector) (map)
	{
		if err := cbg.CborWriteHeader(w, cbg.MajMap, uint64(len(t.PreCommittedSectors))); err != nil {
			return err
		}

		keys := make([]string, 0, len(t.PreCommittedSectors))
		for k := range t.PreCommittedSectors {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			v := t.PreCommittedSectors[k]

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

	// t.t.Sectors (cid.Cid) (struct)

	if err := cbg.WriteCid(w, t.Sectors); err != nil {
		return xerrors.Errorf("failed to write cid field t.Sectors: %w", err)
	}

	// t.t.ProvingSet (cid.Cid) (struct)

	if err := cbg.WriteCid(w, t.ProvingSet); err != nil {
		return xerrors.Errorf("failed to write cid field t.ProvingSet: %w", err)
	}

	// t.t.Info (cid.Cid) (struct)

	if err := cbg.WriteCid(w, t.Info); err != nil {
		return xerrors.Errorf("failed to write cid field t.Info: %w", err)
	}

	// t.t.CurrentFaultSet (types.BitField) (struct)
	if err := t.CurrentFaultSet.MarshalCBOR(w); err != nil {
		return err
	}

	// t.t.NextFaultSet (types.BitField) (struct)
	if err := t.NextFaultSet.MarshalCBOR(w); err != nil {
		return err
	}

	// t.t.NextDoneSet (types.BitField) (struct)
	if err := t.NextDoneSet.MarshalCBOR(w); err != nil {
		return err
	}

	// t.t.Power (types.BigInt) (struct)
	if err := t.Power.MarshalCBOR(w); err != nil {
		return err
	}

	// t.t.Active (bool) (bool)
	if err := cbg.WriteBool(w, t.Active); err != nil {
		return err
	}

	// t.t.SlashedAt (uint64) (uint64)
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.SlashedAt))); err != nil {
		return err
	}

	// t.t.ProvingPeriodEnd (uint64) (uint64)
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.ProvingPeriodEnd))); err != nil {
		return err
	}
	return nil
}

func (t *StorageMinerActorState) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 11 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.t.PreCommittedSectors (map[string]*strgminr.PreCommittedSector) (map)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajMap {
		return fmt.Errorf("expected a map (major type 5)")
	}
	if extra > 4096 {
		return fmt.Errorf("t.PreCommittedSectors: map too large")
	}

	t.PreCommittedSectors = make(map[string]*PreCommittedSector, extra)

	for i, l := 0, int(extra); i < l; i++ {

		var k string

		{
			sval, err := cbg.ReadString(br)
			if err != nil {
				return err
			}

			k = string(sval)
		}

		var v *PreCommittedSector

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
				v = new(PreCommittedSector)
				if err := v.UnmarshalCBOR(br); err != nil {
					return err
				}
			}

		}

		t.PreCommittedSectors[k] = v

	}
	// t.t.Sectors (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(br)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.Sectors: %w", err)
		}

		t.Sectors = c

	}
	// t.t.ProvingSet (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(br)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.ProvingSet: %w", err)
		}

		t.ProvingSet = c

	}
	// t.t.Info (cid.Cid) (struct)

	{

		c, err := cbg.ReadCid(br)
		if err != nil {
			return xerrors.Errorf("failed to read cid field t.Info: %w", err)
		}

		t.Info = c

	}
	// t.t.CurrentFaultSet (types.BitField) (struct)

	{

		if err := t.CurrentFaultSet.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	// t.t.NextFaultSet (types.BitField) (struct)

	{

		if err := t.NextFaultSet.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	// t.t.NextDoneSet (types.BitField) (struct)

	{

		if err := t.NextDoneSet.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	// t.t.Power (types.BigInt) (struct)

	{

		if err := t.Power.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	// t.t.Active (bool) (bool)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajOther {
		return fmt.Errorf("booleans must be major type 7")
	}
	switch extra {
	case 20:
		t.Active = false
	case 21:
		t.Active = true
	default:
		return fmt.Errorf("booleans are either major type 7, value 20 or 21 (got %d)", extra)
	}
	// t.t.SlashedAt (uint64) (uint64)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajUnsignedInt {
		return fmt.Errorf("wrong type for uint64 field")
	}
	t.SlashedAt = uint64(extra)
	// t.t.ProvingPeriodEnd (uint64) (uint64)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajUnsignedInt {
		return fmt.Errorf("wrong type for uint64 field")
	}
	t.ProvingPeriodEnd = uint64(extra)
	return nil
}

func (t *MinerInfo) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{132}); err != nil {
		return err
	}

	// t.t.Owner (address.Address) (struct)
	if err := t.Owner.MarshalCBOR(w); err != nil {
		return err
	}

	// t.t.Worker (address.Address) (struct)
	if err := t.Worker.MarshalCBOR(w); err != nil {
		return err
	}

	// t.t.PeerID (peer.ID) (string)
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajTextString, uint64(len(t.PeerID)))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(t.PeerID)); err != nil {
		return err
	}

	// t.t.SectorSize (uint64) (uint64)
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.SectorSize))); err != nil {
		return err
	}
	return nil
}

func (t *MinerInfo) UnmarshalCBOR(r io.Reader) error {
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

	// t.t.Owner (address.Address) (struct)

	{

		if err := t.Owner.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	// t.t.Worker (address.Address) (struct)

	{

		if err := t.Worker.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	// t.t.PeerID (peer.ID) (string)

	{
		sval, err := cbg.ReadString(br)
		if err != nil {
			return err
		}

		t.PeerID = peer.ID(sval)
	}
	// t.t.SectorSize (uint64) (uint64)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajUnsignedInt {
		return fmt.Errorf("wrong type for uint64 field")
	}
	t.SectorSize = uint64(extra)
	return nil
}

func (t *PreCommittedSector) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{130}); err != nil {
		return err
	}

	// t.t.Info (strgminr.SectorPreCommitInfo) (struct)
	if err := t.Info.MarshalCBOR(w); err != nil {
		return err
	}

	// t.t.ReceivedEpoch (uint64) (uint64)
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.ReceivedEpoch))); err != nil {
		return err
	}
	return nil
}

func (t *PreCommittedSector) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 2 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.t.Info (strgminr.SectorPreCommitInfo) (struct)

	{

		if err := t.Info.UnmarshalCBOR(br); err != nil {
			return err
		}

	}
	// t.t.ReceivedEpoch (uint64) (uint64)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajUnsignedInt {
		return fmt.Errorf("wrong type for uint64 field")
	}
	t.ReceivedEpoch = uint64(extra)
	return nil
}

func (t *SectorPreCommitInfo) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{132}); err != nil {
		return err
	}

	// t.t.SectorNumber (uint64) (uint64)
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.SectorNumber))); err != nil {
		return err
	}

	// t.t.CommR ([]uint8) (slice)
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajByteString, uint64(len(t.CommR)))); err != nil {
		return err
	}
	if _, err := w.Write(t.CommR); err != nil {
		return err
	}

	// t.t.SealEpoch (uint64) (uint64)
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajUnsignedInt, uint64(t.SealEpoch))); err != nil {
		return err
	}

	// t.t.DealIDs ([]uint64) (slice)
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajArray, uint64(len(t.DealIDs)))); err != nil {
		return err
	}
	for _, v := range t.DealIDs {
		if err := cbg.CborWriteHeader(w, cbg.MajUnsignedInt, v); err != nil {
			return err
		}
	}
	return nil
}

func (t *SectorPreCommitInfo) UnmarshalCBOR(r io.Reader) error {
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

	// t.t.SectorNumber (uint64) (uint64)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajUnsignedInt {
		return fmt.Errorf("wrong type for uint64 field")
	}
	t.SectorNumber = uint64(extra)
	// t.t.CommR ([]uint8) (slice)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}

	if extra > cbg.ByteArrayMaxLen {
		return fmt.Errorf("t.CommR: byte array too large (%d)", extra)
	}
	if maj != cbg.MajByteString {
		return fmt.Errorf("expected byte array")
	}
	t.CommR = make([]byte, extra)
	if _, err := io.ReadFull(br, t.CommR); err != nil {
		return err
	}
	// t.t.SealEpoch (uint64) (uint64)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajUnsignedInt {
		return fmt.Errorf("wrong type for uint64 field")
	}
	t.SealEpoch = uint64(extra)
	// t.t.DealIDs ([]uint64) (slice)

	maj, extra, err = cbg.CborReadHeader(br)
	if err != nil {
		return err
	}

	if extra > cbg.MaxLength {
		return fmt.Errorf("t.DealIDs: array too large (%d)", extra)
	}

	if maj != cbg.MajArray {
		return fmt.Errorf("expected cbor array")
	}
	if extra > 0 {
		t.DealIDs = make([]uint64, extra)
	}
	for i := 0; i < int(extra); i++ {

		maj, val, err := cbg.CborReadHeader(br)
		if err != nil {
			return xerrors.Errorf("failed to read uint64 for t.DealIDs slice: %w", err)
		}

		if maj != cbg.MajUnsignedInt {
			return xerrors.Errorf("value read for array t.DealIDs was not a uint, instead got %d", maj)
		}

		t.DealIDs[i] = val
	}

	return nil
}

func (t *UpdatePeerIDParams) MarshalCBOR(w io.Writer) error {
	if t == nil {
		_, err := w.Write(cbg.CborNull)
		return err
	}
	if _, err := w.Write([]byte{129}); err != nil {
		return err
	}

	// t.t.PeerID (peer.ID) (string)
	if _, err := w.Write(cbg.CborEncodeMajorType(cbg.MajTextString, uint64(len(t.PeerID)))); err != nil {
		return err
	}
	if _, err := w.Write([]byte(t.PeerID)); err != nil {
		return err
	}
	return nil
}

func (t *UpdatePeerIDParams) UnmarshalCBOR(r io.Reader) error {
	br := cbg.GetPeeker(r)

	maj, extra, err := cbg.CborReadHeader(br)
	if err != nil {
		return err
	}
	if maj != cbg.MajArray {
		return fmt.Errorf("cbor input should be of type array")
	}

	if extra != 1 {
		return fmt.Errorf("cbor input had wrong number of fields")
	}

	// t.t.PeerID (peer.ID) (string)

	{
		sval, err := cbg.ReadString(br)
		if err != nil {
			return err
		}

		t.PeerID = peer.ID(sval)
	}
	return nil
}
