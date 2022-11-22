package pebble

import (
	"bytes"
	"encoding/binary"

	"github.com/hashicorp/go-msgpack/codec"
)

// Decode reverses the encode operation on a byte slice input
func decodeMsgPack(buf []byte, out interface{}) error {
	r := bytes.NewBuffer(buf)
	hd := codec.MsgpackHandle{}
	dec := codec.NewDecoder(r, &hd)
	return dec.Decode(out)
}

// Encode writes an encoded object to a new bytes buffer
func encodeMsgPack(in interface{}) (*bytes.Buffer, error) {
	buf := bytes.NewBuffer(nil)
	hd := codec.MsgpackHandle{}
	enc := codec.NewEncoder(buf, &hd)
	err := enc.Encode(in)
	return buf, err
}

// Converts bytes to an integer
func bytesToUint64(b []byte) uint64 {
	return binary.BigEndian.Uint64(b)
}

// Converts a uint to a byte slice
func uint64ToBytes(u uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, u)
	return buf
}

func genLogKey(idx uint64) []byte {
	key := make([]byte, 0, len(dbLogsPrefix))
	key = append(key, dbLogsPrefix...)
	key = append(key, uint64ToBytes(idx)...)
	return key
}

func genConfKey(key []byte) []byte {
	k := make([]byte, 0, len(dbConfPrefix)+len(key))
	k = append(k, dbConfPrefix...)
	k = append(k, key...)
	return k
}

func parseLogKey(key []byte) ([]byte, error) {
	if len(key) < len(dbLogsPrefix) {
		return nil, ErrKeyInvalid
	}
	key = key[len(dbLogsPrefix):]
	return key, nil
}
