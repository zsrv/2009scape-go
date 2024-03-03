package bytebuffer

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"math/big"
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

// Deprecated: TODO
func Alloc() {
}

// offset-related functions

func (b *ByteBuffer) Front() {
	b.off = 0
	b.bitOff = 0
}

func (b *ByteBuffer) Seek(num int) {
	b.off += num
}

// read-related functions

func (b *ByteBuffer) GData(numBytes int) ([]byte, error) {
	v := make([]byte, numBytes)
	_, err := b.Read(v)
	return v, err
}

func (b *ByteBuffer) GJStr() string {
	start := b.off
	for ; !b.empty() && b.Peek1() != 0; b.Seek(1) {
	}
	length := b.off - start
	if length == 0 {
		b.Seek(1)
		return ""
	}
	b.Front()
	b.Seek(start)
	str, err := b.GData(length)
	if err != nil {
		fmt.Println(err)
	}
	b.Seek(1)
	return string(str)
}

// write-related functions

func (b *ByteBuffer) PData(v *ByteBuffer) {
	_, _ = b.Write(v.Bytes())
}

func (b *ByteBuffer) PSize1(length int) {
	b.buf[len(b.buf)-length-1] = byte(length)
}

func (b *ByteBuffer) PSize2(length int) {
	b.buf[len(b.buf)-length-2] = byte(length >> 8)
	b.buf[len(b.buf)-length-1] = byte(length)
}

func (b *ByteBuffer) PJStr(s string) {
	for _, r := range s {
		b.P1(uint8(r))
	}
	b.P1(0)
}

func (b *ByteBuffer) PJStr2(s string) {
	b.P1(0) // version prepended
	for _, r := range s {
		b.P1(uint8(r))
	}
	b.P1(0) // null-terminated
}

func (b *ByteBuffer) PBool(v bool) {
	if v {
		b.P1(1)
	} else {
		b.P1(0)
	}
}

// 0 to 32767
func (b *ByteBuffer) PSmart(v uint16) {
	if v < 0x80 {
		b.P1(uint8(v))
	} else {
		b.P2(v + 0x8000)
	}
}

// -16384 to 16383
func (b *ByteBuffer) PSmartS(v int16) {
	if v < 0x80 {
		b.P1(uint8(v) + 0x40)
	} else {
		b.P2(uint16(v) + 0xC000)
	}
}

// bit-level access

// TODO

func (b *ByteBuffer) AccessBits() {
	b.bitOff = len(b.buf) << 3
}

func (b *ByteBuffer) PBit(n int, value int) {
	// TODO: lots of unsigned right shifts here (>>>) - find out how to replicate behavior in go
	bytePos := b.bitOff >> 3
	remaining := 8 - (b.bitOff & 7)
	b.bitOff += n

	// grow if necessary
	if bytePos+1 > b.Len() {
		//fmt.Printf("bytePos %v, b.Len() %v\n", bytePos, b.Len())
		//b.Grow((bytePos + 1) - b.Len())
		b.Write(make([]byte, (bytePos+1)-b.Len()))
	}

	for ; n > remaining; remaining = 8 {
		b.buf[bytePos] &= byte(^Bitmask[remaining])
		b.buf[bytePos] |= byte(uint32(value>>(n-remaining)) & Bitmask[remaining])
		bytePos += 1
		n -= remaining

		// grow if necessary
		if bytePos+1 > b.Len() {
			//b.Grow((bytePos + 1) - b.Len())
			b.Write(make([]byte, (bytePos+1)-b.Len()))
		}
	}

	if n == remaining {
		b.buf[bytePos] &= byte(^Bitmask[remaining])
		b.buf[bytePos] |= byte(value) & byte(Bitmask[remaining])
	} else {
		b.buf[bytePos] &= byte(int(^Bitmask[n]) << (remaining - n))
		b.buf[bytePos] |= byte((uint32(value) & Bitmask[n]) << (remaining - n))
	}
	// this.accessBytes(); // just in case mixed bit/byte access occurs
}

// TODO: gBit()

func (b *ByteBuffer) SeekBits(n int) {
	b.bitOff += n
}

// TODO: peekBits()

// similar to align() for bits
func (b *ByteBuffer) AccessBytes() {
	// TODO ??
}

func (b *ByteBuffer) TinyDec(key []uint32, len, off int) error {
	// TODO: missing offset stuff

	blocks := int(math.Ceil(float64((len - off) / 8)))
	for i := 0; i < blocks; i++ {
		v0, err := b.G4()
		if err != nil {
			return err
		}
		v1, err := b.G4()
		if err != nil {
			return err
		}
		sum := 0xC6EF3720
		rounds := 32

		for ; rounds > 0; rounds -= 1 {
			// TODO: unsigned right shifts
			v1 -= (v0<<4 ^ v0>>5) + v0 ^ uint32(sum) + key[sum>>11&3]
			sum -= 0x9E3779B9
			v0 -= (v1<<4 ^ v1>>5) + v1 ^ uint32(sum) + key[sum&3]
		}

		b.P4(v0)
		b.P4(v1)
	}

	return nil
}

func (b *ByteBuffer) RSADec() (*ByteBuffer, error) {
	// we aren't using BigInteger, so we have to do this manually
	numBytes, err := b.G1()
	if err != nil {
		return nil, err
	}
	rsax, err := b.GData(int(numBytes))
	if err != nil {
		return nil, err
	}
	if len(rsax) == 65 && rsax[0] == 0 {
		// Java BigInteger adds a 0 to indicate it's unsigned
		rsax = rsax[1:]
	} else if len(rsax) == 63 {
		// Java BigInteger didn't pad to 64
		temp := make([]byte, 64)
		copy(temp[1:], rsax)
		rsax = temp
	}

	// TODO: move this into an init() or something, and make key a package-level var?
	// private exponent
	keyD, ok := new(big.Int).SetString("571fb062048b61721ebfcf1e877153241b70c3aa26edb0f9f06a1b2be07c4e45eaba4fc356ea806cbed298d38613590a53fde0383c3a411758516293240925e5", 16)
	if !ok {
		return nil, errors.New("bad keyD")
	}
	// modulus
	keyN, ok := new(big.Int).SetString("0088c38748a58228f7261cdc340b5691d7d0975dee0ecdb717609e6bf971eb3fe723ef9d130e4686813739768ad9472eb46d8bfcc042c1a5fcb05e931f632eea5d", 16)
	if !ok {
		return nil, errors.New("bad keyN")
	}

	// RSA raw decryption (no padding)
	// better: take decrypt() from crypto/rsa/rsa.go
	c := new(big.Int).SetBytes(rsax)
	decrypted := c.Exp(c, keyD, keyN).Bytes()
	decryptedBuf := NewBuffer(decrypted)

	// BigInteger would also remove all the preceding 0s so we seek past them
	for decryptedBuf.Peek1() == 0 {
		decryptedBuf.Seek(1)
	}

	return decryptedBuf, nil
}

func (b *ByteBuffer) IG2() uint16 {
	v := uint16(b.buf[b.off] | b.buf[b.off+1]<<8)
	b.off += 2
	return v
}

func (b *ByteBuffer) IP2(v uint16) {
	b.P1(uint8(v))
	b.P1(uint8(v >> 8))
}

func (b *ByteBuffer) P2Add(v uint8) {
	b.P1(v >> 8)
	b.P1(v + 128)
}

func (b *ByteBuffer) P1Neg(v uint8) {
	b.P1(-v)
}

func (b *ByteBuffer) P1Sub(v uint8) {
	b.P1(128 - v)
}

func (b *ByteBuffer) IPData(data []uint8) { // TODO: or bytebuffer?
	for i := len(data) - 1; i >= 0; i-- {
		b.P1(data[i])
	}
}

// readers

// Get 1 byte
func (b *ByteBuffer) G1() (uint8, error) {
	res := make([]byte, 1)
	_, err := b.Read(res)
	return res[0], err
}

func (b *ByteBuffer) G1B() (int8, error) {
	res := make([]byte, 1)
	_, err := b.Read(res)
	return int8(res[0]), err
}

func (b *ByteBuffer) G2() (uint16, error) {
	res := make([]byte, 2)
	_, err := b.Read(res)
	return binary.BigEndian.Uint16(res), err
}

func (b *ByteBuffer) G2B() (int16, error) {
	res := make([]byte, 2)
	_, err := b.Read(res)
	return int16(binary.BigEndian.Uint16(res)), err
}

func (b *ByteBuffer) G4() (uint32, error) {
	res := make([]byte, 4)
	_, err := b.Read(res)
	return binary.BigEndian.Uint32(res), err
}

func (b *ByteBuffer) G4B() (int32, error) {
	res := make([]byte, 4)
	_, err := b.Read(res)
	return int32(binary.BigEndian.Uint32(res)), err
}

func (b *ByteBuffer) G8() (uint64, error) {
	res := make([]byte, 8)
	_, err := b.Read(res)
	return binary.BigEndian.Uint64(res), err
}

func (b *ByteBuffer) G8B() (int64, error) {
	res := make([]byte, 8)
	_, err := b.Read(res)
	return int64(binary.BigEndian.Uint64(res)), err
}

// peekers

func (b *ByteBuffer) Peek1() uint8 {
	return b.buf[b.off]
}

func (b *ByteBuffer) Peek4() uint32 {
	res := b.buf[b.off : b.off+4]
	return binary.BigEndian.Uint32(res)
}

// writers

// Put 1 byte
func (b *ByteBuffer) P1(v uint8) {
	_ = b.WriteByte(v)
}

func (b *ByteBuffer) P2(v uint16) {
	res := make([]byte, 2)
	binary.BigEndian.PutUint16(res, v)
	_, _ = b.Write(res)
}

func (b *ByteBuffer) P4(v uint32) {
	res := make([]byte, 4)
	binary.BigEndian.PutUint32(res, v)
	_, _ = b.Write(res)
}

func (b *ByteBuffer) P8(v uint64) {
	res := make([]byte, 8)
	binary.BigEndian.PutUint64(res, v)
	_, _ = b.Write(res)
}
