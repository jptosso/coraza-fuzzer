package fuzzer

import (
	"encoding/base64"
	"encoding/hex"
	"math/rand"
	"strconv"
)

type Rule func(length uint64) string

var Rules = map[string]Rule{
	"ascii":      RuleAscii,
	"hexascii":   RuleHexAscii,
	"binary":     RuleBinary,
	"base64":     RuleBase64,
	"urlhex":     RuleUrlHex,
	"urlunicode": RuleUrlUnicode,
	"number":     RuleNumber,
}

func RuleAscii(length uint64) string {
	//0x20-0x7E
	b := []byte{}
	for i := uint64(0); i < length; i++ {
		c := byte(rand.Int()%(0x7E-0x20) + 0x20)
		b = append(b, c)
	}
	return string(b)
}

func RuleHexAscii(length uint64) string {
	return hex.EncodeToString(randomBinary(length))
}

func RuleBinary(length uint64) string {
	return string(randomBinary(length))
}

func RuleUnicode(length uint64) string {
	b := []byte{}
	for i := uint64(0); i < length; i++ {
		bytes := make([]byte, 4)
		rand.Read(bytes)
		//n := bytes //TODO limit to % 0x0010FFFF
		b = append(b, bytes...)
	}
	return string(b)
}

func RuleBase64(length uint64) string {
	return base64.StdEncoding.EncodeToString(randomBinary(length))
}

func RuleUrlHex(length uint64) string {
	if length < 3 {
		return ""
	}
	b := make([]byte, length/3)
	for i := uint64(0); i < length/3; i++ {
		b2 := rand.Int() % 255
		bb := []byte(hex.EncodeToString([]byte{byte(b2)}))
		b = append(b, '%')
		b = append(b, bb...)
	}
	return string(b)
}

func RuleUrlUnicode(length uint64) string {
	if length < 6 {
		return ""
	}
	b := make([]byte, length/6)
	for i := uint64(0); i < length/6; i++ {
		b1 := byte(rand.Int() % 255)
		b2 := byte(rand.Int() % 255)
		bb := []byte(hex.EncodeToString([]byte{b1, b2}))
		b = append(b, '%')
		b = append(b, 'u')
		b = append(b, bb...)
	}
	return string(b)
}

func RuleNumber(length uint64) string {
	r := ""
	for i := 0; i < 100; i++ {
		// TODO, too heavy
		r += strconv.Itoa(i)
	}
	return ""
}

func randomBinary(length uint64) []byte {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return []byte{}
	}
	return bytes
}
