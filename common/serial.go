package common

import (
	"bytes"
	"encoding/gob"
)

// Serializable is an interface for types that can
// be encoded to and decoded from a stream of bytes.
type Serializable interface {
	// Serialize converts the object into a stream of bytes
	Serialize() ([]byte, error)

	// Deserialize converts a stream of bytes into the object
	Deserialize([]byte) error
}

// GobEncode encodes an object into a gob encoded stream of bytes.
// Returns an error if the gob encoder fails
func GobEncode(object any) ([]byte, error) {
	// Declare a buffer and a new Gob encoder
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)

	// Encode the object into the buffer
	if err := encoder.Encode(object); err != nil {
		return nil, err
	}

	// Return the bytes in the buffer
	return buffer.Bytes(), nil
}

// GobDecode decodes a stream of bytes into a given object.
// Returns an error if the gob decoded fails.
func GobDecode(data []byte, object any) (any, error) {
	// Declare a new reader from the data and a new Gob decoder
	reader := bytes.NewReader(data)
	decoder := gob.NewDecoder(reader)

	// Decode the data into the object
	if err := decoder.Decode(object); err != nil {
		return nil, err
	}

	// Return the object
	return object, nil
}
