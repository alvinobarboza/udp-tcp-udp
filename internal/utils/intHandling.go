package utils

import "encoding/binary"

func Int64ToByte(n uint64) []byte {
	bs := make([]byte, 8)
	binary.LittleEndian.PutUint64(bs, n)
	return bs
}

func Int32ToByte(n uint32) []byte {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, n)
	return bs
}

func Int16ToByte(n uint16) []byte {
	bs := make([]byte, 2)
	binary.LittleEndian.PutUint16(bs, n)
	return bs
}
