package utils

import (
	"log"
	"net/url"
	"strings"
)

func ParseCallbackData(data string) (string, map[string]string) {
	result := make(map[string]string)

	parts := strings.SplitN(data, "?", 2)
	action := parts[0]

	if len(parts) < 2 {
		return action, result
	}

	values, err := url.ParseQuery(parts[1])
	if err != nil {
		log.Println("❌ Erro ao parsear callback:", err)
		return action, result
	}

	for k, v := range values {
		if len(v) > 0 {
			result[k] = v[0]
		}
	}

	return action, result
}
