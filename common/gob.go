package common

import (
	"bytes"
	"encoding/gob"
)

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
