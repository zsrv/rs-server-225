// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package packet

// Simple byte buffer for marshaling data.

import (
	"errors"
	"io"
	"unicode/utf8"
)

// smallBufferSize is an initial allocation minimal capacity.
const smallBufferSize = 64

// A Packet is a variable-sized buffer of bytes with [Packet.Read] and [ByteBuffer.Write] methods.
// The zero value for Packet is an empty buffer ready to use.
type Packet struct {
	Buf      []byte // contents are the bytes Buf[Pos : len(Buf)]
	Pos      int    // read at &Buf[Pos], write at &Buf[len(Buf)]
	lastRead readOp // last read operation, so that Unread* can work correctly.
	BitPos   int
}

// mine - *read* offset - not advanced by write operations
func (p *Packet) Offset() int {
	return p.Pos
}

// The readOp constants describe the last action performed on
// the buffer, so that UnreadRune and UnreadByte can check for
// invalid usage. opReadRuneX constants are chosen such that
// converted to int they correspond to the rune size that was read.
type readOp int8

// Don't use iota for these, as the values need to correspond with the
// names and comments, which is easier to see when being explicit.
const (
	opRead      readOp = -1 // Any other read operation.
	opInvalid   readOp = 0  // Non-read operation.
	opReadRune1 readOp = 1  // Read rune of size 1.
	opReadRune2 readOp = 2  // Read rune of size 2.
	opReadRune3 readOp = 3  // Read rune of size 3.
	opReadRune4 readOp = 4  // Read rune of size 4.
)

// ErrTooLarge is passed to panic if memory cannot be allocated to store data in a buffer.
var ErrTooLarge = errors.New("bytes.Packet: too large")
var errNegativeRead = errors.New("bytes.Packet: reader returned negative count from Read")

const maxInt = int(^uint(0) >> 1)

// Bytes returns a slice of length b.Len() holding the unread portion of the buffer.
// The slice is valid for use only until the next buffer modification (that is,
// only until the next call to a method like [Buffer.Read], [Buffer.Write], [Packet.Reset], or [Buffer.Truncate]).
// The slice aliases the buffer content at least until the next buffer modification,
// so immediate changes to the slice will affect the result of future reads.
func (p *Packet) Bytes() []byte { return p.Buf[p.Pos:] }

// AvailableBuffer returns an empty buffer with b.Available() capacity.
// This buffer is intended to be appended to and
// passed to an immediately succeeding [Packet.Write] call.
// The buffer is only valid until the next write operation on b.
func (p *Packet) AvailableBuffer() []byte { return p.Buf[len(p.Buf):] }

// String returns the contents of the unread portion of the buffer
// as a string. If the [Packet] is a nil pointer, it returns "<nil>".
//
// To build strings more efficiently, see the [strings.Builder] type.
func (p *Packet) String() string {
	if p == nil {
		// Special case, useful in debugging.
		return "<nil>"
	}
	return string(p.Buf[p.Pos:])
}

// empty reports whether the unread portion of the buffer is empty.
func (p *Packet) empty() bool { return len(p.Buf) <= p.Pos }

// Len returns the number of bytes of the unread portion of the buffer;
// b.Len() == len(b.Bytes()).
func (p *Packet) Len() int { return len(p.Buf) - p.Pos }

// Cap returns the capacity of the buffer's underlying byte slice, that is, the
// total space allocated for the buffer's data.
func (p *Packet) Cap() int { return cap(p.Buf) }

// Available returns how many bytes are unused in the buffer.
func (p *Packet) Available() int { return cap(p.Buf) - len(p.Buf) }

// Truncate discards all but the first n unread bytes from the buffer
// but continues to use the same allocated storage.
// It panics if n is negative or greater than the length of the buffer.
func (p *Packet) Truncate(n int) {
	if n == 0 {
		p.Reset()
		return
	}
	p.lastRead = opInvalid
	if n < 0 || n > p.Len() {
		panic("bytes.Packet: truncation out of range")
	}
	p.Buf = p.Buf[:p.Pos+n]
}

// Reset resets the buffer to be empty,
// but it retains the underlying storage for use by future writes.
// Reset is the same as [Packet.Truncate](0).
func (p *Packet) Reset() {
	p.Buf = p.Buf[:0]
	p.Pos = 0
	p.lastRead = opInvalid
}

// tryGrowByReslice is an inlineable version of grow for the fast-case where the
// internal buffer only needs to be resliced.
// It returns the index where bytes should be written and whether it succeeded.
func (p *Packet) tryGrowByReslice(n int) (int, bool) {
	if l := len(p.Buf); n <= cap(p.Buf)-l {
		p.Buf = p.Buf[:l+n]
		return l, true
	}
	return 0, false
}

// grow grows the buffer to guarantee space for n more bytes.
// It returns the index where bytes should be written.
// If the buffer can't grow it will panic with ErrTooLarge.
func (p *Packet) grow(n int) int {
	m := p.Len()
	// If buffer is empty, reset to recover space.
	if m == 0 && p.Pos != 0 {
		p.Reset()
	}
	// Try to grow by means of a reslice.
	if i, ok := p.tryGrowByReslice(n); ok {
		return i
	}
	if p.Buf == nil && n <= smallBufferSize {
		p.Buf = make([]byte, n, smallBufferSize)
		return 0
	}
	c := cap(p.Buf)
	if n <= c/2-m {
		// We can slide things down instead of allocating a new
		// slice. We only need m+n <= c to slide, but
		// we instead let capacity get twice as large so we
		// don't spend all our time copying.
		copy(p.Buf, p.Buf[p.Pos:])
	} else if c > maxInt-c-n {
		panic(ErrTooLarge)
	} else {
		// Add p.Pos to account for p.Buf[:p.Pos] being sliced off the front.
		p.Buf = growSlice(p.Buf[p.Pos:], p.Pos+n)
	}
	// Restore p.Pos and len(p.Buf).
	p.Pos = 0
	p.Buf = p.Buf[:m+n]
	return m
}

// Grow grows the buffer's capacity, if necessary, to guarantee space for
// another n bytes. After Grow(n), at least n bytes can be written to the
// buffer without another allocation.
// If n is negative, Grow will panic.
// If the buffer can't grow it will panic with [ErrTooLarge].
func (p *Packet) Grow(n int) {
	if n < 0 {
		panic("bytes.Packet.Grow: negative count")
	}
	m := p.grow(n)
	p.Buf = p.Buf[:m]
}

// Write appends the contents of b to the buffer, growing the buffer as
// needed. The return value n is the length of b; err is always nil. If the
// buffer becomes too large, Write will panic with [ErrTooLarge].
func (p *Packet) Write(b []byte) (n int, err error) {
	p.lastRead = opInvalid
	m, ok := p.tryGrowByReslice(len(b))
	if !ok {
		m = p.grow(len(b))
	}
	return copy(p.Buf[m:], b), nil
}

// WriteString appends the contents of s to the buffer, growing the buffer as
// needed. The return value n is the length of s; err is always nil. If the
// buffer becomes too large, WriteString will panic with [ErrTooLarge].
func (p *Packet) WriteString(s string) (n int, err error) {
	p.lastRead = opInvalid
	m, ok := p.tryGrowByReslice(len(s))
	if !ok {
		m = p.grow(len(s))
	}
	return copy(p.Buf[m:], s), nil
}

// MinRead is the minimum slice size passed to a [Packet.Read] call by
// [Buffer.ReadFrom]. As long as the [Packet] has at least MinRead bytes beyond
// what is required to hold the contents of r, [Packet.ReadFrom] will not grow the
// underlying buffer.
const MinRead = 512

// ReadFrom reads data from r until EOF and appends it to the buffer, growing
// the buffer as needed. The return value n is the number of bytes read. Any
// error except io.EOF encountered during the read is also returned. If the
// buffer becomes too large, ReadFrom will panic with [ErrTooLarge].
func (p *Packet) ReadFrom(r io.Reader) (n int64, err error) {
	p.lastRead = opInvalid
	for {
		i := p.grow(MinRead)
		p.Buf = p.Buf[:i]
		m, e := r.Read(p.Buf[i:cap(p.Buf)])
		if m < 0 {
			panic(errNegativeRead)
		}

		p.Buf = p.Buf[:i+m]
		n += int64(m)
		if e == io.EOF {
			return n, nil // e is EOF, so return nil explicitly
		}
		if e != nil {
			return n, e
		}
	}
}

// growSlice grows b by n, preserving the original content of b.
// If the allocation fails, it panics with ErrTooLarge.
func growSlice(b []byte, n int) []byte {
	defer func() {
		if recover() != nil {
			panic(ErrTooLarge)
		}
	}()
	// TODO(http://golang.org/issue/51462): We should rely on the append-make
	// pattern so that the compiler can call runtime.growslice. For example:
	//	return append(b, make([]byte, n)...)
	// This avoids unnecessary zero-ing of the first len(b) bytes of the
	// allocated slice, but this pattern causes b to escape onto the heap.
	//
	// Instead use the append-make pattern with a nil slice to ensure that
	// we allocate buffers rounded up to the closest size class.
	c := len(b) + n // ensure enough space for n elements
	if c < 2*cap(b) {
		// The growth rate has historically always been 2x. In the future,
		// we could rely purely on append to determine the growth rate.
		c = 2 * cap(b)
	}
	b2 := append([]byte(nil), make([]byte, c)...)
	copy(b2, b)
	return b2[:len(b)]
}

// WriteTo writes data to w until the buffer is drained or an error occurs.
// The return value n is the number of bytes written; it always fits into an
// int, but it is int64 to match the [io.WriterTo] interface. Any error
// encountered during the write is also returned.
func (p *Packet) WriteTo(w io.Writer) (n int64, err error) {
	p.lastRead = opInvalid
	if nBytes := p.Len(); nBytes > 0 {
		m, e := w.Write(p.Buf[p.Pos:])
		if m > nBytes {
			panic("bytes.Packet.WriteTo: invalid Write count")
		}
		p.Pos += m
		n = int64(m)
		if e != nil {
			return n, e
		}
		// all bytes should have been written, by definition of
		// Write method in io.Writer
		if m != nBytes {
			return n, io.ErrShortWrite
		}
	}
	// Packet is now empty; reset.
	p.Reset()
	return n, nil
}

// WriteByte appends the byte c to the buffer, growing the buffer as needed.
// The returned error is always nil, but is included to match [bufio.Writer]'s
// WriteByte. If the buffer becomes too large, WriteByte will panic with
// [ErrTooLarge].
func (p *Packet) WriteByte(c byte) error {
	p.lastRead = opInvalid
	m, ok := p.tryGrowByReslice(1)
	if !ok {
		m = p.grow(1)
	}
	p.Buf[m] = c
	return nil
}

// WriteRune appends the UTF-8 encoding of Unicode code point r to the
// buffer, returning its length and an error, which is always nil but is
// included to match [bufio.Writer]'s WriteRune. The buffer is grown as needed;
// if it becomes too large, WriteRune will panic with [ErrTooLarge].
func (p *Packet) WriteRune(r rune) (n int, err error) {
	// Compare as uint32 to correctly handle negative runes.
	if uint32(r) < utf8.RuneSelf {
		p.WriteByte(byte(r))
		return 1, nil
	}
	p.lastRead = opInvalid
	m, ok := p.tryGrowByReslice(utf8.UTFMax)
	if !ok {
		m = p.grow(utf8.UTFMax)
	}
	p.Buf = utf8.AppendRune(p.Buf[:m], r)
	return len(p.Buf) - m, nil
}

// Read reads the next len(b) bytes from the buffer or until the buffer
// is drained. The return value n is the number of bytes read. If the
// buffer has no data to return, err is [io.EOF] (unless len(b) is zero);
// otherwise it is nil.
func (p *Packet) Read(b []byte) (n int, err error) {
	p.lastRead = opInvalid
	if p.empty() {
		// Packet is empty, reset to recover space.
		p.Reset()
		if len(b) == 0 {
			return 0, nil
		}
		return 0, io.EOF
	}
	n = copy(b, p.Buf[p.Pos:])
	p.Pos += n
	if n > 0 {
		p.lastRead = opRead
	}
	return n, nil
}

// Next returns a slice containing the next n bytes from the buffer,
// advancing the buffer as if the bytes had been returned by [Packet.Read].
// If there are fewer than n bytes in the buffer, Next returns the entire buffer.
// The slice is only valid until the next call to a read or write method.
func (p *Packet) Next(n int) []byte {
	p.lastRead = opInvalid
	m := p.Len()
	if n > m {
		n = m
	}
	data := p.Buf[p.Pos : p.Pos+n]
	p.Pos += n
	if n > 0 {
		p.lastRead = opRead
	}
	return data
}

// ReadByte reads and returns the next byte from the buffer.
// If no byte is available, it returns error [io.EOF].
func (p *Packet) ReadByte() (byte, error) {
	if p.empty() {
		// Packet is empty, reset to recover space.
		p.Reset()
		return 0, io.EOF
	}
	c := p.Buf[p.Pos]
	p.Pos++
	p.lastRead = opRead
	return c, nil
}

// ReadRune reads and returns the next UTF-8-encoded
// Unicode code point from the buffer.
// If no bytes are available, the error returned is io.EOF.
// If the bytes are an erroneous UTF-8 encoding, it
// consumes one byte and returns U+FFFD, 1.
func (p *Packet) ReadRune() (r rune, size int, err error) {
	if p.empty() {
		// Packet is empty, reset to recover space.
		p.Reset()
		return 0, 0, io.EOF
	}
	c := p.Buf[p.Pos]
	if c < utf8.RuneSelf {
		p.Pos++
		p.lastRead = opReadRune1
		return rune(c), 1, nil
	}
	r, n := utf8.DecodeRune(p.Buf[p.Pos:])
	p.Pos += n
	p.lastRead = readOp(n)
	return r, n, nil
}

// UnreadRune unreads the last rune returned by [Packet.ReadRune].
// If the most recent read or write operation on the buffer was
// not a successful [Packet.ReadRune], UnreadRune returns an error.  (In this regard
// it is stricter than [Packet.UnreadByte], which will unread the last byte
// from any read operation.)
func (p *Packet) UnreadRune() error {
	if p.lastRead <= opInvalid {
		return errors.New("bytes.Packet: UnreadRune: previous operation was not a successful ReadRune")
	}
	if p.Pos >= int(p.lastRead) {
		p.Pos -= int(p.lastRead)
	}
	p.lastRead = opInvalid
	return nil
}

var errUnreadByte = errors.New("bytes.Packet: UnreadByte: previous operation was not a successful read")

// UnreadByte unreads the last byte returned by the most recent successful
// read operation that read at least one byte. If a write has happened since
// the last read, if the last read returned an error, or if the read read zero
// bytes, UnreadByte returns an error.
func (p *Packet) UnreadByte() error {
	if p.lastRead == opInvalid {
		return errUnreadByte
	}
	p.lastRead = opInvalid
	if p.Pos > 0 {
		p.Pos--
	}
	return nil
}

// MINE
func (p *Packet) Rewind(n int) error {
	if p.Pos > 0 {
		p.Pos -= n
	}
	return nil
}

// ReadBytes reads until the first occurrence of delim in the input,
// returning a slice containing the data up to and including the delimiter.
// If ReadBytes encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often [io.EOF]).
// ReadBytes returns err != nil if and only if the returned data does not end in
// delim.
func (p *Packet) ReadBytes(delim byte) (line []byte, err error) {
	slice, err := p.readSlice(delim)
	// return a copy of slice. The buffer's backing array may
	// be overwritten by later calls.
	line = append(line, slice...)
	return line, err
}

// readSlice is like ReadBytes but returns a reference to internal buffer data.
func (p *Packet) readSlice(delim byte) (line []byte, err error) {
	i := IndexByte(p.Buf[p.Pos:], delim)
	end := p.Pos + i + 1
	if i < 0 {
		end = len(p.Buf)
		err = io.EOF
	}
	line = p.Buf[p.Pos:end]
	p.Pos = end
	p.lastRead = opRead
	return line, err
}

// mine - from internal/bytealg/indexbyte_generic.go - not optimized
func IndexByte(b []byte, c byte) int {
	for i, x := range b {
		if x == c {
			return i
		}
	}
	return -1
}

// ReadString reads until the first occurrence of delim in the input,
// returning a string containing the data up to and including the delimiter.
// If ReadString encounters an error before finding a delimiter,
// it returns the data read before the error and the error itself (often [io.EOF]).
// ReadString returns err != nil if and only if the returned data does not end
// in delim.
func (p *Packet) ReadString(delim byte) (line string, err error) {
	slice, err := p.readSlice(delim)
	return string(slice), err
}

// NewPacket creates and initializes a new [Packet] using Buf as its
// initial contents. The new [Packet] takes ownership of Buf, and the
// caller should not use Buf after this call. NewPacket is intended to
// prepare a [Packet] to read existing data. It can also be used to set
// the initial size of the internal buffer for writing. To do that,
// Buf should have the desired capacity but a length of zero.
//
// In most cases, new([Packet]) (or just declaring a [Packet] variable) is
// sufficient to initialize a [Packet].
func NewPacket(buf []byte) *Packet { return &Packet{Buf: buf} }

// NewPacketString creates and initializes a new [Packet] using string s as its
// initial contents. It is intended to prepare a buffer to read an existing
// string.
//
// In most cases, new([Packet]) (or just declaring a [Packet] variable) is
// sufficient to initialize a [Packet].
func NewPacketString(s string) *Packet {
	return &Packet{Buf: []byte(s)}
}
