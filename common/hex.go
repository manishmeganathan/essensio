package common

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

// HexEncode encodes b as a hex string with 0x prefix.
func HexEncode(b []byte) string {
	enc := make([]byte, len(b)*2+2)
	copy(enc, "0x")

	hex.Encode(enc[2:], b)
	return string(enc)
}

// HexDecode decodes a hex string with 0x prefix.
func HexDecode(input string) ([]byte, error) {
	if len(input) == 0 {
		return nil, fmt.Errorf("empty hex string")
	}

	if !hasHexPrefix(input) {
		return nil, fmt.Errorf("hex string without 0x prefix")
	}

	b, err := hex.DecodeString(input[2:])
	if err != nil {
		err = hexError(err)
	}

	return b, err
}

// hasHexPrefix checks if input begins with '0x'
func hasHexPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}

// hexError checks the error and converts into a clear hex decode error
func hexError(err error) error {
	if err, ok := err.(*strconv.NumError); ok {
		switch err.Err {
		case strconv.ErrRange:
			return fmt.Errorf("hex number greater than 64 bits")
		case strconv.ErrSyntax:
			return fmt.Errorf("invalid hex string")
		}
	}

	if _, ok := err.(hex.InvalidByteError); ok {
		return fmt.Errorf("invalid hex string")
	}

	if err == hex.ErrLength {
		return fmt.Errorf("hex string of odd length")
	}

	return err
}
