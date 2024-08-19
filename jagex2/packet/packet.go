package packet

import (
	"bufio"
	"container/list"
	"errors"
	"hash/crc32"
	"log"
	"math/big"
	"os"
	"path/filepath"
)

//type Packet struct {
//	Buf []byte
//	Pos int
//}

// TODO: Add Packet constructors

//func New(size int) *Packet {
//
//}

///////////////////////////////////

// static vars taken out of Packet
var cacheMinCount int
var cacheMidCount int
var cacheMaxCount int
var cacheBigCount int
var cacheHugeCount int
var cacheUnimaginableCount int

var cacheMin list.List
var cacheMid list.List
var cacheMax list.List
var cacheBig list.List
var cacheHuge list.List
var cacheUnimaginable list.List

func (p *Packet) Release() {
	// TODO: make Release() for PacketBit, see if I can super() for this one
	p.Pos = 0
	//p.bitPos = 0

	if len(p.Buf) == 100 && cacheMinCount < 1000 {
		cacheMin.PushBack(p)
		cacheMinCount++
	} else if len(p.Buf) == 5000 && cacheMidCount < 250 {
		cacheMid.PushBack(p)
		cacheMidCount++
	} else if len(p.Buf) == 30000 && cacheMaxCount < 50 {
		cacheMax.PushBack(p)
		cacheMaxCount++
	} else if len(p.Buf) == 100000 && cacheBigCount < 10 {
		cacheBig.PushBack(p)
		cacheBigCount++
	} else if len(p.Buf) == 500000 && cacheHugeCount < 5 {
		cacheHuge.PushBack(p)
		cacheHugeCount++
	} else if len(p.Buf) == 2000000 && cacheUnimaginableCount < 2 {
		cacheUnimaginable.PushBack(p)
		cacheUnimaginableCount++
	}
}

func Load(path string, seekToEnd bool) (*Packet, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	p := &Packet{Buf: file}

	if seekToEnd {
		p.Pos = len(p.Buf)
	}
	return p, nil
}

func (p *Packet) Save(filePath string, length int, start int) error {
	// TODO: make Save() for PacketBit, see if I can super() for this one

	dir := filepath.Dir(filePath)
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	_, err = w.Write(p.Buf[start : start+length])
	if err != nil {
		return err
	}
	w.Flush()

	return nil
}

//////////////////////////////////

// Readers

// GetCRC calculate checksum
func GetCRC(src []uint8, offset int, length int) uint32 {
	return crc32.ChecksumIEEE(src[offset : offset+length])
}

// CheckCRC
// TODO: should GetCRC and this be returning int64 or something?
func CheckCRC(src []uint8, offset int, length int, expected uint32) bool {
	checksum := GetCRC(src, offset, length)
	return checksum == expected
}

func AllocPacket(typ int) Packet {
	var p *Packet = nil

	if typ == 0 && cacheMinCount > 0 {
		p = cacheMin.Remove(cacheMin.Front()).(*Packet)
		cacheMinCount--
	} else if typ == 1 && cacheMidCount > 0 {
		p = cacheMid.Remove(cacheMid.Front()).(*Packet)
		cacheMidCount--
	} else if typ == 2 && cacheMaxCount > 0 {
		p = cacheMax.Remove(cacheMax.Front()).(*Packet)
		cacheMaxCount--
	} else if typ == 3 && cacheBigCount > 0 {
		p = cacheBig.Remove(cacheBig.Front()).(*Packet)
		cacheBigCount--
	} else if typ == 4 && cacheHugeCount > 0 {
		p = cacheHuge.Remove(cacheHuge.Front()).(*Packet)
		cacheHugeCount--
	} else if typ == 5 && cacheUnimaginableCount > 0 {
		p = cacheUnimaginable.Remove(cacheUnimaginable.Front()).(*Packet)
		cacheUnimaginableCount--
	}

	if p != nil {
		p.Pos = 0
		p.BitPos = 0
		return *p
	}

	switch typ {
	case 0:
		return Packet{Buf: make([]byte, 0, 100)}
	case 1:
		return Packet{Buf: make([]byte, 0, 5000)}
	case 2:
		return Packet{Buf: make([]byte, 0, 30000)}
	case 3:
		return Packet{Buf: make([]byte, 0, 100000)}
	case 4:
		return Packet{Buf: make([]byte, 0, 500000)}
	case 5:
		return Packet{Buf: make([]byte, 0, 2000000)}
	default:
		return Packet{Buf: make([]byte, 0, typ)}
	}

}

// G1 gets 1 unsigned byte.
// TODO: error isn't returned if there are no bytes to read sometimes. handle this for all getters somehow
func (p *Packet) G1() uint8 {
	b := make([]byte, 1)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
		//fmt.Println("ERROR: EOF") // this is hit sometimes with js5
	}
	return b[0]
}

// G1B gets 1 signed byte.
func (p *Packet) G1B() int8 {
	return int8(p.G1())
}

// G2 gets 2 unsigned bytes.
func (p *Packet) G2() uint16 {
	b := make([]byte, 2)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
		//fmt.Println("ERROR: G2 EOF")
	}
	return uint16(b[0])<<8 | uint16(b[1])
}

// G2S gets 2 signed bytes.
func (p *Packet) G2S() int16 {
	return int16(p.G2())
}

// IG2 gets 2 unsigned bytes represented in little-endian byte order.
func (p *Packet) IG2() uint16 {
	b := make([]byte, 2)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
	}
	return uint16(b[0]) | uint16(b[1])<<8
}

// G3 gets 3 unsigned bytes.
func (p *Packet) G3() uint32 {
	b := make([]byte, 3)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
	}
	return uint32(b[0])<<16 | uint32(b[1])<<8 | uint32(b[2])
}

// G4 gets 4 unsigned bytes.
func (p *Packet) G4() uint32 {
	b := make([]byte, 4)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
		//fmt.Println("ERROR: G4 PANIC")
	}
	return uint32(b[0])<<24 | uint32(b[1])<<16 | uint32(b[2])<<8 | uint32(b[3])
}

// IG4 gets 4 unsigned bytes represented in little-endian byte order.
func (p *Packet) IG4() uint32 {
	b := make([]byte, 4)
	_, err := p.Read(b)
	if err != nil {
		panic(err)
	}
	return uint32(b[3])<<24 | uint32(b[2])<<16 | uint32(b[1])<<8 | uint32(b[0])
}

// G8 gets 8 unsigned bytes.
func (p *Packet) G8() uint64 {
	return (uint64(p.G4()) << 32) + uint64(p.G4())
}

// GBool gets one byte and returns true if the value is 1,
// or false if the value is anything else.
func (p *Packet) GBool() bool {
	return p.G1() == 1
}

// GJStr gets a JagString, reading from the Packet
// until terminator is reached.
func (p *Packet) GJStr(terminator byte) string {
	// TODO: optimize this
	if p.Len() == 0 {
		log.Println("NO BYTES AVAILABLE IN GJSTR")
		return ""
	}
	// TODO: review the Packet.java version for charset
	start := p.Pos
	for p.Buf[p.Pos] != terminator {
		p.Pos++
	}
	p.Pos++
	length := p.Pos - start - 1
	return string(p.Buf[start : start+length])
}

// GJStrLF gets a newline-terminated JagString.
func (p *Packet) GJStrLF() string {
	return p.GJStr(10)
}

// GJStrNUL gets a NUL-terminated JagString.
func (p *Packet) GJStrNUL() string {
	return p.GJStr(0)
}

// GData gets data.
func (p *Packet) GData(dest []byte, length int) {
	// TODO: optimize
	for i := 0; i < length; i++ {
		dest[i] = p.Buf[p.Pos]
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
	// TODO: 2004scape server has this as uint.. maybe? maybe not
	if p.Buf[p.Pos] >= 128 {
		return int32(p.G2() - 49152)
	} else {
		return int32(p.G1() - 64)
	}
}

// BITS

////////////////////////////

// Writers

// P1 puts 1 unsigned byte.
func (p *Packet) P1(value uint8) {
	_, err := p.Write([]byte{value})
	if err != nil {
		panic(err)
	}
}

// P2 puts 2 unsigned bytes.
func (p *Packet) P2(value uint16) {
	_, err := p.Write([]byte{
		uint8(value >> 8),
		uint8(value),
	})
	if err != nil {
		panic(err)
	}
}

// IP2 puts 2 unsigned bytes represented in little-endian byte order.
func (p *Packet) IP2(value uint16) {
	_, err := p.Write([]byte{
		uint8(value),
		uint8(value >> 8),
	})
	if err != nil {
		panic(err)
	}
}

// P3 puts 3 unsigned bytes.
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

// P4 puts 4 unsigned bytes.
// TODO: 2004scape has this as int32
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

// IP4 puts 4 unsigned bytes represented in little-endian byte order.
// TODO: 2004scape has this as int32
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

// P8 puts 8 unsigned bytes.
// TODO: 2004scape has this as int64
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

// PBool puts 1 if the value is true, or 0 if the value is false.
func (p *Packet) PBool(value bool) {
	v := uint8(0)
	if value {
		v = 1
	}
	err := p.WriteByte(v)
	if err != nil {
		panic(err)
	}
}

// PJStr puts a JagString, terminated by terminator.
func (p *Packet) PJStr(str string, terminator byte) {
	//if firstNul := strings.IndexByte(str, 0); firstNul >= 0 {
	//	panic(fmt.Sprintf("NUL character at %v - cannot PJStr", firstNul))
	//}

	// TODO: Use client Cp1252Charset
	for _, r := range str {
		_, err := p.Write([]byte{uint8(r)})
		if err != nil {
			panic(err)
		}
	}
	_, err := p.Write([]byte{terminator})
	if err != nil {
		panic(err)
	}
}

// PJStrLF puts a newline-terminated JagString.
func (p *Packet) PJStrLF(str string) {
	p.PJStr(str, 10)
}

// PJStrNUL puts a NUL-terminated JagString.
func (p *Packet) PJStrNUL(str string) {
	p.PJStr(str, 0)
}

// PData puts data.
// TODO: might have to add offset arg
func (p *Packet) PData(src []byte, length int) {
	_, err := p.Write(src[:length])
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

// PSmart puts a Smart value.
// TODO: does it make sense to convert to unsigned?
func (p *Packet) PSmart(value int32) {
	if value >= 0 && value < 128 {
		p.P1(uint8(value))
	} else if value >= 0 && value < 32768 {
		p.P2(uint16(value + 32768))
	} else {
		panic("value out of range")
	}
}

// PSmartS puts a Smart value (signed?).
// TODO: does it make sense to convert to unsigned?
func (p *Packet) PSmartS(value int32) {
	if value < 64 && value >= -64 {
		p.P1(uint8(value + 64))
	} else if value < 16384 && value >= -16384 {
		p.P2(uint16(value + 49152))
	} else {
		panic("value out of range")
	}
}

/////////////////////////

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

func (p *Packet) RSADec() (*Packet, error) {
	// TODO: add a test for this
	// TODO: make two funcs: one that can use raw key components (for the original key)
	// and one that can use a PEM/DER key from disk or something (normal keys for later)

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
