package cuei

import (
	"fmt"
	bitter "github.com/futzu/bitter"
)

var uriUpids = map[uint8]string{
	0x01: "Deprecated",
	0x02: "Deprecated",
	0x03: "AdID",
	0x07: "TID",
	0x08: "AiringID",
	0x09: "ADI",
	0x10: "UUID",
	0x11: "ACR",
	0x0a: "EIDR",
	0x0b: "ATSC",
	0x0c: "MPU",
	0x0d: "MID",
	0x0e: "ADS Info",
	0x0f: "URI",
}

/*
Upid is the Struct for Segmentation Upids

Non-standard UPID types are returned as bytes.
*/
type Upid struct {
	Name             string `json:",omitempty"`
	UpidType         uint8  `json:",omitempty"`
	Value            string `json:",omitempty"`
	TSID             uint16 `json:",omitempty"`
	Reserved         uint8  `json:",omitempty"`
	EndOfDay         uint8  `json:",omitempty"`
	UniqueFor        uint16 `json:",omitempty"`
	ContentID        []byte `json:",omitempty"`
	Upids            []Upid `json:",omitempty"`
	FormatIdentifier string `json:",omitempty"`
	PrivateData      []byte `json:",omitempty"`
}

// Decode Upids
func (upid *Upid) Decode(bd *bitter.Decoder, upidType uint8, upidlen uint8) {

	upid.UpidType = upidType

	name, ok := uriUpids[upidType]
	if ok {
		upid.Name = name
		upid.uri(bd, upidlen)
	} else {

		switch upidType {
		case 0x05, 0x06:
			upid.Name = "ISAN"
			upid.isan(bd, upidlen)
		case 0x08:
			upid.Name = "AiringID"
			upid.airid(bd, upidlen)
		case 0x0a:
			upid.Name = "EIDR"
			upid.eidr(bd, upidlen)
		case 0x0b:
			upid.Name = "ATSC"
			upid.atsc(bd, upidlen)
		case 0x0c:
			upid.Name = "MPU"
			upid.mpu(bd, upidlen)
		case 0x0d:
			upid.Name = "MID"
			upid.mid(bd, upidlen)
		default:
			upid.Name = "UPID"
			upid.uri(bd, upidlen)
		}
	}
}

// Decode for AirId
func (upid *Upid) airid(bd *bitter.Decoder, upidlen uint8) {
	upid.Value = bd.Hex(uint(upidlen << 3))
}

// Decode for Isan Upid
func (upid *Upid) isan(bd *bitter.Decoder, upidlen uint8) {
	upid.Value = bd.Ascii(uint(upidlen << 3))
}

// Decode for URI Upid
func (upid *Upid) uri(bd *bitter.Decoder, upidlen uint8) {
	upid.Value = bd.Ascii(uint(upidlen) << 3)
}

// Decode for ATSC Upid
func (upid *Upid) atsc(bd *bitter.Decoder, upidlen uint8) {
	upid.TSID = bd.UInt16(16)
	upid.Reserved = bd.UInt8(2)
	upid.EndOfDay = bd.UInt8(5)
	upid.UniqueFor = bd.UInt16(9)
	upid.ContentID = bd.Bytes(uint((upidlen - 4) << 3))
}

// Decode for EIDR Upid
func (upid *Upid) eidr(bd *bitter.Decoder, upidlen uint8) {
	if upidlen == 12 {
		head := bd.UInt64(16)
		tail := bd.Hex(80)
		upid.Value = fmt.Sprintf("10%v/%v", head, tail)
	}
}

// Decode for MPU Upid
func (upid *Upid) mpu(bd *bitter.Decoder, upidlen uint8) {
	ulb := uint(upidlen) << 3
	upid.FormatIdentifier = bd.Hex(32)
	upid.PrivateData = bd.Bytes(ulb - 32)
}

// Decode for MID Upid
func (upid *Upid) mid(bd *bitter.Decoder, upidlen uint8) {
	var i uint8
	i = 0
	for i < upidlen {
		utype := bd.UInt8(8)
		i++
		ulen := bd.UInt8(8)
		i++
		i += ulen
		var mupid Upid
		upid.Decode(bd, utype, ulen)
		upid.Upids = append(upid.Upids, mupid)
	}
}

// Encode Upids
func (upid *Upid) Encode(be *bitter.Encoder) {

	name, ok := uriUpids[upid.UpidType]
	if ok {
		upid.Name = name
		upid.encodeUri(be)
	}

}

func (upid *Upid) encodeUri(be *bitter.Encoder) {
	bites := []byte(upid.Value)
	bitlen := uint(len(bites) << 3)
	be.AddBytes(bites, bitlen)
}
