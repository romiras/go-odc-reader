// Package oberon defines primitive types used in the Oberon system.
// These types map Oberon/BlackBox types to Go equivalents.
package oberon

// Boolean represents a 1-byte boolean value (0 = FALSE, 1 = TRUE)
type Boolean = bool

// ShortChar represents a 1-byte character in the Latin-1 character set (Unicode page 0)
type ShortChar = byte

// Char represents a 2-byte character in the Unicode character set (0000X..0FFFFX)
// Note: Using uint16 as Go's rune is int32
type Char = uint16

// Byte represents a 1-byte signed integer (-128..127)
type Byte = int8

// ShortInt represents a 2-byte signed integer (-32768..32767)
type ShortInt = int16

// Integer represents a 4-byte signed integer (-2147483648..2147483647)
type Integer = int32

// LongInt represents an 8-byte signed integer (-9223372036854775808..9223372036854775807)
type LongInt = int64

// ShortReal represents a 4-byte IEEE 754 floating point number
type ShortReal = float32

// Real represents an 8-byte IEEE 754 floating point number
type Real = float64

// Set represents a 4-byte set (least significant bit = element 0)
type Set = uint32
