package packet

import (
	"errors"
	"fmt"
	"hash/crc32"
	"math/big"
	"strings"
)

//type Packet struct {
//	Buf []byte
//	Pos int
//}

// TODO: Add Packet constructors

//func New(size int) *Packet {
//
//}

// Readers

// G1 gets 1 byte.
func (p *Packet) G1() uint8 {
	b := make([]byte, 1)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
	}
	return b[0]
}

// G1B gets 1 signed byte.
func (p *Packet) G1B() int8 {
	return int8(p.G1())
}

// G1Alt1 gets 1 byte using alternate method 1.
func (p *Packet) G1Alt1() uint8 {
	return p.G1() - 128
}

// G1Alt2 gets 1 byte using alternate method 2.
func (p *Packet) G1Alt2() uint8 {
	return -p.G1()
}

// G1Alt3 gets 1 byte using alternate method 3.
func (p *Packet) G1Alt3() uint8 {
	return 128 - p.G1()
}

// G1BAlt1 gets 1 signed byte using alternate method 1.
// alt: returns 1 byte that is encoded using alt method 1 as a signed byte?
func (p *Packet) G1BAlt1() int8 {
	return int8(p.G1() - 128)
}

// G1BAlt2 gets 1 byte using alternate method 2.
func (p *Packet) G1BAlt2() int8 {
	return int8(-p.G1())
}

// G1BAlt3 gets 1 signed byte using alternate method 3.
func (p *Packet) G1BAlt3() int8 {
	return int8(128 - p.G1())
}

// G2 gets 2 bytes.
func (p *Packet) G2() uint16 {
	b := make([]byte, 2)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
	}
	return uint16(b[0])<<8 | uint16(b[1])
}

// G2S gets 2 signed bytes.
func (p *Packet) G2S() int16 {
	return int16(p.G2())
}

// G2 gets 2 bytes using alternate method 1.
// TODO: add unit test
// mine
func (p *Packet) G2Alt1() uint16 {
	b := make([]byte, 2)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
	}
	return uint16(b[0]) | uint16(b[1])<<8
}

// G2Alt2 gets 2 bytes using alternate method 2.
func (p *Packet) G2Alt2() uint16 {
	b := make([]byte, 2)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
	}
	return uint16(b[0])<<8 | uint16(b[1]-128)
}

// G2Alt3 gets 2 bytes using alternate method 3.
func (p *Packet) G2Alt3() uint16 {
	b := make([]byte, 2)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
	}
	return uint16(b[1])<<8 | uint16(b[0]-128)
}

// G2SAlt1 gets 2 signed bytes using alternate method 1.
// alt write: gets 2 signed bytes that are encoded using alternate method 1?
// alternate encoding method 1?
func (p *Packet) G2SAlt1() int16 {
	b := make([]byte, 2)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
	}
	return int16(b[1])<<8 | int16(b[0])
}

// G2SAlt2 gets 2 signed bytes using alternate method 2.
func (p *Packet) G2SAlt2() int16 {
	b := make([]byte, 2)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
	}
	return int16(b[0])<<8 | int16(b[1]-128)
}

// G2SAlt3 gets 2 signed bytes using alternate method 3.
func (p *Packet) G2SAlt3() int16 {
	b := make([]byte, 2)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
	}
	return int16(b[1])<<8 | int16(b[0]-128)
}

// G3 gets 3 bytes.
func (p *Packet) G3() uint32 {
	b := make([]byte, 3)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
	}
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

// G4 gets 4 bytes.
func (p *Packet) G4() uint32 {
	b := make([]byte, 4)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
	}
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

// G4Alt1 gets 4 bytes using alternate method 1.
func (p *Packet) G4Alt1() uint32 {
	b := make([]byte, 4)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
	}
	return uint32(b[3])<<24 | uint32(b[2])<<16 | uint32(b[1])<<8 | uint32(b[0])
}

// G4Alt2 gets 4 bytes using alternate method 2.
func (p *Packet) G4Alt2() uint32 {
	b := make([]byte, 4)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
	}
	return uint32(b[2])<<24 | uint32(b[3])<<16 | uint32(b[0])<<8 | uint32(b[1])
}

// G4Alt3 gets 4 bytes using alternate method 3.
func (p *Packet) G4Alt3() uint32 {
	b := make([]byte, 4)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
	}
	return uint32(b[1])<<24 | uint32(b[0])<<16 | uint32(b[3])<<8 | uint32(b[2])
}

// G8 gets 8 bytes.
func (p *Packet) G8() uint64 {
	return (uint64(p.G4()) << 32) + uint64(p.G4())
}

// GVarInt gets a variable-length int.
func (p *Packet) GVarInt() int32 {
	b := int8(p.Buf[p.Pos])
	p.Pos++
	var value int32 = 0

	for b < 0 {
		value = (int32(b)&0x7F | value) << 7
		b = int8(p.Buf[p.Pos])
		p.Pos++
	}

	return value | int32(b)
}

// GVarLong gets a variable-length long.
func (p *Packet) GVarLong(length int) int64 {
	bytes := length - 1
	if bytes < 0 || bytes > 7 {
		panic("length must be between 1 and 8 inclusive")
	}

	var value int64 = 0
	for shift := bytes * 8; shift >= 0; shift -= 8 {
		value |= int64(p.Buf[p.Pos]) << shift
		p.Pos++
	}
	return value
}

// GData gets data.
func (p *Packet) GData(dest []byte, length int) {
	for i := 0; i < length; i++ {
		dest[i] = p.Buf[p.Pos]
		p.Pos++
	}
}

// GDataAlt1 gets data using alternate method 1.
func (p *Packet) GDataAlt1(dest []byte, length int) {
	for i := length - 1; i >= 0; i-- {
		dest[i] = p.Buf[p.Pos]
		p.Pos++
	}
}

// GDataAlt3 gets data using alternate method 3.
func (p *Packet) GDataAlt3(dest []byte, length int) {
	for i := length - 1; i >= 0; i-- {
		dest[i] = p.Buf[p.Pos] - 128
		p.Pos++
	}
}

// GSmart gets a Smart value (range 0 to 32767).
func (p *Packet) GSmart() uint16 {
	if p.Buf[p.Pos] >= 128 {
		return p.G2() - 32768
	} else {
		return uint16(p.G1())
	}
}

// GSmartS gets a signed Smart value (range -16384 to 16383).
func (p *Packet) GSmartS() int32 {
	if p.Buf[p.Pos] >= 128 {
		return int32(p.G2() - 49152)
	} else {
		return int32(p.G1() - 64)
	}
}

// GExtended1or2 gets an extended range of Smart values.
func (p *Packet) GExtended1or2() uint32 {
	var value uint32 = 0
	next := p.GSmart()

	for next == 32767 {
		next = p.GSmart()
		value += 32767
	}

	return value + uint32(next)
}

// GJStr gets a JagString.
func (p *Packet) GJStr() string {
	// TODO: review the Packet.java version for charset
	start := p.Pos
	for p.Buf[p.Pos] != 0 {
		p.Pos++
	}
	p.Pos++
	length := p.Pos - start - 1
	return string(p.Buf[start : start+length])
}

// GJStr2 gets a versioned JagString.
func (p *Packet) GJStr2() string {
	// TODO: review the Packet.java version for charset
	version := p.Buf[p.Pos]
	p.Pos++
	if version != 0 {
		// TODO: recover the panic upstream and close client conn?
		panic("bad version number")
	}

	return p.GJStr()
}

// FastGJStr gets a JagString quickly?
func (p *Packet) FastGJStr() string {
	if p.Buf[p.Pos] == 0 {
		p.Pos++
		return ""
	} else {
		return p.GJStr()
	}
}

////////////////////////////

// Writers

// P1 puts 1 byte.
func (p *Packet) P1(value uint8) {
	_, err := p.Write([]byte{value})
	if err != nil {
		panic(err)
	}
}

// P1Alt1 puts 1 byte using alternate method 1.
func (p *Packet) P1Alt1(value uint8) {
	_, err := p.Write([]byte{value + 128})
	if err != nil {
		panic(err)
	}
}

// P1Alt2 puts 1 byte using alternate method 2.
func (p *Packet) P1Alt2(value uint8) {
	_, err := p.Write([]byte{-value})
	if err != nil {
		panic(err)
	}
}

// P1Alt3 puts 1 byte using alternate method 3.
func (p *Packet) P1Alt3(value uint8) {
	_, err := p.Write([]byte{128 - value})
	if err != nil {
		panic(err)
	}
}

// P2 puts 2 bytes.
func (p *Packet) P2(value uint16) {
	_, err := p.Write([]byte{
		uint8(value >> 8),
		uint8(value),
	})
	if err != nil {
		panic(err)
	}
}

// P2Alt1 puts 2 bytes using alternate method 1.
// TODO: this is the same as IP2
func (p *Packet) P2Alt1(value uint16) {
	_, err := p.Write([]byte{
		uint8(value),
		uint8(value >> 8),
	})
	if err != nil {
		panic(err)
	}
}

// P2Alt2 puts 2 bytes using alternate method 2.
func (p *Packet) P2Alt2(value uint16) {
	_, err := p.Write([]byte{
		uint8(value >> 8),
		uint8(value + 128),
	})
	if err != nil {
		panic(err)
	}
}

// P2Alt3 puts 2 bytes using alternate method 3.
func (p *Packet) P2Alt3(value uint16) {
	_, err := p.Write([]byte{
		uint8(value + 128),
		uint8(value >> 8),
	})
	if err != nil {
		panic(err)
	}
}

// P3 puts 3 bytes.
func (p *Packet) P3(value uint32) {
	_, err := p.Write([]byte{
		uint8(value >> 16),
		uint8(value >> 8),
		uint8(value),
	})
	if err != nil {
		panic(err)
	}
}

// P4 puts 4 bytes.
func (p *Packet) P4(value uint32) {
	_, err := p.Write([]byte{
		uint8(value >> 24),
		uint8(value >> 16),
		uint8(value >> 8),
		uint8(value),
	})
	if err != nil {
		panic(err)
	}
}

// P4Alt1 puts 4 bytes using alternate method 1.
func (p *Packet) P4Alt1(value uint32) {
	_, err := p.Write([]byte{
		uint8(value),
		uint8(value >> 8),
		uint8(value >> 16),
		uint8(value >> 24),
	})
	if err != nil {
		panic(err)
	}
}

// P4Alt2 puts 4 bytes using alternate method 2.
func (p *Packet) P4Alt2(value uint32) {
	_, err := p.Write([]byte{
		uint8(value >> 8),
		uint8(value),
		uint8(value >> 24),
		uint8(value >> 16),
	})
	if err != nil {
		panic(err)
	}
}

// P4Alt3 puts 4 bytes using alternate method 3.
func (p *Packet) P4Alt3(value uint32) {
	_, err := p.Write([]byte{
		uint8(value >> 16),
		uint8(value >> 24),
		uint8(value),
		uint8(value >> 8),
	})
	if err != nil {
		panic(err)
	}
}

// P8 puts 8 bytes.
func (p *Packet) P8(value uint64) {
	_, err := p.Write([]byte{
		uint8(value >> 56),
		uint8(value >> 48),
		uint8(value >> 40),
		uint8(value >> 32),
		uint8(value >> 24),
		uint8(value >> 16),
		uint8(value >> 8),
		uint8(value),
	})
	if err != nil {
		panic(err)
	}
}

// PVarInt puts a variable length int.
func (p *Packet) PVarInt(value uint32) {
	if value&0xFFFFFF80 != 0 {
		if value&0xFFFFC000 != 0 {
			if value&0xFFE00000 != 0 {
				if value&0xF0000000 != 0 {
					// TODO: is the uint32 conv needed? check in debugger
					// TODO: check in debugger how large these values end up being
					// might be able to optimize
					p.P1(uint8(value>>28 | 0x80))
				}
				p.P1(uint8(value>>21 | 0x80))
			}
			p.P1(uint8(value>>14 | 0x80))
		}
		p.P1(uint8(value>>7 | 0x80))
	}
	p.P1(uint8(value & 0x7F))
}

// PVarLong puts a variable length long.
func (p *Packet) PVarLong(length int, value int64) {
	bytes := length - 1
	if bytes < 0 || bytes > 7 {
		panic("length must be between 1 and 8 inclusive")
	}

	for shift := bytes * 8; shift >= 0; shift -= 8 {
		_, err := p.Write([]byte{byte(value >> shift)})
		if err != nil {
			panic(err)
		}
	}
}

// PData puts data.
func (p *Packet) PData(src []byte, length int) {
	_, err := p.Write(src[:length])
	if err != nil {
		panic(err)
	}
}

// write data in inverse order
// mine
// TODO: write unit test
// TODO: rename to PDataAlt1 or something? check Gdata funcs for the patterns
func (p *Packet) IPData(src []byte, length int) {
	for i := length - 1; i >= 0; i-- {
		err := p.WriteByte(src[i])
		if err != nil {
			panic(err)
		}
	}
}

// PSmart puts a Smart value.
func (p *Packet) PSmart(value uint16) {
	if value >= 0 && value < 128 {
		p.P1(uint8(value))
	} else if value >= 0 && value < 32768 {
		p.P2(value + 32768)
	} else {
		panic("bad value")
	}
}

// PJStr puts a JagString.
func (p *Packet) PJStr(str string) {
	if firstNul := strings.IndexByte(str, 0); firstNul >= 0 {
		panic(fmt.Sprintf("NUL character at %v - cannot PJStr", firstNul))
	}

	// TODO: Use client Cp1252Charset
	for _, r := range str {
		_, err := p.Write([]byte{uint8(r)})
		if err != nil {
			panic(err)
		}
	}
	_, err := p.Write([]byte{0})
	if err != nil {
		panic(err)
	}
}

// TODO: add a test for this
// TODO: make it as robust as PJStr
// mine
// PJStr2 puts a versioned JagString.
func (p *Packet) PJStr2(str string) {
	p.P1(0) // version prepended
	for _, r := range str {
		_, err := p.Write([]byte{uint8(r)})
		if err != nil {
			panic(err)
		}
	}
	_, err := p.Write([]byte{0}) // null-terminated
	if err != nil {
		panic(err)
	}
}

// PSize1 puts a 1 byte size?
func (p *Packet) PSize1(length int) {
	p.Buf[len(p.Buf)-length-1] = uint8(length)
}

// PSize2 puts a size of 2 bytes?
func (p *Packet) PSize2(length int) {
	p.Buf[len(p.Buf)-length-2] = uint8(length >> 8)
	p.Buf[len(p.Buf)-length-1] = uint8(length)
}

// PSize4 puts the size of a byte sequence in the buffer
// as 4 bytes preceding the sequence.
func (p *Packet) PSize4(length int) {
	p.Buf[len(p.Buf)-length-4] = uint8(length >> 24)
	p.Buf[len(p.Buf)-length-3] = uint8(length >> 16)
	p.Buf[len(p.Buf)-length-2] = uint8(length >> 8)
	p.Buf[len(p.Buf)-length-1] = uint8(length)
}

/////////////////////////

// AddCRC adds a checksum.
func (p *Packet) AddCRC(off int) uint32 {
	checksum := GetCRC(p.Len()-off, off, p.Buf)
	p.P4(checksum)
	return checksum
}

// CheckCRC compares checksums.
func (p *Packet) CheckCRC() bool {
	p.Pos += p.Len() - 4
	thisCrc := GetCRC(p.Pos, 0, p.Buf)
	otherCrc := p.G4()
	p.Pos -= len(p.Buf)
	return otherCrc == thisCrc
}

// TinyEnc XTEA encrypt
func (p *Packet) TinyEnc(key []uint32) {
	blocks := p.Len() / 8

	for i := 0; i < blocks; i++ {
		v0 := p.G4()
		v1 := p.G4()
		sum := uint32(0)
		rounds := 32

		for rounds > 0 {
			rounds--
			v0 += sum + key[sum&0x3] ^ (v1 + ((v1 >> 5) ^ (v1 << 4)))
			sum += 0x9E3779B9
			v1 += (v0 + (v0>>5 ^ v0<<4)) ^ (key[sum>>11&0x3] + sum)
		}

		p.P4(v0)
		p.P4(v1)
	}
}

// TinyDec Extended Tiny Encryption Algorithm decrypt
func (p *Packet) TinyDec(offset int, key []uint32, length int) {
	blocks := (length - offset) / 8

	for i := 0; i < blocks; i++ {
		v0 := p.G4()
		v1 := p.G4()
		sum := uint32(0xC6EF3720)
		rounds := 32

		for rounds > 0 {
			rounds--
			v1 -= (key[sum>>11&0x3] + sum) ^ (((v0 >> 5) ^ (v0 << 4)) + v0)
			sum -= 0x9E3779B9
			v0 -= ((v1>>5 ^ v1<<4) + v1) ^ (key[sum&0x3] + sum)
		}

		p.P4(v0)
		p.P4(v1)
	}

	//p.Pos = start
}

// RSAEnc RSA-encrypts the buffer contents.
func (p *Packet) RSAEnc(modulus *big.Int, exponent *big.Int) {
	//length := p.Pos
	length := p.Len()
	//p.Pos = 0

	plaintextBytes := make([]byte, length)
	p.GData(plaintextBytes, length)

	plaintext := new(big.Int).SetBytes(plaintextBytes)
	ciphertext := plaintext.Exp(plaintext, exponent, modulus)
	ciphertextBytes := ciphertext.Bytes()

	//p.Pos = 0
	p.Reset()
	p.P1(uint8(len(ciphertextBytes)))
	p.PData(ciphertextBytes, len(ciphertextBytes))
}

// TODO: add a test for this
func (p *Packet) RSADec() (*Packet, error) {
	// we aren't using BigInteger, so we have to do this manually
	numBytes := p.G1()
	rsax := make([]byte, numBytes)
	p.GData(rsax, int(numBytes))
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
	//decryptedBuf := NewBuffer(decrypted)
	decryptedBuf := NewPacket(decrypted)

	// BigInteger would also remove all the preceding 0s, so we seek past them
	//for decryptedBuf.Peek1() == 0 {
	//	decryptedBuf.Seek(1)
	//}
	for decryptedBuf.Buf[decryptedBuf.Pos] == 0 {
		decryptedBuf.G1()
	}

	return decryptedBuf, nil
}

// IP2 puts 2 bytes in inverse order. ?
// TODO: this is the same as P2Alt1
func (p *Packet) IP2(value uint16) {
	_, err := p.Write([]byte{
		uint8(value),
		uint8(value >> 8),
	})
	if err != nil {
		panic(err)
	}
}

// IP4 puts 4 bytes in inverse order. ?
func (p *Packet) IP4(value uint32) {
	_, err := p.Write([]byte{
		uint8(value),
		uint8(value >> 8),
		uint8(value >> 16),
		uint8(value >> 24),
	})
	if err != nil {
		panic(err)
	}
}

// PJStrLen returns the size of PJStr output for s.
func PJStrLen(str string) int {
	return len(str) + 1
}

// GetCRC calculate checksum
func GetCRC(length int, offset int, src []uint8) uint32 {
	return crc32.ChecksumIEEE(src[offset : offset+length])
}
