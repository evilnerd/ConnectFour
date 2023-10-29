package server

import (
	"github.com/stretchr/testify/assert"
	"log"
	"strings"
	"testing"
)

func TestGenerateKey(t *testing.T) {

	for it := 0; it < 100; it++ {
		var result1 = GenerateKey(1)

		log.Println(result1)

		var result2 = GenerateKey(2)

		log.Println(result2)

		var result3 = GenerateKey(3)

		log.Println(result3)

		assert.Len(t, strings.Split(result1, " "), 1)
		assert.Len(t, strings.Split(result2, " "), 2)
		assert.Len(t, strings.Split(result3, " "), 3)
	}

}
