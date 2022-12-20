package utils

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
	"unicode"

	"github.com/decred/base58"
	"golang.org/x/crypto/blake2b"
)

var (
	SSPrefix        = []byte{0x53, 0x53, 0x35, 0x38, 0x50, 0x52, 0x45}
	SubstratePrefix = []byte{0x2a}
	CessPrefix      = []byte{0x50, 0xac}
)

func DecodePublicKeyOfCessAccount(address string) ([]byte, error) {
	err := VerityAddress(address, CessPrefix)
	if err != nil {
		return nil, errors.New("invalid account")
	}
	data := base58.Decode(address)
	if len(data) != (34 + len(CessPrefix)) {
		return nil, errors.New("public key decoding failed")
	}
	return data[len(CessPrefix) : len(data)-2], nil
}
func VerityAddress(address string, prefix []byte) error {
	decodeBytes := base58.Decode(address)
	if len(decodeBytes) != (34 + len(prefix)) {
		return errors.New("base58 decode error")
	}
	if decodeBytes[0] != prefix[0] {
		return errors.New("prefix valid error")
	}
	pub := decodeBytes[len(prefix) : len(decodeBytes)-2]

	data := append(prefix, pub...)
	input := append(SSPrefix, data...)
	ck := blake2b.Sum512(input)
	checkSum := ck[:2]
	for i := 0; i < 2; i++ {
		if checkSum[i] != decodeBytes[32+len(prefix)+i] {
			return errors.New("checksum valid error")
		}
	}
	if len(pub) != 32 {
		return errors.New("decode public key length is not equal 32")
	}
	return nil
}

func GetStringSize(size int64) string {
	return strconv.FormatInt(size, 10)
}

func GetIntSize(size string) (int64, error) {
	return strconv.ParseInt(size, 10, 64)
}

func MatchStatus(status, stype string) bool {
	return status == stype
}

func StringNotEmpty(args ...string) bool {
	for _, s := range args {
		if s == "" {
			return false
		}
	}
	return true
}

func Camel2Case(name string) string {
	buffer := bytes.Buffer{}
	for i, r := range name {
		if unicode.IsUpper(r) {
			if i != 0 {
				buffer.WriteRune('_')
			}
			buffer.WriteRune(unicode.ToLower(r))
		} else {
			buffer.WriteRune(r)
		}
	}
	return buffer.String()
}

func IsImage(ext string) bool {
	ext = strings.ToLower(ext)
	switch ext {
	case ".jpg", ".jpeg", ".png", ".tiff", ".gif", ".bmp", ".svg":
		return true
	}
	return false
}
