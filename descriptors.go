package cuei

import (
	bitter "github.com/futzu/bitter"
)

// audioCmpt is a struct for audioDscptr Components
type audioCmpt struct {
	ComponentTag  uint8
	ISOCode       uint32
	BitstreamMode uint8
	NumChannels   uint8
	FullSrvcAudio bool
}

// segCmpt Segmentation Descriptor Component
type segCmpt struct {
	ComponentTag uint8
	PtsOffset    float64
}

type Descriptor struct {
	Tag                              uint8       `json:",omitempty"`
	Length                           uint8       `json:",omitempty"`
	Identifier                       string      `json:",omitempty"`
	Name                             string      `json:",omitempty"`
	AudioComponents                  []audioCmpt `json:",omitempty"`
	ProviderAvailID                  uint32      `json:",omitempty"`
	PreRoll                          uint8       `json:",omitempty"`
	DTMFCount                        uint8       `json:",omitempty"`
	DTMFChars                        uint64      `json:",omitempty"`
	TAISeconds                       uint64      `json:",omitempty"`
	TAINano                          uint32      `json:",omitempty"`
	UTCOffset                        uint16      `json:",omitempty"`
	SegmentationEventID              string      `json:",omitempty"`
	SegmentationEventCancelIndicator bool        `json:",omitempty"`
	ProgramSegmentationFlag          bool        `json:",omitempty"`
	SegmentationDurationFlag         bool        `json:",omitempty"`
	DeliveryNotRestrictedFlag        bool        `json:",omitempty"`
	WebDeliveryAllowedFlag           bool        `json:",omitempty"`
	NoRegionalBlackoutFlag           bool        `json:",omitempty"`
	ArchiveAllowedFlag               bool        `json:",omitempty"`
	DeviceRestrictions               string      `json:",omitempty"`
	Components                       []segCmpt   `json:",omitempty"`
	SegmentationDuration             float64     `json:",omitempty"`
	SegmentationMessage              string      `json:",omitempty"`
	SegmentationUpidType             uint8       `json:",omitempty"`
	SegmentationUpidLength           uint8       `json:",omitempty"`
	SegmentationUpid                 *Upid       `json:",omitempty"`
	SegmentationTypeID               uint8       `json:",omitempty"`
	SegmentNum                       uint8       `json:",omitempty"`
	SegmentsExpected                 uint8       `json:",omitempty"`
	SubSegmentNum                    uint8       `json:",omitempty"`
	SubSegmentsExpected              uint8       `json:",omitempty"`
}

/*
*
Decode returns a Splice Descriptor by tag.

	The following Splice Descriptors are recognized.

	    0x0: Avail Descriptor,
	    0x1: DTMF Descriptor,
	    0x2: Segmentation Descriptor,
	    0x3: Time Descriptor,
	    0x4: Audio Descrioptor,

*
*/
func (dscptr *Descriptor) Decode(bd *bitter.Decoder, tag uint8, length uint8) {
	switch tag {
	case 0x0:
		dscptr.Tag = 0x0
		dscptr.availDescriptor(bd, tag, length)
	case 0x1:
		dscptr.Tag = 0x1
		dscptr.dtmfDescriptor(bd, tag, length)
	case 0x2:
		dscptr.Tag = 0x2
		dscptr.segmentationDescriptor(bd, tag, length)
	case 0x3:
		dscptr.Tag = 0x3
		dscptr.timeDescriptor(bd, tag, length)
	case 0x4:
		dscptr.Tag = 0x4
		dscptr.audioDescriptor(bd, tag, length)
	}
}

func (dscptr *Descriptor) audioDescriptor(bd *bitter.Decoder, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = bd.Ascii(32)
	ccount := bd.UInt8(4)
	bd.Forward(4)
	for ccount > 0 {
		ccount--
		ct := bd.UInt8(8)
		iso := bd.UInt32(24)
		bsm := bd.UInt8(3)
		nc := bd.UInt8(4)
		fsa := bd.Flag()
		dscptr.AudioComponents = append(dscptr.AudioComponents, audioCmpt{ct, iso, bsm, nc, fsa})
	}
}

// Decode for the avail Splice Descriptors
func (dscptr *Descriptor) availDescriptor(bd *bitter.Decoder, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = bd.Ascii(32)
	dscptr.Name = "Avail Descriptor"
	dscptr.ProviderAvailID = bd.UInt32(32)
}

// DTMF Splice Descriptor
func (dscptr *Descriptor) dtmfDescriptor(bd *bitter.Decoder, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = bd.Ascii(32)
	dscptr.Name = "DTMF Descriptor"
	dscptr.PreRoll = bd.UInt8(8)
	dscptr.DTMFCount = bd.UInt8(3)
	//bd.Forward(5)
	dscptr.DTMFChars = bd.UInt64(uint(8 * dscptr.DTMFCount))

}

// Decode for the Time Descriptor
func (dscptr *Descriptor) timeDescriptor(bd *bitter.Decoder, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = bd.Ascii(32)
	dscptr.Name = "Time Descriptor"
	dscptr.TAISeconds = bd.UInt64(48)
	dscptr.TAINano = bd.UInt32(32)
	dscptr.UTCOffset = bd.UInt16(16)
}

// Decode for the Segmentation Descriptor
func (dscptr *Descriptor) segmentationDescriptor(bd *bitter.Decoder, tag uint8, length uint8) {
	dscptr.Tag = tag
	dscptr.Length = length
	dscptr.Identifier = bd.Ascii(32)
	dscptr.Name = "Segmentation Descriptor"
	dscptr.SegmentationEventID = bd.Hex(32)
	dscptr.SegmentationEventCancelIndicator = bd.Flag()
	bd.Forward(7)
	if !dscptr.SegmentationEventCancelIndicator {
		dscptr.decodeSegFlags(bd)
		if !dscptr.ProgramSegmentationFlag {
			dscptr.decodeSegCmpnts(bd)
		}
		dscptr.decodeSegmentation(bd)
	}
}

func (dscptr *Descriptor) decodeSegFlags(bd *bitter.Decoder) {
	dscptr.ProgramSegmentationFlag = bd.Flag()
	dscptr.SegmentationDurationFlag = bd.Flag()
	dscptr.DeliveryNotRestrictedFlag = bd.Flag()
	if !dscptr.DeliveryNotRestrictedFlag {
		dscptr.WebDeliveryAllowedFlag = bd.Flag()
		dscptr.NoRegionalBlackoutFlag = bd.Flag()
		dscptr.ArchiveAllowedFlag = bd.Flag()
		dscptr.DeviceRestrictions = table20[bd.UInt8(2)]
		return
	}
	bd.Forward(5)
}

func (dscptr *Descriptor) decodeSegCmpnts(bd *bitter.Decoder) {
	ccount := bd.UInt8(8)
	for ccount > 0 { // 6 bytes each
		ccount--
		ct := bd.UInt8(8)
		bd.Forward(7)
		po := bd.As90k(33)
		dscptr.Components = append(dscptr.Components, segCmpt{ct, po})
	}
}

func (dscptr *Descriptor) decodeSegmentation(bd *bitter.Decoder) {
	if dscptr.SegmentationDurationFlag {
		dscptr.SegmentationDuration = bd.As90k(40)
	}
	dscptr.SegmentationUpidType = bd.UInt8(8)
	dscptr.SegmentationUpidLength = bd.UInt8(8)
	dscptr.SegmentationUpid = &Upid{}
	dscptr.SegmentationUpid.Decode(bd, dscptr.SegmentationUpidType, dscptr.SegmentationUpidLength)
	dscptr.SegmentationTypeID = bd.UInt8(8)

	mesg, ok := table22[dscptr.SegmentationTypeID]
	if ok {
		dscptr.SegmentationMessage = mesg
	}
	dscptr.SegmentNum = bd.UInt8(8)
	dscptr.SegmentsExpected = bd.UInt8(8)
	subSegIDs := []uint8{0x34, 0x36, 0x38, 0x3a}
	for _, ssid := range subSegIDs {
		if dscptr.SegmentationTypeID == ssid {
			dscptr.SubSegmentNum = bd.UInt8(8)
			dscptr.SubSegmentsExpected = bd.UInt8(8)
		}
	}
}
