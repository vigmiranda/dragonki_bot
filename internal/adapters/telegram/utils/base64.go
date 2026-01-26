package utils

import (
	"encoding/base64"
	"strings"
)

func DecodeBase64(data string) []byte {
	// Remove prefixos comuns se vierem junto
	data = strings.TrimPrefix(data, "data:image/png;base64,")
	data = strings.TrimPrefix(data, "data:image/jpeg;base64,")

	decoded, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil
	}

	return decoded
}
