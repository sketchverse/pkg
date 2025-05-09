package security

import (
	"encoding/base64"
	"strings"
)

type EncodingType int

const (
	StdEncoding EncodingType = iota
	URLEncoding
	RawStdEncoding
	RawURLEncoding
)

func EncodeToString(data []byte, encType EncodingType) string {
	switch encType {
	case URLEncoding:
		return base64.URLEncoding.EncodeToString(data)
	case RawStdEncoding:
		return base64.RawStdEncoding.EncodeToString(data)
	case RawURLEncoding:
		return base64.RawURLEncoding.EncodeToString(data)
	default:
		return base64.StdEncoding.EncodeToString(data)
	}
}

func DecodeString(s string, encType EncodingType) ([]byte, error) {
	cleanStr := strings.Map(func(r rune) rune {
		switch {
		case r >= 'A' && r <= 'Z':
			return r
		case r >= 'a' && r <= 'z':
			return r
		case r >= '0' && r <= '9':
			return r
		case r == '+' || r == '-':
			return r
		case r == '/' || r == '_':
			return r
		case r == '=':
			return r
		default:
			return -1
		}
	}, s)

	if strings.ContainsAny(cleanStr, "+/") {
		if strings.HasSuffix(cleanStr, "=") {
			return base64.StdEncoding.DecodeString(cleanStr)
		}
		return base64.RawStdEncoding.DecodeString(cleanStr)
	} else {
		if strings.HasSuffix(cleanStr, "=") {
			return base64.URLEncoding.DecodeString(cleanStr)
		}
		return base64.RawURLEncoding.DecodeString(cleanStr)
	}
}

func MustDecode(s string) []byte {
	data, err := DecodeString(s, StdEncoding)
	if err != nil {
		if d, e := base64.RawStdEncoding.DecodeString(s); e == nil {
			return d
		}
		if d, e := base64.URLEncoding.DecodeString(s); e == nil {
			return d
		}
		if d, e := base64.RawURLEncoding.DecodeString(s); e == nil {
			return d
		}
		return nil
	}
	return data
}

func BatchEncode(dataList [][]byte, encType EncodingType) []string {
	result := make([]string, len(dataList))
	encoder := selectEncoder(encType)

	for i, data := range dataList {
		result[i] = encoder(data)
	}
	return result
}

func selectEncoder(encType EncodingType) func([]byte) string {
	switch encType {
	case URLEncoding:
		return base64.URLEncoding.EncodeToString
	case RawStdEncoding:
		return base64.RawStdEncoding.EncodeToString
	case RawURLEncoding:
		return base64.RawURLEncoding.EncodeToString
	default:
		return base64.StdEncoding.EncodeToString
	}
}
