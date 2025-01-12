package types

import (
	"bytes"
	"encoding/binary"
)

const u32Size uint32 = 4

// MolSerializer molecule serialize interface
type MolSerializer interface {
	Serialize() ([]byte, error)
}

// serializeUint32 serialize uint32
func serializeUint32(n uint32) []byte {
	b := make([]byte, 4)
	binary.LittleEndian.PutUint32(b, n)

	return b
}

// SerializeArray serialize array
func SerializeArray(items []MolSerializer) ([][]byte, error) {
	ret := make([][]byte, len(items))
	for i := 0; i < len(items); i++ {
		bytes, err := items[i].Serialize()
		if err != nil {
			return nil, err
		}

		ret[i] = bytes
	}

	return ret, nil
}

// SerializeStruct serialize struct
func SerializeStruct(fields [][]byte) []byte {
	b := new(bytes.Buffer)

	for i := 0; i < len(fields); i++ {
		b.Write(fields[i])
	}

	return b.Bytes()
}

// SerializeFixVec serialize fixvec vector
/*
 * There are two steps of serializing a fixvec:
 *
 *     Serialize the length as a 32 bit unsigned integer in little-endian.
 *     Serialize all items in it.
 */
func SerializeFixVec(items [][]byte) []byte {
	// Empty fix vector bytes
	if len(items) == 0 {
		return []byte{00, 00, 00, 00}
	}

	l := serializeUint32(uint32(len(items)))

	b := new(bytes.Buffer)

	b.Write(l)

	for i := 0; i < len(items); i++ {
		b.Write(items[i])
	}

	return b.Bytes()
}

// SerializeDynVec serialize dynvec
/*
 * There are three steps of serializing a dynvec:
 *
 *     Serialize the full size in bytes as a 32 bit unsigned integer in little-endian.
 *     Serialize all offset of items as 32 bit unsigned integer in little-endian.
 *     Serialize all items in it.
 */
func SerializeDynVec(items [][]byte) []byte {
	// Start with u32Size
	size := u32Size

	// Empty dyn vector, just return size's bytes
	if len(items) == 0 {
		return serializeUint32(size)
	}

	offsets := make([]uint32, len(items))

	// Calculate first offset then loop for rest items offsets
	offsets[0] = size + u32Size*uint32(len(items))
	for i := 0; i < len(items); i++ {
		size += u32Size + uint32(len(items[i]))

		if i != 0 {
			offsets[i] = offsets[i-1] + uint32(len(items[i-1]))
		}
	}

	b := new(bytes.Buffer)

	b.Write(serializeUint32(size))

	for i := 0; i < len(items); i++ {
		b.Write(serializeUint32(offsets[i]))
	}

	for i := 0; i < len(items); i++ {
		b.Write(items[i])
	}

	return b.Bytes()
}

// SerializeTable serialize table
/*
 * The serializing steps are same as dynvec:
 *
 *     Serialize the full size in bytes as a 32 bit unsigned integer in little-endian.
 *     Serialize all offset of fields as 32 bit unsigned integer in little-endian.
 *     Serialize all fields in it in the order they are declared.
 */
func SerializeTable(fields [][]byte) []byte {
	size := u32Size
	offsets := make([]uint32, len(fields))

	// Calculate first offset then loop for rest items offsets
	offsets[0] = u32Size + u32Size*uint32(len(fields))
	for i := 0; i < len(fields); i++ {
		size += u32Size + uint32(len(fields[i]))

		if i != 0 {
			offsets[i] = offsets[i-1] + uint32(len(fields[i-1]))
		}
	}

	b := new(bytes.Buffer)

	b.Write(serializeUint32(size))

	for i := 0; i < len(fields); i++ {
		b.Write(serializeUint32(offsets[i]))
	}

	for i := 0; i < len(fields); i++ {
		b.Write(fields[i])
	}

	return b.Bytes()
}

// SerializeOption serialize option
func SerializeOption(o MolSerializer) ([]byte, error) {
	if o == nil {
		return []byte{}, nil
	}

	return o.Serialize()
}
