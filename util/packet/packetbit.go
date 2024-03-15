package packet

import (
	"github.com/zsrv/rt5-server-go/util"
)

var Bitmask = []uint32{
	0,
	0x1, 0x3, 0x7, 0xF,
	0x1F, 0x3F, 0x7F, 0xFF,
	0x1FF, 0x3FF, 0x7FF, 0xFFF,
	0x1FFF, 0x3FFF, 0x7FFF, 0xFFFF,
	0x1FFFF, 0x3FFFF, 0x7FFFF, 0xFFFFF,
	0x1FFFFF, 0x3FFFFF, 0x7FFFFF, 0xFFFFFF,
	0x1FFFFFF, 0x3FFFFFF, 0x7FFFFFF, 0xFFFFFFF,
	0x1FFFFFFF, 0x3FFFFFFF, 0x7FFFFFFF, 0xFFFFFFFF,
}

type PacketBit struct {
	Packet
	bitOffset int
	random    *util.IsaacRandom
}

func NewPacketBit() {
}

// initialize isaac
func (p *PacketBit) SetKey(key []uint32) {
	p.random = util.NewIsaacRandom(key)
}

// get remaining bits at a position
func (p *PacketBit) AvailableBits(pos int) int {
	return pos*8 - p.bitOffset
}

// change stream position for bit access
func (p *PacketBit) AccessBits() {
	// p.Len() was p.Pos (which was wrong since we use Pos only for reads
	// and use len(Buf) for write pos
	p.bitOffset = p.Len() * 8
}

// change stream position for byte access
// must only be called after calling AccessBits!
func (p *PacketBit) AccessBytes() {
	// TODO: don't write to p.Pos here!!! it screws it up, maybe use a different
	//p.Pos = (p.bitOffset + 7) / 8
	// no op?
}

// Readers

// peek (do not advance), 1 byte, encrypted with isaac
func (p *PacketBit) Peek1Isaac() bool {
	value := uint32(p.Buf[p.Pos]) - p.random.PeekNext()&0xFF
	return value >= 128
}

func (p *PacketBit) G1Isaac() uint32 {
	value := uint32(p.Buf[p.Pos]) - p.random.GetNext()&0xFF
	p.Pos++
	if value < 128 {
		return value
	} else {
		p.Pos++
		return (uint32(p.Buf[p.Pos-1]) - p.random.GetNext()&0xFF) + (value - 128<<8)
	}
}

// get data encrypted with isaac
func (p *PacketBit) GIsaac(dest []uint8, length int) {
	for i := 0; i < length; i++ {
		// TODO: compare both
		//dest[i] = p.Buf[p.Pos] - uint8(p.random.GetNext())
		dest[i] = uint8(uint32(p.Buf[p.Pos]) - p.random.GetNext())
		p.Pos++
	}
}

// get n bits
func (p *PacketBit) GBit(n int) uint8 {
	bitPos := p.bitOffset >> 3
	bytePos := 8 - (p.bitOffset & 0x7)
	p.bitOffset += n
	value := uint8(0)
	for bytePos < n {
		value += (p.Buf[bitPos]&uint8(Bitmask[bytePos]))<<n - byte(bytePos)
		bitPos++
		n -= bytePos
		bytePos = 8
	}
	if n == bytePos {
		value += uint8(Bitmask[bytePos]) & p.Buf[bitPos]
	} else {
		value += p.Buf[bitPos]>>bytePos - uint8(n)&uint8(Bitmask[n])
	}
	return value
}

// Writers

// put 1 byte encrypted with isaac
func (p *PacketBit) P1Isaac(op uint8) {
	// TODO: compare both
	//p.Buf[p.Pos] = uint8(p.random.GetNext()) + op
	p.Buf[p.Pos] = uint8(p.random.GetNext() + uint32(op))
	p.Pos++
}

func (p *PacketBit) PBit(n int, value int) {
	// TODO: lots of unsigned right shifts here (>>>) - find out how to replicate behavior in go
	bytePos := p.bitOffset >> 3
	remaining := 8 - (p.bitOffset & 7)
	p.bitOffset += n

	// grow if necessary
	if bytePos+1 > p.Len() {
		//fmt.Printf("bytePos %v, b.Len() %v\n", bytePos, b.Len())
		//b.Grow((bytePos + 1) - b.Len())
		_, err := p.Write(make([]byte, (bytePos+1)-p.Len()))
		if err != nil {
			panic(err)
		}
		//fmt.Println(n)
	}

	for ; n > remaining; remaining = 8 {
		p.Buf[bytePos] &= byte(^Bitmask[remaining])
		p.Buf[bytePos] |= byte(uint32(value>>(n-remaining)) & Bitmask[remaining])
		bytePos += 1
		n -= remaining

		// grow if necessary
		if bytePos+1 > p.Len() {
			//b.Grow((bytePos + 1) - b.Len())
			p.Write(make([]byte, (bytePos+1)-p.Len()))
		}
	}

	if n == remaining {
		p.Buf[bytePos] &= byte(^Bitmask[remaining])
		p.Buf[bytePos] |= byte(value) & byte(Bitmask[remaining])
	} else {
		p.Buf[bytePos] &= byte(int(^Bitmask[n]) << (remaining - n))
		p.Buf[bytePos] |= byte((uint32(value) & Bitmask[n]) << (remaining - n))
	}
	// this.accessBytes(); // just in case mixed bit/byte access occurs
}

// custom
//func (p *PacketBit) PBitX(n int, value int) {
//	bytePos := int(uint(p.bitOffset) >> 3)
//	remaining := 8 - (p.bitOffset & 7)
//	p.bitOffset += n
//
//	// grow if necessary
//	//if bytePos+1 > p.Len() {
//	if bytePos+1 > len(p.Buf) {
//		//fmt.Printf("bytePos %v, p.Len() %v\n", bytePos, p.Len())
//		//p.Grow((bytePos + 1) - p.Len())
//		//p.Grow((bytePos + 1) - len(p.Buf))
//		x := make([]byte, (bytePos+1)-len(p.Buf))
//		p.Write(x)
//		//p.Write(make([]byte, (bytePos+1)-len(p.Buf)))
//	}
//
//	for ; n > remaining; remaining = 8 {
//		p.Buf[bytePos] &= byte(^Bitmask[remaining])
//		p.Buf[bytePos] |= byte(uint32(uint(value)>>(n-remaining)) & Bitmask[remaining])
//		bytePos++
//		n -= remaining
//
//		// grow if necessary
//		if bytePos+1 > len(p.Buf) {
//			//p.Grow((bytePos + 1) - p.Len())
//			p.Write(make([]byte, (bytePos+1)-len(p.Buf)))
//		}
//	}
//
//	if n == remaining {
//		p.Buf[bytePos] &= byte(^Bitmask[remaining])
//		//p.Buf[bytePos] |= byte(value) & byte(Bitmask[remaining])
//		p.Buf[bytePos] |= uint8(uint32(value) & Bitmask[remaining])
//	} else {
//		p.Buf[bytePos] &= byte(int(^Bitmask[n]) << (remaining - n))
//		p.Buf[bytePos] |= byte((uint32(value) & Bitmask[n]) << (remaining - n))
//	}
//	// this.accessBytes(); // just in case mixed bit/byte access occurs
//	return
//}
