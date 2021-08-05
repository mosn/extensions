package proxy

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
)

const (
	maxInt                = int(^uint(0) >> 1) // maxInt is max buffer length limit
	smallBufferSize       = 128                // smallBufferSize is an initial allocation minimal capacity.
	resetOffMark          = -1                 // buffer mark flag
	BigEndian       Order = 1                  // big endian byteOrder
	LittleEndian    Order = 2                  // little endian byteOrder
)

var BufferTooLarge = errors.New("wasm plugin Buffer: too large")
var IndexOutOfBound = errors.New("wasm plugin Buffer: index out of bound")

type ConfigMap Header

type Order int

type Header interface {
	// Get value of key
	// If multiple values associated with this key, first one will be returned.
	Get(key string) (string, bool)

	// Set key-value pair in header map, the previous pair will be replaced if exists
	Set(key, value string)

	// Add value for given key.
	// Multiple headers with the same key may be added with this function.
	// Use Set for setting a single header for the given key.
	Add(key, value string)

	// Del delete pair of specified key
	Del(key string)

	// Range calls f sequentially for each key and value present in the map.
	// If f returns false, range stops the iteration.
	Range(f func(key, value string) bool)

	// Size header key value pair count
	Size() int

	ToMap() map[string]string
}

type CommonHeader struct {
	m       map[string]string
	Changed bool
}

func NewConfigMap() ConfigMap {
	return &CommonHeader{}
}

func NewHeader() Header {
	return &CommonHeader{}
}

// Get value of key
func (h *CommonHeader) Get(key string) (value string, ok bool) {
	if len(h.m) == 0 {
		return "", false
	}
	value, ok = h.m[key]
	return
}

// Set key-value pair in header map, the previous pair will be replaced if exists
func (h *CommonHeader) Set(key string, value string) {
	h.Changed = true
	if len(h.m) == 0 {
		h.m = make(map[string]string, 8)
	}
	h.m[key] = value
}

// Add value for given key.
// Multiple headers with the same key may be added with this function.
// Use Set for setting a single header for the given key.
func (h *CommonHeader) Add(key string, value string) {
	panic("not supported")
}

// Del delete pair of specified key
func (h *CommonHeader) Del(key string) {
	h.Changed = true
	if len(h.m) == 0 {
		return
	}
	delete(h.m, key)
}

func (h *CommonHeader) Size() int {
	return len(h.m)
}

// Range calls f sequentially for each key and value present in the map.
// If f returns false, range stops the iteration.
func (h *CommonHeader) Range(f func(key, value string) bool) {
	if len(h.m) == 0 {
		return
	}

	for k, v := range h.m {
		// stop if f return false
		if !f(k, v) {
			break
		}
	}
}
func (h *CommonHeader) Clone() *CommonHeader {
	clone := &CommonHeader{}
	for k, v := range h.m {
		clone.Set(k, v)
	}
	return clone
}

func (h *CommonHeader) ToMap() map[string]string {
	if len(h.m) == 0 {
		h.m = make(map[string]string, 8)
	}
	return h.m
}

func (h *CommonHeader) String() string {
	if len(h.m) == 0 {
		return "{}"
	}
	buf := bytes.Buffer{}
	buf.WriteString("{")
	h.Range(func(key, value string) bool {
		buf.WriteString(key)
		buf.WriteString("=")
		buf.WriteString(value)
		return true
	})
	buf.WriteString("}")
	return buf.String()
}

type Buffer interface {
	Bytes() []byte
	Len() int
	Cap() int
	Pos() int
	Move(int)
	Grow(n int)
	Reset()
	Peek(n int) []byte
	Drain(len int)
	Mark()
	ResetMark()
	ByteOrder(order Order) Buffer

	WriteByte(value byte) error
	WriteUint16(value uint16) error
	WriteUint32(value uint32) error
	WriteUint(value uint) error
	WriteUint64(value uint64) error
	WriteInt16(value int16) error
	WriteInt32(value int32) error
	WriteInt(value int) error
	WriteInt64(value int64) error

	PutByte(index int, value byte) error
	PutUint16(index int, value uint16) error
	PutUint32(index int, value uint32) error
	PutUint(index int, value uint) error
	PutUint64(index int, value uint64) error
	PutInt16(index int, value int16) error
	PutInt32(index int, value int32) error
	PutInt(index int, value int) error
	PutInt64(index int, value int64) error

	Write(p []byte) (n int, err error)
	WriteString(s string) (n int, err error)

	ReadByte() (byte, error)
	ReadUint16() (uint16, error)
	ReadUint32() (uint32, error)
	ReadUint() (uint, error)
	ReadUint64() (uint64, error)
	ReadInt16() (int16, error)
	ReadInt32() (int32, error)
	ReadInt() (int, error)
	ReadInt64() (int64, error)

	GetByte(index int) (byte, error)
	GetUint16(index int) (uint16, error)
	GetUint32(index int) (uint32, error)
	GetUint(index int) (uint, error)
	GetUint64(index int) (uint64, error)
	GetInt16(index int) (int16, error)
	GetInt32(index int) (int32, error)
	GetInt(index int) (int, error)
	GetInt64(index int) (int64, error)
}

func AllocateBuffer() Buffer {
	return NewBuffer(smallBufferSize)
}

func NewBuffer(size int) Buffer {
	cap := size
	if cap < smallBufferSize {
		cap = smallBufferSize
	}
	return &byteBuffer{
		// be sure to update the index on write, where the length is set to 0
		buf:       make([]byte, 0, cap),
		pos:       0,
		mark:      resetOffMark,
		byteOrder: binary.BigEndian,
	}
}

func WrapBuffer(buf []byte) Buffer {
	return &byteBuffer{
		buf:       buf,
		pos:       0,
		mark:      resetOffMark,
		byteOrder: binary.BigEndian,
	}
}

type byteBuffer struct {
	buf       []byte           // contents are the bytes buf[pos : len(buf)]
	pos       int              // read at &buf[pos], write at &buf[len(buf)]
	mark      int              // mark flag
	byteOrder binary.ByteOrder // byte byteOrder
}

func (b *byteBuffer) Bytes() []byte {
	return b.buf[b.pos:]
}

func (b *byteBuffer) Len() int {
	return len(b.buf) - b.pos
}

func (b *byteBuffer) Cap() int {
	return cap(b.buf)
}

// Grow grows the buffer to guarantee space for n more bytes.
// If the buffer can't grow it will panic with ErrTooLarge.
func (b *byteBuffer) Grow(n int) {
	if n < 0 {
		panic("bytes.Buffer.Grow: negative count")
	}
	m := b.grow(n)
	b.buf = b.buf[:m]
}

func (b *byteBuffer) Reset() {
	b.buf = b.buf[:0]
	b.pos = 0
}

func (b *byteBuffer) ByteOrder(order Order) Buffer {
	switch order {
	case BigEndian:
		b.byteOrder = binary.BigEndian
	case LittleEndian:
		b.byteOrder = binary.LittleEndian
	default:
		b.byteOrder = binary.BigEndian
	}
	return b
}

func (b *byteBuffer) Peek(n int) []byte {
	if len(b.buf)-b.pos < n {
		return nil
	}

	return b.buf[b.pos : b.pos+n]
}

func (b *byteBuffer) Drain(n int) {
	if b.pos+n > len(b.buf) {
		// update pos to the end of b.buf
		b.pos = len(b.buf) - 1
		b.mark = resetOffMark
		return
	}

	b.pos += n
	b.mark = resetOffMark
}

func (b *byteBuffer) Mark() {
	b.mark = b.pos
}

func (b *byteBuffer) ResetMark() {
	if b.mark != resetOffMark {
		b.pos = b.mark
		b.mark = resetOffMark
	}
}

func (b *byteBuffer) WriteByte(value byte) error {
	m, ok := b.tryGrowBySlice(1)
	if !ok {
		m = b.grow(1)
	}

	b.buf[m] = value
	return nil
}

func (b *byteBuffer) WriteUint16(value uint16) error {
	m, ok := b.tryGrowBySlice(2)
	if !ok {
		m = b.grow(2)
	}

	b.byteOrder.PutUint16(b.buf[m:], value)
	return nil
}

func (b *byteBuffer) WriteUint32(value uint32) error {
	m, ok := b.tryGrowBySlice(4)
	if !ok {
		m = b.grow(4)
	}

	b.byteOrder.PutUint32(b.buf[m:], value)
	return nil
}

func (b *byteBuffer) WriteUint(value uint) error {
	return b.WriteUint32(uint32(value))
}

func (b *byteBuffer) WriteUint64(value uint64) error {
	m, ok := b.tryGrowBySlice(8)
	if !ok {
		m = b.grow(8)
	}

	b.byteOrder.PutUint64(b.buf[m:], value)
	return nil
}

func (b *byteBuffer) WriteInt16(value int16) error {
	return b.WriteUint16(uint16(value))
}

func (b *byteBuffer) WriteInt32(value int32) error {
	return b.WriteUint32(uint32(value))
}

func (b *byteBuffer) WriteInt(value int) error {
	return b.WriteUint32(uint32(value))
}

func (b *byteBuffer) WriteInt64(value int64) error {
	return b.WriteUint64(uint64(value))
}

func (b *byteBuffer) PutByte(i int, value byte) error {
	if i < 0 {
		panic("bytes.Buffer.PutByte: negative put index")
	}
	b.tryGrowSlice0(i, 1)
	b.buf[i] = value
	return nil
}

func (b *byteBuffer) PutUint16(i int, value uint16) error {
	if i < 0 {
		panic("bytes.Buffer.PutUint16: negative put index")
	}
	b.tryGrowSlice0(i, 2)
	b.byteOrder.PutUint16(b.buf[i:], value)
	return nil
}

func (b *byteBuffer) PutUint32(i int, value uint32) error {
	if i < 0 {
		panic("bytes.Buffer.PutUint32: negative put index")
	}
	b.tryGrowSlice0(i, 4)
	b.byteOrder.PutUint32(b.buf[i:], value)
	return nil
}

func (b *byteBuffer) PutUint(i int, value uint) error {
	if i < 0 {
		panic("bytes.Buffer.PutUint: negative put index")
	}
	b.tryGrowSlice0(i, 4)
	b.byteOrder.PutUint32(b.buf[i:], uint32(value))
	return nil
}

func (b *byteBuffer) PutUint64(i int, value uint64) error {
	if i < 0 {
		panic("bytes.Buffer.PutUint64: negative put index")
	}
	b.tryGrowSlice0(i, 8)
	b.byteOrder.PutUint64(b.buf[i:], value)
	return nil
}

func (b *byteBuffer) PutInt16(i int, value int16) error {
	if i < 0 {
		panic("bytes.Buffer.PutInt16: negative put index")
	}
	b.tryGrowSlice0(i, 2)
	b.byteOrder.PutUint16(b.buf[i:], uint16(value))
	return nil
}

func (b *byteBuffer) PutInt32(i int, value int32) error {
	if i < 0 {
		panic("bytes.Buffer.PutInt32: negative put index")
	}
	b.tryGrowSlice0(i, 4)
	b.byteOrder.PutUint32(b.buf[i:], uint32(value))
	return nil
}

func (b *byteBuffer) PutInt(i int, value int) error {
	if i < 0 {
		panic("bytes.Buffer.PutInt: negative put index")
	}
	b.tryGrowSlice0(i, 4)
	b.byteOrder.PutUint32(b.buf[i:], uint32(value))
	return nil
}

func (b *byteBuffer) PutInt64(i int, value int64) error {
	if i < 0 {
		panic("bytes.Buffer.PutInt64: negative put index")
	}
	b.tryGrowSlice0(i, 8)
	b.byteOrder.PutUint64(b.buf[i:], uint64(value))
	return nil
}

func (b *byteBuffer) Write(p []byte) (n int, err error) {
	m, ok := b.tryGrowBySlice(len(p))
	if !ok {
		m = b.grow(len(p))
	}

	return copy(b.buf[m:], p), nil
}

func (b *byteBuffer) WriteString(s string) (n int, err error) {
	m, ok := b.tryGrowBySlice(len(s))
	if !ok {
		m = b.grow(len(s))
	}
	return copy(b.buf[m:], s), nil
}

func (b *byteBuffer) ReadByte() (byte, error) {
	if b.pos >= len(b.buf) {
		return 0, io.EOF
	}
	v := b.buf[b.pos]
	b.pos++
	return v, nil
}

func (b *byteBuffer) ReadUint16() (uint16, error) {
	if b.pos+1 >= len(b.buf) {
		return 0, io.EOF
	}

	v := b.byteOrder.Uint16(b.buf[b.pos:])
	b.pos += 2
	return v, nil
}

func (b *byteBuffer) ReadUint32() (uint32, error) {
	if b.pos+3 >= len(b.buf) {
		return 0, io.EOF
	}

	v := b.byteOrder.Uint32(b.buf[b.pos:])
	b.pos += 4
	return v, nil
}

func (b *byteBuffer) ReadUint() (uint, error) {
	v, err := b.ReadUint32()
	if err != nil {
		return 0, err
	}
	return uint(v), nil
}

func (b *byteBuffer) ReadUint64() (uint64, error) {
	if b.pos+7 >= len(b.buf) {
		return 0, io.EOF
	}

	v := b.byteOrder.Uint64(b.buf[b.pos:])
	b.pos += 8
	return v, nil
}

func (b *byteBuffer) ReadInt16() (int16, error) {
	v, err := b.ReadUint16()
	if err != nil {
		return 0, err
	}
	return int16(v), nil
}

func (b *byteBuffer) ReadInt32() (int32, error) {
	v, err := b.ReadUint32()
	if err != nil {
		return 0, err
	}
	return int32(v), nil
}

func (b *byteBuffer) ReadInt() (int, error) {
	v, err := b.ReadUint32()
	if err != nil {
		return 0, err
	}
	return int(v), nil
}

func (b *byteBuffer) ReadInt64() (int64, error) {
	v, err := b.ReadUint64()
	if err != nil {
		return 0, err
	}
	return int64(v), nil
}

func (b *byteBuffer) GetByte(i int) (byte, error) {
	if i < 0 {
		panic("bytes.Buffer.GetByte: negative get index")
	}

	if i >= len(b.buf) {
		return 0, io.EOF
	}

	v := b.buf[i]
	return v, nil
}

func (b *byteBuffer) GetUint16(i int) (uint16, error) {
	if i < 0 {
		panic("bytes.Buffer.GetUint16: negative get index")
	}

	if i+1 >= len(b.buf) {
		return 0, io.EOF
	}

	v := b.byteOrder.Uint16(b.buf[i:])
	return v, nil
}

func (b *byteBuffer) GetUint32(i int) (uint32, error) {
	if i < 0 {
		panic("bytes.Buffer.GetUint32: negative get index")
	}

	if i+3 >= len(b.buf) {
		return 0, io.EOF
	}

	v := b.byteOrder.Uint32(b.buf[i:])
	return v, nil
}

func (b *byteBuffer) GetUint(i int) (uint, error) {
	if i < 0 {
		panic("bytes.Buffer.GetUint: negative get index")
	}

	if i+3 >= len(b.buf) {
		return 0, io.EOF
	}

	v := b.byteOrder.Uint32(b.buf[i:])
	return uint(v), nil
}

func (b *byteBuffer) GetUint64(i int) (uint64, error) {
	if i < 0 {
		panic("bytes.Buffer.GetUint64: negative get index")
	}

	if i+7 >= len(b.buf) {
		return 0, io.EOF
	}

	v := b.byteOrder.Uint64(b.buf[i:])
	return v, nil
}

func (b *byteBuffer) GetInt16(i int) (int16, error) {
	if i < 0 {
		panic("bytes.Buffer.GetInt16: negative get index")
	}

	if i+1 >= len(b.buf) {
		return 0, io.EOF
	}

	v := b.byteOrder.Uint16(b.buf[i:])
	return int16(v), nil
}

func (b *byteBuffer) GetInt32(i int) (int32, error) {
	if i < 0 {
		panic("bytes.Buffer.GetInt32: negative get index")
	}

	if i+3 >= len(b.buf) {
		return 0, io.EOF
	}

	v := b.byteOrder.Uint32(b.buf[i:])
	return int32(v), nil
}

func (b *byteBuffer) GetInt(i int) (int, error) {
	if i < 0 {
		panic("bytes.Buffer.GetInt: negative get index")
	}

	if i+3 >= len(b.buf) {
		return 0, io.EOF
	}

	v := b.byteOrder.Uint32(b.buf[i:])
	return int(v), nil
}

func (b *byteBuffer) GetInt64(i int) (int64, error) {
	if i < 0 {
		panic("bytes.Buffer.GetInt: negative get index")
	}

	if i+3 >= len(b.buf) {
		return 0, io.EOF
	}

	v := b.byteOrder.Uint64(b.buf[i:])
	return int64(v), nil
}

func (b *byteBuffer) Pos() int {
	return b.pos
}

func (b *byteBuffer) Move(p int) {
	if p < 0 || p > len(b.buf) {
		panic("bytes.Buffer.Move: bad index")
	}
	b.pos = p
}

// ======================== private method impl ========================

// empty reports whether the unread portion of the buffer is empty.
func (b *byteBuffer) empty() bool { return b.pos >= len(b.buf) }

func (b *byteBuffer) grow(n int) int {
	m := b.Len()

	// If buffer is empty, reset to recover space.
	if m == 0 && b.pos != 0 {
		b.Reset()
	}

	// Try to grow by means of a re-slice.
	if i, ok := b.tryGrowBySlice(n); ok {
		return i
	}
	if b.buf == nil && n <= smallBufferSize {
		b.buf = make([]byte, n, smallBufferSize)
		return 0
	}

	c := cap(b.buf)
	if n <= c/2-m {
		// The current position can be moved.
		if b.pos > 0 {
			// We can slide things down instead of allocating a new
			// slice. We only need m+n <= c to slide, but
			// we instead let capacity get twice as large so we
			// don't spend all our time copying.
			copy(b.buf, b.buf[b.pos:])
		}
	} else if c > maxInt-c-n {
		panic(BufferTooLarge)
	} else {
		// Not enough space anywhere, we need to allocate.
		//log.Infof("allocate len %v, cap %v", m +n, 2*c+n)
		buf := make([]byte, m+n, 2*c+n)
		//log.Infof("allocate len %v, cap %v, success", m +n, 2*c+n)
		copy(buf, b.buf[b.pos:])
		b.buf = buf
	}
	// Restore b.pos and len(b.buf).
	b.pos = 0
	b.buf = b.buf[:m+n]
	return m
}

func (b *byteBuffer) tryGrowBySlice(n int) (int, bool) {
	if l := len(b.buf); n <= cap(b.buf)-l {
		b.buf = b.buf[:l+n]
		return l, true
	}
	return 0, false
}

func (b *byteBuffer) tryGrowSlice0(i, l int) {
	n := i + l - len(b.buf)
	if n > 0 {
		_, ok := b.tryGrowBySlice(n)
		if !ok {
			b.grow(n)
		}
	}
}
