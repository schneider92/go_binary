package binary

// BinaryWriter is a helper struct to write binary data to a buffer
type BinaryWriter struct {
	buf []byte
}

// Create a binary writer that will append the data to the end of the given buffer. The buffer is
// not copied, so any changes to the buffer will be reflected in the writer.
func NewBinaryWriterFromBuffer(buf []byte) *BinaryWriter {
	return &BinaryWriter{
		buf: buf,
	}
}

// Create a binary writer that will write to a new buffer with the given capacity.
func NewBinaryWriter(capacity int) *BinaryWriter {
	return &BinaryWriter{
		buf: make([]byte, 0, capacity),
	}
}

// Get the underlying buffer.
func (bw *BinaryWriter) Bytes() []byte {
	return bw.buf
}

// Write a 64-bit unsigned value as varint
func (bw *BinaryWriter) WriteUnsignedVarint(val uint64) {
	// save buffer length
	origlen := len(bw.buf)

	// write by byte until it becomes 0
	first := true
	for {
		b0 := byte(val & 0x7f)
		if !first {
			b0 |= 0x80
		}
		bw.buf = append(bw.buf, b0)
		first = false
		val >>= 7
		if val == 0 {
			break
		}
	}

	// reverse the bytes to get the correct order
	for i, j := origlen, len(bw.buf)-1; i < j; i, j = i+1, j-1 {
		bw.buf[i], bw.buf[j] = bw.buf[j], bw.buf[i]
	}
}

// Write a 64-bit signed value as varint using zigzag encoding
func (bw *BinaryWriter) WriteSignedVarint(val int64) {
	// zigzag encode
	var uval uint64
	if val >= 0 {
		uval = uint64(val) * 2
	} else {
		uval = uint64(-val*2 - 1)
	}
	bw.WriteUnsignedVarint(uval)
}

// Write a 8-bit unsigned value to the buffer
func (bw *BinaryWriter) WriteUint8(val uint8) {
	bw.buf = append(bw.buf, val)
}

// Write a 16-bit unsigned value to the buffer in big-endian order
func (bw *BinaryWriter) WriteUint16(val uint16) {
	bw.buf = append(bw.buf, byte(val>>8), byte(val))
}

// Write a 32-bit unsigned value to the buffer in big-endian order
func (bw *BinaryWriter) WriteUint32(val uint32) {
	bw.buf = append(bw.buf, byte(val>>24), byte(val>>16), byte(val>>8), byte(val))
}

// Write a 64-bit unsigned value to the buffer in big-endian order
func (bw *BinaryWriter) WriteUint64(val uint64) {
	bw.buf = append(bw.buf,
		byte(val>>56), byte(val>>48), byte(val>>40), byte(val>>32),
		byte(val>>24), byte(val>>16), byte(val>>8), byte(val))
}

// Write a 8-bit 2's complement signed value to the buffer
func (bw *BinaryWriter) WriteInt8(val int8) {
	bw.WriteUint8(uint8(val))
}

// Write a 16-bit 2's complement signed value to the buffer in big-endian order
func (bw *BinaryWriter) WriteInt16(val int16) {
	bw.WriteUint16(uint16(val))
}

// Write a 32-bit 2's complement signed value to the buffer in big-endian order
func (bw *BinaryWriter) WriteInt32(val int32) {
	bw.WriteUint32(uint32(val))
}

// Write a 64-bit 2's complement signed value to the buffer in big-endian order
func (bw *BinaryWriter) WriteInt64(val int64) {
	bw.WriteUint64(uint64(val))
}

// Write a byte slice without any length information
func (bw *BinaryWriter) WriteRawBytes(data []byte) {
	bw.buf = append(bw.buf, data...)
}

// Write a byte slice to the buffer, prefixed with its length as a varint
func (bw *BinaryWriter) WriteBytes(data []byte) {
	bw.WriteUnsignedVarint(uint64(len(data)))
	bw.buf = append(bw.buf, data...)
}

// Write a string to the buffer, prefixed with its length as a varint
func (bw *BinaryWriter) WriteString(data string) {
	bw.WriteUnsignedVarint(uint64(len(data)))
	bw.buf = append(bw.buf, data...)
}

// BinaryParser is a helper struct to read binary data from a buffer
type BinaryParser struct {
	buf []byte
}

// Create a binary parser that will read from the given buffer. The buffer is not copied, so any
// changes to the buffer will be reflected in the parser.
func NewBinaryParser(buf []byte) *BinaryParser {
	return &BinaryParser{
		buf: buf,
	}
}

// Read a 64-bit unsigned value as varint. Returns the value and a boolean indicating success.
func (bp *BinaryParser) ReadUnsignedVarint() (uint64, bool) {
	var val uint64
	i := 0
	l := len(bp.buf)
	for {
		if i >= l {
			bp.buf = nil
			return 0, false
		}
		b := bp.buf[i]
		val = (val << 7) | uint64(b&0x7f) // TODO overflow check
		if b&0x80 == 0 {
			bp.buf = bp.buf[i+1:]
			return val, true
		}
		i++
	}
}

// Read a 64-bit signed value as varint using zigzag decoding. Returns the value and a boolean
// indicating success.
func (bp *BinaryParser) ReadSignedVarint() (int64, bool) {
	uval, ok := bp.ReadUnsignedVarint()
	if !ok {
		bp.buf = nil
		return 0, false
	}

	// zigzag decode
	var val int64
	if uval%2 == 0 {
		val = int64(uval / 2)
	} else {
		val = -int64(uval/2) - 1
	}
	return val, true
}

// Read a 8-bit unsigned value from the buffer. Returns the value and a boolean indicating success.
func (bp *BinaryParser) ReadUint8() (uint8, bool) {
	if len(bp.buf) < 1 {
		bp.buf = nil
		return 0, false
	}
	val := bp.buf[0]
	bp.buf = bp.buf[1:]
	return val, true
}

// Read a 16-bit unsigned value from the buffer in big-endian order. Returns the value and a boolean
// indicating success.
func (bp *BinaryParser) ReadUint16() (uint16, bool) {
	if len(bp.buf) < 2 {
		bp.buf = nil
		return 0, false
	}
	val := uint16(bp.buf[0])<<8 | uint16(bp.buf[1])
	bp.buf = bp.buf[2:]
	return val, true
}

// Read a 32-bit unsigned value from the buffer in big-endian order. Returns the value and a boolean
// indicating success.
func (bp *BinaryParser) ReadUint32() (uint32, bool) {
	if len(bp.buf) < 4 {
		bp.buf = nil
		return 0, false
	}
	val := uint32(bp.buf[0])<<24 | uint32(bp.buf[1])<<16 | uint32(bp.buf[2])<<8 | uint32(bp.buf[3])
	bp.buf = bp.buf[4:]
	return val, true
}

// Read a 64-bit unsigned value from the buffer in big-endian order. Returns the value and a boolean
// indicating success.
func (bp *BinaryParser) ReadUint64() (uint64, bool) {
	if len(bp.buf) < 8 {
		bp.buf = nil
		return 0, false
	}
	val := uint64(bp.buf[0])<<56 | uint64(bp.buf[1])<<48 | uint64(bp.buf[2])<<40 | uint64(bp.buf[3])<<32 |
		uint64(bp.buf[4])<<24 | uint64(bp.buf[5])<<16 | uint64(bp.buf[6])<<8 | uint64(bp.buf[7])
	bp.buf = bp.buf[8:]
	return val, true
}

// Read a 8-bit 2's complement signed value from the buffer. Returns the value and a boolean
// indicating success.
func (bp *BinaryParser) ReadInt8() (int8, bool) {
	val, ok := bp.ReadUint8()
	return int8(val), ok
}

// Read a 16-bit 2's complement signed value from the buffer in big-endian order. Returns the value
// and a boolean indicating success.
func (bp *BinaryParser) ReadInt16() (int16, bool) {
	val, ok := bp.ReadUint16()
	return int16(val), ok
}

// Read a 32-bit 2's complement signed value from the buffer in big-endian order. Returns the value
// and a boolean indicating success.
func (bp *BinaryParser) ReadInt32() (int32, bool) {
	val, ok := bp.ReadUint32()
	return int32(val), ok
}

// Read a 64-bit 2's complement signed value from the buffer in big-endian order. Returns the value
// and a boolean indicating success.
func (bp *BinaryParser) ReadInt64() (int64, bool) {
	val, ok := bp.ReadUint64()
	return int64(val), ok
}

// Read a byte slice of the given length from the buffer. Returns the byte slice and a boolean
// indicating success.
func (bp *BinaryParser) ReadRawBytes(length int) ([]byte, bool) {
	if len(bp.buf) < length {
		bp.buf = nil
		return nil, false
	}
	data := bp.buf[:length]
	bp.buf = bp.buf[length:]
	return data, true
}

// Read a byte slice from the buffer, prefixed with its length as a varint. Returns the byte slice
// and a boolean indicating success.
func (bp *BinaryParser) ReadBytes() ([]byte, bool) {
	length, ok := bp.ReadUnsignedVarint()
	if !ok {
		bp.buf = nil
		return nil, false
	}
	if uint64(len(bp.buf)) < length {
		bp.buf = nil
		return nil, false
	}
	data := bp.buf[:length]
	bp.buf = bp.buf[length:]
	return data, true
}

// Read a string from the buffer, prefixed with its length as a varint. Returns the string and a
// boolean indicating success.
func (bp *BinaryParser) ReadString() (string, bool) {
	data, ok := bp.ReadBytes()
	if !ok {
		bp.buf = nil
		return "", false
	}
	return string(data), true
}

// Get the remaining bytes in the buffer. This can be used to check if there are any bytes left to read.
func (bp *BinaryParser) Remaining() []byte {
	return bp.buf
}

// Success flag of the read operations can be ignored, then the following functions return the
// zero value of the type and the success flag can be checked by calling IsOK() after all the reads
// are done. If any read operation fails, the buffer will be set to nil and IsOK() will return false.
func (bp *BinaryParser) IsOK() bool {
	return bp.buf != nil
}
