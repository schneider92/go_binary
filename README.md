# Summary
This small project is a binary serializer written in go.

# Varint
Variable-length integer format stores 64-bit integers in that many bytes that is needed for its actual value.

### Format
Varints are stored in 7 bits per byte, because the most significant bit is used as continuation bit. This bit always has a value of 1 unless the current byte is the last byte of the value.

Examples:
```
Values below 128 does not change
0   (0x00) -> 0000 0000
1   (0x01) -> 0000 0001
127 (0x7F) -> 0111 1111

128 (0x80) -> 1000 0001  0000 0000
       ^          ^
       --------------- continuation bits
       x000 0001  x000 0000 - this is the actual data

54321 (0xD431)
    in binary:    1101 0100 0011 0001
    grouped:      x000 0011  x010 1000  x011 0001
    cont bits:    1000 0011  1010 1000  0011 0001
    so the final value is this: 0x83A831
```

# Signed varint values
Signed values are first converted to unsigned ones using zigzag encoding, so the small negative values still occupy a little space as varint:
```
-3  =>  5
-2  =>  3
-1  =>  1
 0  =>  0
 1  =>  2
 2  =>  4
 3  =>  6
```

# Strings and byte slices
String and byte slice values are stored byte by byte, but they are prepended with the length encoded in varint:
```
"Hello" => 5, 'H', 'e', 'l', 'l', 'o'
```

# Fixed length
Fixed length integer values are also supported in 8, 16, 32, and 64 bits, both signed and unsigned. These values are stored most significant byte first order, and negative values are stored as 2s complement.

# TODOs
Currently the overflow condition in varint parsing is not handled, but if the value comes from an untrusted source, it is technically possible. Also it might be better to disallow values where a lot of zero bytes with continuation bits (0x80) can be prepended to a valid value.
