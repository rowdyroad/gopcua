// Copyright 2018 gopcua authors. All rights reserved.
// Use of this source code is governed by a MIT-style license that can be
// found in the LICENSE file.

package datatypes

import "encoding/binary"

// ExtensionObject is encoded as sequence of bytes prefixed by the NodeId of its DataTypeEncoding
// and the number of bytes encoded.
//
// Specification: Part 6, 5.2.2.15
type ExtensionObject struct {
	TypeID       *ExpandedNodeID
	EncodingMask byte
	Length       int32
	Body         *ByteString
}

// NewExtensionObject creates a new ExtensionObject.
func NewExtensionObject(typeID *ExpandedNodeID, mask uint8, body []byte) *ExtensionObject {
	e := &ExtensionObject{
		TypeID:       typeID,
		EncodingMask: mask,
		Body:         NewByteString(body),
	}
	e.SetLength()

	return e
}

// DecodeExtensionObject decodes given bytes into ExtensionObject.
func DecodeExtensionObject(b []byte) (*ExtensionObject, error) {
	e := &ExtensionObject{}
	if err := e.DecodeFromBytes(b); err != nil {
		return nil, err
	}
	return e, nil
}

// DecodeFromBytes decodes given bytes into ExtensionObject.
func (e *ExtensionObject) DecodeFromBytes(b []byte) error {
	// type id
	nodeID, err := DecodeExpandedNodeID(b)
	if err != nil {
		return err
	}
	e.TypeID = nodeID
	offset := e.TypeID.Len()

	// encoding mask
	e.EncodingMask = b[offset]
	offset++

	// length
	e.Length = int32(binary.LittleEndian.Uint32(b[offset : offset+4]))
	offset += 4

	// body
	e.Body = &ByteString{}
	if err := e.Body.DecodeFromBytes(b[offset:]); err != nil {
		return err
	}

	return nil
}

// Serialize serializes ExtensionObject into bytes.
func (e *ExtensionObject) Serialize() ([]byte, error) {
	b := make([]byte, e.Len())
	if err := e.SerializeTo(b); err != nil {
		return nil, err
	}

	return b, nil
}

// SerializeTo serializes ExtensionObject into bytes.
func (e *ExtensionObject) SerializeTo(b []byte) error {
	offset := 0

	// type id
	if e.TypeID != nil {
		if err := e.TypeID.SerializeTo(b); err != nil {
			return err
		}
		offset += e.TypeID.Len()
	}

	// encoding mask
	b[offset] = e.EncodingMask
	offset++

	// length
	binary.LittleEndian.PutUint32(b[offset:offset+4], uint32(e.Length))
	offset += 4

	// body
	if e.Body != nil {
		if err := e.Body.SerializeTo(b[offset:]); err != nil {
			return err
		}
	}

	return nil
}

// Len returns the actual length of ExtensionObject in int.
func (e *ExtensionObject) Len() int {
	// encoding mask byte + length
	length := 1 + 4

	if e.TypeID != nil {
		length += e.TypeID.Len()
	}

	if e.Body != nil {
		length += e.Body.Len()
	}

	return length
}

// SetLength sets the length of Body in Length field.
func (e *ExtensionObject) SetLength() {
	e.Length = int32(e.Body.Len())
}
