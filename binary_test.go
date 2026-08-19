package binary

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnsignedVarint(t *testing.T) {
	wr := NewBinaryWriter(64)

	wr.WriteUnsignedVarint(0)
	assert.Equal(t, []byte{0x00}, wr.Bytes())

	wr.WriteUnsignedVarint(1)
	assert.Equal(t, []byte{0x00, 0x01}, wr.Bytes())

	wr.WriteUnsignedVarint(127)
	assert.Equal(t, []byte{0x00, 0x01, 0x7f}, wr.Bytes())

	wr.WriteUnsignedVarint(128)
	assert.Equal(t, []byte{0x00, 0x01, 0x7f, 0x81, 0x00}, wr.Bytes())

	rd := NewBinaryParser(wr.Bytes())

	v, ok := rd.ReadUnsignedVarint()
	require.True(t, ok)
	require.Equal(t, uint64(0), v)
	require.True(t, rd.IsOK())

	v, ok = rd.ReadUnsignedVarint()
	require.True(t, ok)
	require.Equal(t, uint64(1), v)
	require.True(t, rd.IsOK())

	v, ok = rd.ReadUnsignedVarint()
	require.True(t, ok)
	require.Equal(t, uint64(127), v)
	require.True(t, rd.IsOK())

	v, ok = rd.ReadUnsignedVarint()
	require.True(t, ok)
	require.Equal(t, uint64(128), v)
	require.True(t, rd.IsOK())

	// read one more when empty
	require.Equal(t, []byte{}, rd.Remaining())
	v, ok = rd.ReadUnsignedVarint()
	require.False(t, ok)
	require.Equal(t, uint64(0), v)
	require.False(t, rd.IsOK())
}

func TestUnsignedVarint2(t *testing.T) {
	wr := NewBinaryWriter(64)

	// 123456789 = 75BCD15
	// 0111 0101 1011 1100 1101 0001 0101
	// 1'011_1010 1'110_1111 1'001_1010 0'001_0101
	// BA EF 9A 15
	wr.WriteUnsignedVarint(123456789)
	assert.Equal(t, []byte{0xBA, 0xEF, 0x9A, 0x15}, wr.Bytes())

	rd := NewBinaryParser(wr.Bytes())
	v, ok := rd.ReadUnsignedVarint()
	require.True(t, ok)
	require.Equal(t, uint64(123456789), v)
	require.True(t, rd.IsOK())

	// can have a lot of leading bytes with no consequence
	rd = NewBinaryParser([]byte{0x80, 0x80, 0x80, 0x80, 11})
	ui, ok := rd.ReadUnsignedVarint()
	require.True(t, ok)
	require.Equal(t, uint64(11), ui)
	require.True(t, rd.IsOK())
	require.Equal(t, []byte{}, rd.Remaining())
}

func TestSignedVarint(t *testing.T) {
	wr := NewBinaryWriter(64)

	wr.WriteSignedVarint(0)
	wr.WriteSignedVarint(1)
	wr.WriteSignedVarint(2)
	wr.WriteSignedVarint(3)
	wr.WriteSignedVarint(-1)
	wr.WriteSignedVarint(-2)
	wr.WriteSignedVarint(-3)

	bytes := wr.Bytes()

	// read them back as signed varints
	rd := NewBinaryParser(bytes)
	sv, ok := rd.ReadSignedVarint()
	require.True(t, ok)
	require.Equal(t, int64(0), sv)

	sv, ok = rd.ReadSignedVarint()
	require.True(t, ok)
	require.Equal(t, int64(1), sv)

	sv, ok = rd.ReadSignedVarint()
	require.True(t, ok)
	require.Equal(t, int64(2), sv)

	sv, ok = rd.ReadSignedVarint()
	require.True(t, ok)
	require.Equal(t, int64(3), sv)

	sv, ok = rd.ReadSignedVarint()
	require.True(t, ok)
	require.Equal(t, int64(-1), sv)

	sv, ok = rd.ReadSignedVarint()
	require.True(t, ok)
	require.Equal(t, int64(-2), sv)

	sv, ok = rd.ReadSignedVarint()
	require.True(t, ok)
	require.Equal(t, int64(-3), sv)

	sv, ok = rd.ReadSignedVarint()
	require.False(t, ok)
	require.Equal(t, int64(0), sv)

	// now read them back as unsigned varints
	rd = NewBinaryParser(bytes)
	uv, ok := rd.ReadUnsignedVarint()
	require.True(t, ok)
	require.Equal(t, uint64(0), uv)

	uv, ok = rd.ReadUnsignedVarint()
	require.True(t, ok)
	require.Equal(t, uint64(2), uv)

	uv, ok = rd.ReadUnsignedVarint()
	require.True(t, ok)
	require.Equal(t, uint64(4), uv)

	uv, ok = rd.ReadUnsignedVarint()
	require.True(t, ok)
	require.Equal(t, uint64(6), uv)

	uv, ok = rd.ReadUnsignedVarint()
	require.True(t, ok)
	require.Equal(t, uint64(1), uv)

	uv, ok = rd.ReadUnsignedVarint()
	require.True(t, ok)
	require.Equal(t, uint64(3), uv)

	uv, ok = rd.ReadUnsignedVarint()
	require.True(t, ok)
	require.Equal(t, uint64(5), uv)

	uv, ok = rd.ReadUnsignedVarint()
	require.False(t, ok)
	require.Equal(t, uint64(0), uv)
}

func TestFixedLengthIntegers(t *testing.T) {
	wr := NewBinaryWriter(64)

	wr.WriteUint8(0x12)
	wr.WriteUint16(0x1234)
	wr.WriteUint32(0x12345678)
	wr.WriteUint64(0x1234567890abcdef)

	wr.WriteInt8(-1)
	wr.WriteInt16(-2)
	wr.WriteInt32(-3)
	wr.WriteInt64(-4)

	bytes := wr.Bytes()
	require.Equal(t, []byte{
		0x12,
		0x12, 0x34,
		0x12, 0x34, 0x56, 0x78,
		0x12, 0x34, 0x56, 0x78, 0x90, 0xab, 0xcd, 0xef,
		0xff,
		0xff, 0xfe,
		0xff, 0xff, 0xff, 0xfd,
		0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xfc,
	}, bytes)

	rd := NewBinaryParser(bytes)
	u8, ok := rd.ReadUint8()
	require.True(t, ok)
	require.Equal(t, uint8(0x12), u8)

	u16, ok := rd.ReadUint16()
	require.True(t, ok)
	require.Equal(t, uint16(0x1234), u16)

	u32, ok := rd.ReadUint32()
	require.True(t, ok)
	require.Equal(t, uint32(0x12345678), u32)

	u64, ok := rd.ReadUint64()
	require.True(t, ok)
	require.Equal(t, uint64(0x1234567890abcdef), u64)

	i8, ok := rd.ReadInt8()
	require.True(t, ok)
	require.Equal(t, int8(-1), i8)

	i16, ok := rd.ReadInt16()
	require.True(t, ok)
	require.Equal(t, int16(-2), i16)

	i32, ok := rd.ReadInt32()
	require.True(t, ok)
	require.Equal(t, int32(-3), i32)

	i64, ok := rd.ReadInt64()
	require.True(t, ok)
	require.Equal(t, int64(-4), i64)

	// read more when empty
	require.Equal(t, []byte{}, rd.Remaining())
	u8, ok = rd.ReadUint8()
	require.False(t, ok)
	require.Equal(t, uint8(0), u8)

	u16, ok = rd.ReadUint16()
	require.False(t, ok)
	require.Equal(t, uint16(0), u16)

	u32, ok = rd.ReadUint32()
	require.False(t, ok)
	require.Equal(t, uint32(0), u32)

	u64, ok = rd.ReadUint64()
	require.False(t, ok)
	require.Equal(t, uint64(0), u64)
}

func TestBytesAndString(t *testing.T) {
	wr := NewBinaryWriterFromBuffer([]byte{5, 'h', 'e', 'l', 'l', 'o'})
	wr.WriteBytes([]byte{})
	wr.WriteBytes([]byte{'1', '2'})
	wr.WriteString("")
	wr.WriteString("world")

	rd := NewBinaryParser(wr.Bytes())
	s, ok := rd.ReadString()
	require.True(t, ok)
	require.Equal(t, "hello", s)

	b, ok := rd.ReadBytes()
	require.True(t, ok)
	require.Equal(t, []byte{}, b)

	s, ok = rd.ReadString()
	require.True(t, ok)
	require.Equal(t, "12", s)

	b, ok = rd.ReadBytes()
	require.True(t, ok)
	require.Equal(t, []byte{}, b)

	s, ok = rd.ReadString()
	require.True(t, ok)
	require.Equal(t, "world", s)

	s, ok = rd.ReadString()
	require.False(t, ok)
	require.Equal(t, "", s)

	b, ok = rd.ReadBytes()
	require.False(t, ok)
	require.Equal(t, []byte(nil), b)
}

func TestRawBytes(t *testing.T) {
	wr := NewBinaryWriterFromBuffer([]byte("hello"))
	wr.WriteBytes([]byte{' '})
	wr.WriteUint8(5)
	wr.WriteRawBytes([]byte("world"))

	bytes := wr.Bytes()
	require.Equal(t, []byte{'h', 'e', 'l', 'l', 'o', 1, ' ', 5, 'w', 'o', 'r', 'l', 'd'}, bytes)

	rd := NewBinaryParser(bytes)
	b, ok := rd.ReadRawBytes(5)
	require.True(t, ok)
	require.Equal(t, []byte("hello"), b)

	b, ok = rd.ReadRawBytes(2)
	require.True(t, ok)
	require.Equal(t, []byte{1, ' '}, b)

	b, ok = rd.ReadBytes()
	require.True(t, ok)
	require.Equal(t, []byte("world"), b)

	b, ok = rd.ReadRawBytes(0)
	require.True(t, ok)
	require.Equal(t, []byte{}, b)

	b = rd.Remaining()
	require.Equal(t, []byte{}, b)
	require.True(t, rd.IsOK())

	b, ok = rd.ReadRawBytes(1)
	require.False(t, ok)
	require.Equal(t, []byte(nil), b)
	require.False(t, rd.IsOK())
}

func TestNotEnoughData(t *testing.T) {
	bytes := []byte{5, 'h', 'e', 'l', 'l', 'o', 2}

	// the leading 5 can be read as varint
	rd := NewBinaryParser(bytes)
	ui, ok := rd.ReadUnsignedVarint()
	require.True(t, ok)
	require.Equal(t, uint64(5), ui)
	require.True(t, rd.IsOK())

	// if the input length is not enough, it is an error
	rd = NewBinaryParser(bytes[:5])
	s, ok := rd.ReadString()
	require.False(t, ok)
	require.Equal(t, "", s)
	require.False(t, rd.IsOK())

	// varint is not terminated
	bytes = []byte{0x80, 0x80, 0x81}
	rd = NewBinaryParser(bytes)
	ui, ok = rd.ReadUnsignedVarint()
	require.False(t, ok)
	require.Equal(t, uint64(0), ui)
	require.False(t, rd.IsOK())
}
